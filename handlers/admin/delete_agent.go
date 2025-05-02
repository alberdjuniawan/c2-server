// deleteAgent.go
package admin

import (
	"c2/database"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// DeleteAgentHandler menghapus agent berdasarkan ID
func DeleteAgentHandler(w http.ResponseWriter, r *http.Request) {
	// Mengambil agent ID dari URL
	vars := mux.Vars(r)
	agentID := vars["id"]

	// Panggil fungsi DeleteAgent dari database
	err := database.DeleteAgent(agentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete agent: %v", err), http.StatusInternalServerError)
		return
	}

	// Kembalikan respon sukses
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Agent deleted successfully"))
}
