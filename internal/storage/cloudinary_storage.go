package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"vp_backend/internal/config"
)

type CloudinaryStorage struct{}

func NewCloudinaryStorage() *CloudinaryStorage {
	return &CloudinaryStorage{}
}

func (c *CloudinaryStorage) Upload(file *multipart.FileHeader, folder string) (string, error) {

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	uploadResult, err := config.Cloudinary.Upload.Upload(
		context.Background(),
		f,
		uploader.UploadParams{
			Folder: folder,
		},
	)

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func (c *CloudinaryStorage) Delete(imageUrl string) error {

	// Extract public_id dari URL Cloudinary
	// Contoh URL:
	// https://res.cloudinary.com/xxx/image/upload/v123/property/abc.jpg

	parts := strings.Split(imageUrl, "/upload/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid cloudinary url")
	}

	path := parts[1]

	// remove version (v123/)
	pathParts := strings.SplitN(path, "/", 2)
	if len(pathParts) < 2 {
		return fmt.Errorf("invalid cloudinary path")
	}

	publicID := strings.TrimSuffix(pathParts[1], ".jpg")
	publicID = strings.TrimSuffix(publicID, ".png")
	publicID = strings.TrimSuffix(publicID, ".jpeg")
	publicID = strings.TrimSuffix(publicID, ".webp")

	_, err := config.Cloudinary.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID: publicID,
		},
	)

	return err
}