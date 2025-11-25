package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"prothomuse-server/internal/handler"
	"prothomuse-server/internal/repository"
	"prothomuse-server/internal/services"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

// Metric represents the data from middleware
type Metric struct {
	ProjectID    string `json:"projectId"`
	Route        string `json:"route"`
	Method       string `json:"method"`
	StatusCode   int    `json:"statusCode"`
	ResponseTime int64  `json:"responseTime"`
	Timestamp    int64  `json:"timestamp"` // Unix timestamp in milliseconds
}

var db *sql.DB

func init() {
	connStr := "host=localhost port=5432 user=postgres password=suddendeath123@ dbname=postgres sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL database (postgres)")
	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections for now
	},
}

// Store metrics in memory
var metrics []Metric

func main() {
	// Initialize repository, service, and handler
	userRepo := repository.NewUserRepository(db)

	// Create users table if it doesn't exist
	if err := userRepo.CreateTable(); err != nil {
		log.Printf("⚠️  Warning: Could not create users table: %v", err)
	} else {
		log.Println("✅ Users table ready")
	}

	authService := services.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// Authentication endpoints
	http.HandleFunc("/api/auth/register", authHandler.RegisterUser)
	http.HandleFunc("/api/auth/login", authHandler.Login)
	http.HandleFunc("/api/auth/update", authHandler.UpdateUser)
	http.HandleFunc("/api/auth/validate-apikey", authHandler.ValidateAPIKey)
	http.HandleFunc("/api/auth/validate-jwt", authHandler.ValidateJWT)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","service":"prothomuse-health-server"}`))
	})

	// WebSocket endpoint - middleware connects here
	http.HandleFunc("/stream", handleWebSocket)

	// API to view metrics (for testing/dashboard)
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total":   len(metrics),
			"metrics": metrics,
		})
	})

	// Get metrics by project ID
	http.HandleFunc("/metrics/", func(w http.ResponseWriter, r *http.Request) {
		projectID := r.URL.Path[len("/metrics/"):]

		var projectMetrics []Metric
		for _, m := range metrics {
			if m.ProjectID == projectID {
				projectMetrics = append(projectMetrics, m)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"projectId": projectID,
			"total":     len(projectMetrics),
			"metrics":   projectMetrics,
		})
	})

	log.Println("╔══════════════════════════════════════════════════╗")
	log.Println("║   Prothomuse Health Monitoring Server           ║")
	log.Println("╚══════════════════════════════════════════════════╝")
	log.Println("")
	log.Println("🚀 Server running on http://localhost:8080")
	log.Println("")
	log.Println("� Authentication API Endpoints:")
	log.Println("   POST   /api/auth/register           - Register a new user")
	log.Println("   POST   /api/auth/login              - Login and get JWT token")
	log.Println("   PUT    /api/auth/update             - Update user profile (requires Bearer token)")
	log.Println("   GET    /api/auth/validate-apikey    - Validate API key (Authorization: ApiKey <key>)")
	log.Println("   GET    /api/auth/validate-jwt       - Validate JWT token (Authorization: Bearer <token>)")
	log.Println("")
	log.Println("� Health & Metrics Endpoints:")
	log.Println("   GET    /health                      - Health check")
	log.Println("   WS     /stream                      - WebSocket for metrics")
	log.Println("   GET    /metrics                     - View all metrics")
	log.Println("   GET    /metrics/{projectId}        - View metrics by project")
	log.Println("")
	log.Println("Waiting for connections...")

	http.ListenAndServe(":8080", nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("✅ New middleware connected!")

	// Read messages from middleware
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ Middleware disconnected")
			break
		}

		// Parse metric data
		var metric Metric
		if err := json.Unmarshal(message, &metric); err != nil {
			log.Println("⚠️  Failed to parse metric:", err)
			continue
		}

		// Store metric in memory
		metrics = append(metrics, metric)
		// Store metric in PostgreSQL
		_, err = db.Exec(`
	INSERT INTO metrics (project_id, route, method, status_code, response_time, timestamp)
	VALUES ($1, $2, $3, $4, $5, $6)
`, metric.ProjectID, metric.Route, metric.Method, metric.StatusCode, metric.ResponseTime, metric.Timestamp)

		if err != nil {
			log.Println("❌ Failed to insert metric into DB:", err)
		} else {
			log.Println("✅ Metric saved to DB")
		}

		// Log received metric
		log.Printf("📊 [%s] %s %s -> %d (%dms)",
			metric.ProjectID,
			metric.Method,
			metric.Route,
			metric.StatusCode,
			metric.ResponseTime,
		)

		// Send acknowledgment back to middleware
		ack := map[string]string{
			"status":  "received",
			"message": "Metric saved successfully",
		}
		ackJSON, _ := json.Marshal(ack)
		conn.WriteMessage(websocket.TextMessage, ackJSON)
	}
}
