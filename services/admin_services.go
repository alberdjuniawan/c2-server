package services

import (
	"c2/database"
	"c2/utils"
	"errors"
	"log"
)

func RegisterAdmin(username, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return errors.New("error hashing password")
	}

	db := database.Connect()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("username already taken")
	}

	_, err = tx.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("[*] Admin registered successfully:", username)
	return nil
}
