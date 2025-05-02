package admin

import (
	"c2/database"
	"encoding/json"
	"log"
	"net/http"
)

// GetAllAgentsHandler menangani permintaan untuk mengambil semua agent
func GetAllAgentsHandler(w http.ResponseWriter, r *http.Request) {
	agents, err := database.GetAllAgents()
	if err != nil {
		log.Println("Gagal mengambil daftar agent:", err)
		http.Error(w, "Gagal mengambil daftar agent", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}
