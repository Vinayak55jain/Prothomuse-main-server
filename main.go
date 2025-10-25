package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Metric struct {
	ProjectId    string `json:"project_id"`
	Route        string `json:"route"`
	Method       string `json:"method"`
	StatusCode   int    `json:"status_code"`
	ResponseTime int    `json:"response_time"`
	Timestamp    int64  `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	metrics []Metric
	mu      sync.Mutex
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	http.HandleFunc("/stream", handleWebsocket)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mu.Lock()
		defer mu.Unlock()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total":   len(metrics),
			"metrics": metrics,
		})
	})

	log.Println("🚀 Server running on http://localhost:8080")
	log.Println("📡 WebSocket: ws://localhost:8080/stream")
	log.Println("🔍 Health: http://localhost:8080/health")
	log.Println("📊 Metrics: http://localhost:8080/metrics")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket Read Error:", err)
			break
		}

		var m Metric
		if err := json.Unmarshal(message, &m); err != nil {
			log.Println("JSON Unmarshal Error:", err)
			continue
		}

		mu.Lock()
		metrics = append(metrics, m)
		mu.Unlock()

		log.Printf("Received Metric: %+v\n", m)

		ack := map[string]string{"status": "received"}
		ackJSON, _ := json.Marshal(ack)
		if err := conn.WriteMessage(websocket.TextMessage, ackJSON); err != nil {
			log.Println("WebSocket Write Error:", err)
			break
		}
	}

	log.Println("WebSocket connection closed")
}