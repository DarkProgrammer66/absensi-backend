package models

import (
	"time"

	"gorm.io/gorm"
)

// User merepresentasikan tabel 'users' di database
type User struct {
	gorm.Model        // Menyediakan field bawaan: ID, CreatedAt, UpdatedAt, DeletedAt
	NIS      string `gorm:"unique;not null" json:"nis"` // Nomor Induk Siswa/Guru (harus unik)
	Nama     string `json:"nama"`
	Password string `json:"-"`                         // "-" artinya field password tidak akan dikembalikan saat serialisasi JSON
	Role     string `gorm:"default:siswa" json:"role"` // 'siswa' atau 'admin'
}

// Absensi merepresentasikan tabel 'absensis' di database
type Absensi struct {
	gorm.Model
	UserID     uint      `json:"user_id"` // Foreign key ke tabel users
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tanggal    time.Time `json:"tanggal"`
	JamMasuk   time.Time `json:"jam_masuk"`
	JamKeluar  *time.Time `json:"jam_keluar"` // Menggunakan pointer (*) agar bisa ber-nilai NULL jika belum absen keluar
	Status     string    `json:"status"`     // "Hadir", "Izin", "Sakit"
}