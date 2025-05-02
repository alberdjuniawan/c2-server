package admin

import (
	"c2/database"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler untuk memperbarui Tags dan Notes agent
func UpdateTagsAndNotes(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID agent dari URL
	vars := mux.Vars(r)
	agentID := vars["agent_id"]

	// Menangani parsing data JSON untuk Tags dan Notes
	var requestData struct {
		Tags  []string `json:"tags"`
		Notes string   `json:"notes"`
	}
	// Decode request body ke struct
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Panggil fungsi untuk memperbarui Tags dan Notes
	err := database.UpdateTagsAndNotes(agentID, requestData.Tags, requestData.Notes)
	if err != nil {
		http.Error(w, "Failed to update agent information", http.StatusInternalServerError)
		return
	}

	// Kirim response sukses
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Agent tags and notes updated successfully"))
}
