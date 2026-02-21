package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// InitDB menginisialisasi koneksi database MySQL
// menggunakan konfigurasi dari environment variable.
//
// Environment variable yang digunakan:
// - DB_USERNAME : username database
// - DB_PASSWORD : password database
// - DB_HOST     : host database
// - DB_PORT     : port database
// - DB_NAME     : nama database
//
// Jika koneksi atau proses ping database gagal,
// aplikasi akan dihentikan.
func InitDB() *sql.DB {

	// Menyusun konfigurasi DSN MySQL
	dbConfig := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		GetEnv("DB_USERNAME", "root"),
		GetEnv("DB_PASSWORD", ""),
		GetEnv("DB_HOST", "mysql.railway.internal"),
		GetEnv("DB_PORT", "3306"),
		GetEnv("DB_NAME", "victoria_property"),
	)

	// Membuka koneksi database
	db, err := sql.Open("mysql", dbConfig)
	if err != nil {
		log.Fatal("DB open error:", err)
	}

	// Memastikan koneksi database dapat digunakan
	if err := db.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}

	log.Println("Database connected")
	return db
}

