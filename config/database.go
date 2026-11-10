package config

import (
	"absensi-backend/models"
	"fmt"

	// ✅ Gunakan sub-package dialek sqlite untuk GORM dari glebarez
	"github.com/glebarez/go-sqlite/gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	// Menggunakan sqlite.Open dari package gorm di atas
	DB, err = gorm.Open(sqlite.Open("absensi.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Detail Error:", err)
		panic("Gagal terhubung ke database!")
	}

	fmt.Println("Database berhasil terhubung.")
	DB.AutoMigrate(&models.User{}, &models.Absensi{})
}