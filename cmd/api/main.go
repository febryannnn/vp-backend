package main

import (
	"os"
	"time"

	"vp_backend/internal/config"
	"vp_backend/internal/delivery/http"
	"vp_backend/internal/delivery/http/handler"
	"vp_backend/internal/repository"
	"vp_backend/internal/service"
	"vp_backend/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// main merupakan entry point aplikasi backend Victoria Property.
// Fungsi ini bertanggung jawab untuk:
// - Memuat konfigurasi environment
// - Menginisialisasi database
// - Menyusun dependency (repository, service, handler)
// - Mendaftarkan HTTP routes
// - Menjalankan HTTP server
func main() {

	// Memuat variabel environment dari file .env
	config.LoadEnv()

	// Inisialisasi koneksi database
	db := config.InitDB()
	defer db.Close()

	// Inisialisasi Gin HTTP server dengan default middleware
	// (logger dan recovery)
	r := gin.Default()

	// ==========================
	// STATIC FILES (GLOBAL)
	// ==========================
	r.Static("/static", "./public/uploads")

	config := cors.Config{
		// Hanya izinkan domain frontend resmi kamu
		AllowOrigins: []string{
			"http://localhost:3000",         // Local development
			"https://victoria-properti.com", // Production
			"https://victoria-property-frontend-igi8.vercel.app",
		},
		// Batasi method yang boleh digunakan
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		// Batasi header yang boleh dikirim oleh frontend
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		// Izinkan browser mengirim Cookie atau Auth Header
		AllowCredentials: true,
		// Berapa lama browser boleh menyimpan izin ini (biar gak OPTIONS terus)
		MaxAge: 12 * time.Hour,
	}

	r.Use(cors.New(config))

	imageStorage := storage.NewLocalStorage("./public/uploads", "/static")

	// ==========================
	// REPOSITORY INITIALIZATION
	// ==========================

	// Repository bertanggung jawab untuk interaksi langsung
	// dengan database
	userRepo := &repository.UserRepository{DB: db}
	propertyRepo := &repository.PropertyRepository{DB: db}
	favoriteRepo := &repository.FavoriteRepository{DB: db}
	imageRepo := &repository.ImageRepository{DB: db}

	// ==========================
	// SERVICE INITIALIZATION
	// ==========================

	// Service berisi business logic aplikasi
	authService := &service.AuthService{UserRepo: userRepo}
	userService := &service.UserService{UserRepo: userRepo}
	propertyService := &service.PropertyService{PropertyRepo: propertyRepo, ImageRepo: imageRepo, Storage: imageStorage}
	favoriteService := &service.FavoriteService{FavoriteRepo: favoriteRepo}
	imageService := &service.ImageService{ImageRepo: imageRepo, PropertyRepo: propertyRepo, Storage: imageStorage}

	// ==========================
	// HANDLER INITIALIZATION
	// ==========================

	// Handler berfungsi sebagai layer HTTP
	// yang menangani request dan response
	authHandler := &handler.AuthHandler{AuthService: authService}
	userHandler := &handler.UserHandler{UserService: userService}
	propertyHandler := &handler.PropertyHandler{PropertyService: propertyService}
	favoriteHandler := &handler.FavoriteHandler{FavoriteService: favoriteService}
	imageHandler := &handler.ImageHandler{ImageService: imageService}

	// ==========================
	// INJECT SERVICES TO CONTEXT
	// ==========================
	r.Use(func(c *gin.Context) {
		// Set user_service ke context agar bisa diakses middleware
		c.Set("user_service", userService)
		c.Next()
	})

	// ==========================
	// ROUTE REGISTRATION
	// ==========================

	// Menggabungkan seluruh handler ke dalam
	// satu struct untuk kebutuhan routing
	h := http.Handler{
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		PropertyHandler: propertyHandler,
		FavoriteHandler: favoriteHandler,
		ImageHandler:    imageHandler,
	}

	// Mendaftarkan seluruh endpoint API
	http.RegisterRoutes(r, h)

	// ==========================
	// SERVER STARTUP
	// ==========================

	// Mengambil port dari environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Menjalankan HTTP server
	r.Run(":" + port)
}
