package admin

import (
	"c2/database"
	"encoding/json"
	"log"
	"net/http"
)

// GetAllAgentsHandler godoc
// @Summary Retrieve all registered agents
// @Description Returns a list of all agents currently registered in the system
// @Tags admin
// @Produce json
// @Success 200 {array} database.Agent
// @Failure 500 {object} database.ErrorResponse
// @Router /admin/agents [get]
func GetAllAgentsHandler(w http.ResponseWriter, r *http.Request) {
	agents, err := database.GetAllAgents()
	if err != nil {
		log.Println("Failed to retrieve agent list:", err)
		http.Error(w, "Failed to retrieve agent list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}
