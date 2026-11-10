package handlers

import (
	"absensi-backend/config"
	"absensi-backend/middleware"
	"absensi-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register menyimpan user baru ke database
func Register(c *fiber.Ctx) error {
	var input models.User

	// 1. Parsing body request JSON ke struct input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Format data salah"})
	}

	// 2. Enkripsi password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengamankan password"})
	}
	input.Password = string(hashedPassword)

	// 3. Simpan user ke database via GORM
	if err := config.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "NIS sudah terdaftar"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Registrasi berhasil",
		"data":    input,
	})
}

// Login membuat token JWT untuk user yang valid
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		NIS      string `json:"nis"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	// 1. Cari user berdasarkan NIS di database
	var user models.User
	if err := config.DB.Where("nis = ?", input.NIS).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "NIS atau password salah"})
	}

	// 2. Cocokkan password yang di-input dengan hash di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "NIS atau password salah"})
	}

	// 3. Buat JWT Claims (isi data dalam token)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Kadaluwarsa dalam 24 jam
	}

	// 4. Generate token JWT dengan algoritma HMAC-SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(middleware.SecretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal membuat token"})
	}

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   t,
	})
}