package admin

import (
	"c2/database"
	"c2/utils"
	"database/sql"
	"encoding/json"
	"net/http"
)

// LoginRequest untuk login admin
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginAdmin untuk autentikasi admin
func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.LogError("Invalid input in LoginAdmin")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	db := database.Connect()
	defer db.Close()

	var hashedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username = ?", req.Username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			utils.LogError("Database error in LoginAdmin")
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Verifikasi password menggunakan VerifyPassword
	if !utils.VerifyPassword(hashedPassword, req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Buat JWT token
	token, err := utils.GenerateJWT(req.Username)
	if err != nil {
		utils.LogError("Error generating token in LoginAdmin")
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	utils.LogInfo("Admin logged in successfully: " + req.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
