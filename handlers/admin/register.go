package admin

import (
	"c2/services"
	"c2/utils"
	"encoding/json"
	"net/http"
)

// RegisterRequest untuk pendaftaran admin baru
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterAdmin untuk mendaftarkan admin baru
func RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.LogError("Invalid input in RegisterAdmin")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Panggil service untuk melakukan registrasi admin
	err = services.RegisterAdmin(req.Username, req.Password)
	if err != nil {
		utils.LogError("Error registering admin")
		http.Error(w, "Error registering admin", http.StatusInternalServerError)
		return
	}

	utils.LogInfo("Admin registered successfully: " + req.Username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin registered successfully"})
}
