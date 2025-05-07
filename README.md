# C2 Server

A Command & Control (C2) backend server designed to manage remote agents securely and efficiently. This C2 server handles agent registration, command dispatching, heartbeat monitoring, result retrieval, and admin authentication — all with strong security measures and modular design.

## Features

### **Admin Features**

- **Admin Authentication & Registration**
  - JWT-based admin authentication for secure access.
  - Admin can register with bcrypt-hashed passwords for safety.

- **Agent Management**
  - **View all agents**: Admin can list all registered agents with details such as IP, hostname, and OS.
  - **Delete agents**: Admin can remove agents from the system.
  - **Update agent metadata**: Admin can add tags and notes for agents to manage them better.
  
- **Command Management**
  - **Send commands to agents**: Admin can issue commands to specific agents, e.g., capture screenshots, record keystrokes.
  - **Command history tracking**: Track and retrieve command results for each agent.
  - **File download**: Admin can download the result files of executed commands.
  - **Delete command history**: Admin can delete command history records.

### **Agent Features**

- **Agent Registration**
  - Secure agent registration via signature verification to avoid impersonation.
  
- **Heartbeat Monitoring**
  - Periodic heartbeat messages to indicate the agent is online and active.
  
- **Command Result Submission**
  - Agents can send command execution results back to the server.
  - **File upload and encryption**: Secure command result file upload with AES-GCM encryption.
  
- **File Upload Handling**
  - **Filename sanitization**: Prevent malicious file names by sanitizing uploads.
  - **File size restrictions**: Limits on the size of uploaded files to prevent overflow or excessive load.
  - **Allowed extensions validation**: Ensures that only certain file types are accepted.

### **Security Features**

- **Routing & Access Control**
  - **Dual-server routing**: Public server (`:443`) for agents, internal server (`:8443`) for admin access.
  - **Access control by network boundary**: Ensures admin and agent endpoints are properly separated.
  
- **Signature Verification & HMAC Authentication**
  - **Agent signature verification**: Ensures requests from agents are legitimate via HMAC-based signature checks.
  
- **Encryption**
  - **AES-GCM encryption**: Secure storage and transmission of sensitive command results.
  - **File upload encryption**: Encrypts uploaded files before storing them.

- **Security Middleware**
  - **Rate Limiting**: Limits the number of requests to prevent abuse.
  - **Security Headers**: Adds headers to requests to prevent attacks like XSS, CSRF, and clickjacking.
  - **Replay Attack Prevention**: Ensures requests are unique by validating timestamps and using nonces.
  - **Panic Recovery**: Catches unexpected crashes to prevent leaks of sensitive information.

- **Command Integrity & Validation**
  - **Command result integrity**: Ensures that the results sent back from agents are unmodified.
  - **Upload validation**: Ensures files uploaded during command results are valid and safe.

## Database: SQLite

- **SQLite Database**: Uses SQLite as a lightweight, serverless database for storing agent data, user credentials, commands, logs, and more.
  - **Schema definitions**: Models defined in `models.go` represent database tables (`agents`, `users`, `commands`, etc.)
  - **User Management**: Admin accounts are stored with hashed passwords using bcrypt for secure login.
  - **Logging & Command History**: Tracks each agent's activity, including sent commands and results, with timestamps and metadata.
  
## Project Structure

```bash
C2-Server/
├── .github/
│   └── workflows/
│       └── go.yml                   # GitHub Actions CI configuration
├── config/
│   ├── ssl/
│   │   ├── internal.key             # Private key for HTTPS server
│   │   └── internal.pem             # Certificate file for HTTPS server
│   └── config.go                    # Loads environment variables and config
├── database/
│   ├── db.go                        # Database connection and setup
│   └── models.go                    # DB schema definitions (Agent, User, Command, etc.)
├── docs/
│   ├── docs.go                      # Swagger docs init (via swaggo/swag)
│   ├── swagger.json                 # Generated Swagger spec (JSON)
│   └── swagger.yaml                 # Swagger spec (YAML version)
├── handlers/
│   ├── admin/
│   │   ├── agents.go                # Handler to list agents
│   │   ├── command.go               # Handler to send commands to agents
│   │   ├── delete_agent.go          # Handler to remove agent
│   │   ├── download.go              # File download endpoint for command results
│   │   ├── login.go                 # Admin login handler
│   │   ├── register.go              # Admin registration handler
│   │   └── update_meta.go           # Update agent's metadata (tags, notes)
│   ├── agent/
│   │   ├── heartbeat.go             # Heartbeat handler to mark agent online
│   │   ├── register.go              # Agent registration handler
│   │   ├── result.go                # Handler for agent to send back command result
│   │   └── upload.go                # Handler to upload result from agent
│   └── middleware/
│       ├── auth.go                  # JWT validation middleware for admin
│       ├── cors.go                  # Cross-Origin Resource Sharing config
│       ├── rate_limiter.go          # Rate limiting middleware
│       ├── recover.go               # Panic recovery middleware
│       ├── security.go              # Secure headers middleware
│       └── verify_signature.go      # Signature validation for agent request integrity
├── services/
│   ├── admin_services.go           # Admin-side business logic (user, command management)
│   └── agent_services.go           # Agent-side business logic (heartbeat, results)
├── utils/
│   ├── crypto.go                   # AES-256-GCM encryption for file result
│   ├── jwt.go                      # JWT generation and parsing helpers
│   ├── logger.go                   # Logging utility
│   └── password.go                 # Bcrypt password hashing utilities
├── .env                            # Env vars: JWT_SECRET, DB path, etc.
├── .gitignore                      # Git ignored files
├── app.log                         # Server log file (optional, runtime created)
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksum
├── LICENSE                         # Project license (MIT, assumed)
├── main.go                         # Entry point: HTTPS server setup & routing
├── README.md                       # Project documentation (to be generated)
└── test-api.rest                   # REST client file to test endpoints (e.g., VSCode Thunder Client)
```

