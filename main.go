package main

import (
	"absensi-backend/config"
	"absensi-backend/handlers"
	"absensi-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// 1. Inisialisasi Koneksi Database
	config.ConnectDB()

	// 2. Inisialisasi Framework Fiber
	app := fiber.New()

	// 3. Mengaktifkan CORS (Cross-Origin Resource Sharing)
	// Penting agar aplikasi Flutter di HP/Emulator bisa mengakses backend ini
	app.Use(cors.New())

	// --- PUBLIC ROUTES (Bisa diakses tanpa login) ---
	api := app.Group("/api")
	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)

	// --- PROTECTED ROUTES (Wajib menyertakan Token JWT) ---
	protected := api.Group("/", middleware.Protected())
	protected.Post("/absen/masuk", handlers.AbsenMasuk)
	protected.Post("/absen/keluar", handlers.AbsenKeluar)
	protected.Get("/absen/riwayat", handlers.GetRiwayatAbsensi)

	// 4. Jalankan Server di Port 8080
	app.Listen(":8080")
}