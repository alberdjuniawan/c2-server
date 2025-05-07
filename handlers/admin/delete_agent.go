package admin

import (
	"c2/database"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// DeleteAgentHandler deletes a registered agent by ID.
//
// @Summary Delete an agent
// @Description Permanently removes an agent from the system by its unique ID.
// @Tags Admin - Agents
// @Param id path string true "Agent ID"
// @Success 200 {string} string "Agent deleted successfully"
// @Failure 500 {string} string "Failed to delete agent"
// @Router /admin/agents/{id} [delete]
func DeleteAgentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["id"]

	err := database.DeleteAgent(agentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete agent: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Agent deleted successfully"))
}
