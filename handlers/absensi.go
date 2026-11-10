package handlers

import (
	"absensi-backend/config"
	"absensi-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AbsenMasuk mencatat waktu kehadiran siswa
func AbsenMasuk(c *fiber.Ctx) error {
	// Ambil userID dari token JWT yang sudah di-set oleh Middleware
	userID := c.Locals("userID").(uint)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Cek apakah user sudah absen hari ini
	var absensiAda models.Absensi
	err := config.DB.Where("user_id = ? AND tanggal = ?", userID, today).First(&absensiAda).Error
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Anda sudah melakukan absen masuk hari ini"})
	}

	// Buat record absensi baru
	absensiBaru := models.Absensi{
		UserID:   userID,
		Tanggal:  today,
		JamMasuk: now,
		Status:   "Hadir",
	}

	config.DB.Create(&absensiBaru)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Absen masuk berhasil",
		"data":    absensiBaru,
	})
}

// AbsenKeluar memperbarui catatan jam keluar
func AbsenKeluar(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Cari record absen masuk hari ini
	var absensi models.Absensi
	if err := config.DB.Where("user_id = ? AND tanggal = ?", userID, today).First(&absensi).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Anda belum melakukan absen masuk hari ini"})
	}

	if absensi.JamKeluar != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Anda sudah melakukan absen keluar hari ini"})
	}

	// Update JamKeluar menggunakan pointer jam saat ini
	config.DB.Model(&absensi).Update("jam_keluar", now)

	return c.JSON(fiber.Map{
		"message": "Absen keluar berhasil",
		"data":    absensi,
	})
}

// GetRiwayatAbsensi menampilkan histori absensi milik user login
func GetRiwayatAbsensi(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var listAbsensi []models.Absensi
	// GORM Preload digunakan untuk mengambil detail relasi data User (JOIN tabel)
	config.DB.Where("user_id = ?", userID).Find(&listAbsensi)

	return c.JSON(fiber.Map{
		"data": listAbsensi,
	})
}