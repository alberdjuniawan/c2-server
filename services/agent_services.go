package services

import (
	"c2/database"
	"c2/utils"
	"errors"
	"log"
	"time"
)

func RegisterAgent(id, ip, hostname, os, arch string) (string, error) {
	db := database.Connect()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM agents WHERE id = ?", id).Scan(&count)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", errors.New("agent already registered")
	}

	token, err := utils.GenerateAgentToken(id)
	if err != nil {
		log.Println("Error generating JWT token for agent:", err)
		return "", err
	}

	_, err = tx.Exec(`
		INSERT INTO agents (id, ip, hostname, os, arch, token, registered_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, ip, hostname, os, arch, token, time.Now())
	if err != nil {
		log.Println("Error inserting agent into database:", err)
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return token, nil
}