## Requirements

- **Go 1.18+**
- **SQLite**

## Installation

1. **Clone the repository**
    ```bash
    git clone https://github.com/yourname/c2-server.git
    cd c2-server
    ```

2. **Configure environment variables**
   Create a `.env` file in the root directory (example below):
   ```bash
   DB_PATH=path_to_your_sqlite_db
   JWT_SECRET=your_jwt_secret
   ENCRYPTION_KEY=your_encryption_key
   DOMAIN=yourdomain.com
   AGENT_SECRET=your_agent_secret

   
3. **Run the server**
   ```bash
   go run main.go
   ```

4. **Access the server**
   - **Admin Panel** (localhost): `https://localhost:8443`
   - **Agent Endpoint** (with Let's Encrypt SSL): `https://yourdomain.com:443`
   
   Note: For the agent, you'll need a valid public domain and SSL setup through Let's Encrypt (automatic in your code).

## Endpoints overview
   - **Agent** endpoints are available on port `:443` (Public HTTPS).
   - **Admin** endpoints are available on port `:8443` (Internal HTTPS).

   **Public Endpoints (Agent) - Accessible on Port `:443`**

   ### `POST /agent/register`
   - **Description**: Register a new agent.
   - **Request Body**: 
   - JSON object containing agent details (e.g., signature).
  
   ### `POST /agent/result`
   - **Description**: Submit the results from an agent (e.g., screenshots, keylog data, etc.).
   - **Request Body**: 
     - JSON object containing the result data.

   ### `POST /agent/heartbeat`
   - **Description**: Sends a heartbeat signal to indicate that the agent is still alive and active.
   - **Authentication**: Requires a JWT token in the request header (via `Authorization` header).

   ### `POST /agent/upload`
   - **Description**: Upload results from the agent to the server (e.g., files, screenshots).
   - **Authentication**: Requires a JWT token in the request header (via `Authorization` header).

   **Admin Endpoints - Accessible on Port `:8443`**

   ### Public Endpoints (Admin)

   ### `POST /admin/register`
   - **Description**: Register a new admin.
   - **Request Body**: 
   - JSON object containing admin credentials (e.g., username and password).

   ### `POST /admin/login`
   - **Description**: Admin login to receive a JWT token.
   - **Request Body**: 
   - JSON object containing admin credentials (e.g., username and password).

   **Private Endpoints (Admin) - Requires JWT Authentication**

   ### `GET /admin/agents`
   - **Description**: Retrieve a list of all active agents.
   - **Authentication**: Requires JWT token for authorization.

   ### `DELETE /admin/delete_agent/{agent_id}`
   - **Description**: Delete an agent by its ID.
   - **Authentication**: Requires JWT token for authorization.

   ### `PATCH /admin/update_meta/{agent_id}`
   - **Description**: Update tags and notes associated with an agent by its ID.
   - **Authentication**: Requires JWT token for authorization.

   ### `POST /admin/command/{agent_id}/send`
   - **Description**: Send a command to a specific agent.
   - **Request Body**: 
     - JSON object containing the command to be executed by the agent.
   - **Authentication**: Requires JWT token for authorization.

   ### `GET /admin/command/{agent_id}`
   - **Description**: Get the list of commands executed by a specific agent.
   - **Authentication**: Requires JWT token for authorization.

   ### `GET /admin/command/{command_id}/download`
   - **Description**: Download the result of a specific command by command ID.
   - **Authentication**: Requires JWT token for authorization.

   ### `DELETE /admin/command/{id}`
   - **Description**: Delete a command by its ID.
   - **Authentication**: Requires JWT token for authorization.

## Testing the API

You can test these endpoints using any REST client of your choice (such as Postman, Insomnia, or the REST Client extension in VS Code) with the provided `test-api.rest` file available in the repository. This file contains pre-configured HTTP requests for each endpoint, making it easy to interact with and verify the server functionality.

To use the file:
1. Open the `test-api.rest` file in your preferred REST client.
2. Update any necessary variables (such as agent IDs or JWT tokens).
3. Send the requests directly to the server running on port `:443` (for Agent) or `:8443` (for Admin).

Ensure the server is running and accessible before testing the API.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.