package services

import (
	"c2/database"
	"c2/utils"
	"errors"
	"log"
	"time"
)

// RegisterAgent untuk menangani logika registrasi agent
func RegisterAgent(id, ip, hostname, os, arch string) (string, error) {
	// Koneksi ke database dan menggunakan transaksi untuk memastikan integritas data
	db := database.Connect()
	defer db.Close()

	tx, err := db.Begin() // Memulai transaksi
	if err != nil {
		return "", err
	}
	defer tx.Rollback() // Jika ada kesalahan, rollback transaksi

	// Cek apakah agent sudah terdaftar
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM agents WHERE id = ?", id).Scan(&count)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", errors.New("agent already registered")
	}

	// Generate token JWT untuk agent
	token, err := utils.GenerateAgentToken(id)
	if err != nil {
		log.Println("Error generating JWT token for agent:", err)
		return "", err
	}

	// Insert agent baru ke database
	_, err = tx.Exec(`
		INSERT INTO agents (id, ip, hostname, os, arch, token, registered_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, ip, hostname, os, arch, token, time.Now())
	if err != nil {
		log.Println("Error inserting agent into database:", err)
		return "", err
	}

	// Commit transaksi jika tidak ada error
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return token, nil
}
