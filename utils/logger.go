package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogInfo mencatat log dengan level info
func LogInfo(message string) {
	log.Printf("[INFO] [%s] %s", time.Now().Format(time.RFC3339), message)
}

// LogWarning mencatat log dengan level warning
func LogWarning(message string) {
	log.Printf("[WARNING] [%s] %s", time.Now().Format(time.RFC3339), message)
}

// LogError mencatat log dengan level error
func LogError(message string) {
	log.Printf("[ERROR] [%s] %s", time.Now().Format(time.RFC3339), message)
}

// SetupLogger setup untuk log file jika diperlukan
func SetupLogger() (*os.File, error) {
	// Create or open log file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("unable to open log file: %v", err)
	}

	// Redirect log output ke file dan juga ke konsol
	log.SetOutput(logFile)

	// Set format log untuk waktu yang lebih mudah dibaca
	log.SetFlags(log.LstdFlags)

	return logFile, nil
}
