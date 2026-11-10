package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("kunci-rahasia-sekolah-123") // Gunakan secret key yang aman

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak ditemukan",
			})
		}

		// Header format: "Bearer <TOKEN>"
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 2. Parse dan validasi token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak valid atau kadaluwarsa",
			})
		}

		// 3. Simpan data user (claims) dari token ke dalam konteks lokal Fiber
		claims := token.Claims.(jwt.MapClaims)
		c.Locals("userID", uint(claims["user_id"].(float64)))

		// Lanjutkan ke handler utama
		return c.Next()
	}
}