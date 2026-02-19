package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// WebSocketClient represents a connected websocket client
type WebSocketClient struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

// WebSocketHub manages websocket connections
type WebSocketHub struct {
	clients    map[*WebSocketClient]bool
	broadcast  chan []byte
	register   chan *WebSocketClient
	unregister chan *WebSocketClient
	mutex      sync.RWMutex
}

var hub = &WebSocketHub{
	clients:    make(map[*WebSocketClient]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *WebSocketClient),
	unregister: make(chan *WebSocketClient),
}

// Run starts the WebSocket hub
func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					// Client's send buffer is full, skip this message
					// Client will be cleaned up by unregister
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// BroadcastMessage sends a message to all connected clients
func BroadcastMessage(messageType string, data interface{}) {
	message := map[string]interface{}{
		"type":      messageType,
		"data":      data,
		"timestamp": time.Now().Unix(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return
	}

	hub.broadcast <- jsonMessage
}

// HandleWebSocket handles websocket connections
func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	client := &WebSocketClient{
		ID:   c.Query("id"),
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump reads messages from the websocket connection
func (c *WebSocketClient) readPump() {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// Handle incoming messages
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		// Process message based on type
		if msgType, ok := msg["type"].(string); ok {
			switch msgType {
			case "ping":
				c.Send <- []byte(`{"type":"pong"}`)
			}
		}
	}
}

// writePump writes messages to the websocket connection
func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// StartSystemMonitoring starts broadcasting system stats
func StartSystemMonitoring() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Get system stats
			cpuPercent, _ := cpu.Percent(0, false)
			memInfo, _ := mem.VirtualMemory()

			stats := map[string]interface{}{
				"cpu": map[string]interface{}{
					"usage": cpuPercent,
				},
				"memory": map[string]interface{}{
					"total":       memInfo.Total,
					"used":        memInfo.Used,
					"free":        memInfo.Free,
					"usedPercent": memInfo.UsedPercent,
				},
				"timestamp": time.Now().Unix(),
			}

			BroadcastMessage("system_stats", stats)
		}
	}()
}

// InitWebSocket initializes the WebSocket hub
func InitWebSocket() {
	go hub.Run()
	go StartSystemMonitoring()
}
