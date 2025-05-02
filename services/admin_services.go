package services

import (
	"c2/database"
	"c2/utils"
	"errors"
	"log"
)

// RegisterAdmin untuk menangani logika registrasi admin
func RegisterAdmin(username, password string) error {
	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return errors.New("error hashing password")
	}

	// Koneksi ke database dan menggunakan transaction untuk memastikan integritas data
	db := database.Connect()
	defer db.Close()

	tx, err := db.Begin() // Memulai transaksi
	if err != nil {
		return err
	}
	defer tx.Rollback() // Jika ada kesalahan, rollback transaksi

	// Periksa apakah username sudah ada
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("username already taken")
	}

	// Masukkan data admin baru
	_, err = tx.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		return err
	}

	// Commit transaksi jika tidak ada error
	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("[*] Admin registered successfully:", username)
	return nil
}
