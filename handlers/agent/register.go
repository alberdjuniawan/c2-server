package agent

import (
	"c2/database"
	"c2/services"
	"c2/utils"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// RegisterRequest untuk pendaftaran agent baru (tanpa ID dan IP)
type RegisterRequest struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
}

// RegisterResponse untuk response setelah registrasi agent
type RegisterResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

// RegisterAgent untuk mendaftarkan agent baru
func RegisterAgent(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.LogError("Invalid input in RegisterAgent")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Generate UUID untuk ID agent
	newID := uuid.New().String()

	// Ambil IP dari request (gunakan r.RemoteAddr atau X-Forwarded-For jika menggunakan proxy)
	ip := r.RemoteAddr // Ambil IP asli dari request

	// Validasi apakah agent dengan ID yang sama sudah terdaftar
	db := database.Connect()
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM agents WHERE id = ?", newID).Scan(&count)
	if err != nil {
		utils.LogError("Error checking if agent already exists: " + err.Error())
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		utils.LogError("Agent already registered: " + newID)
		http.Error(w, "Agent already registered", http.StatusConflict)
		return
	}

	// Panggil service untuk melakukan registrasi agent
	token, err := services.RegisterAgent(newID, ip, req.Hostname, req.OS, req.Arch)
	if err != nil {
		utils.LogError("Error registering agent: " + err.Error())
		http.Error(w, "Error registering agent", http.StatusInternalServerError)
		return
	}

	utils.LogInfo("Agent registered successfully: " + newID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{
		Message: "Agent registered successfully",
		Token:   token, // Kirim token yang dihasilkan
	})
}
