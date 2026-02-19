package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type terminalSession struct {
	cmd    *exec.Cmd
	ptmx   *os.File
	conn   *websocket.Conn
	mu     sync.Mutex
	closed bool
}

var (
	sessions   = make(map[string]*terminalSession)
	sessionsMu sync.RWMutex
)

type TerminalMessage struct {
	Type string `json:"type"` // input, resize, ping
	Data string `json:"data"`
	Rows uint16 `json:"rows,omitempty"`
	Cols uint16 `json:"cols,omitempty"`
}

// HandleTerminalWebSocket handles WebSocket connections for terminal sessions
func HandleTerminalWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Create a new terminal session
	session, err := createTerminalSession(conn)
	if err != nil {
		log.Printf("Failed to create terminal session: %v", err)
		conn.WriteJSON(map[string]string{
			"type":  "error",
			"error": "Failed to start terminal: " + err.Error(),
		})
		return
	}
	defer session.cleanup()

	// Start goroutine to read from PTY and send to WebSocket
	go session.readFromPTY()

	// Read from WebSocket and write to PTY
	session.readFromWebSocket()
}

func createTerminalSession(conn *websocket.Conn) (*terminalSession, error) {
	// Determine shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	// Create command
	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// Start the command with a pty
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	// Set initial size
	if err := pty.Setsize(ptmx, &pty.Winsize{
		Rows: 24,
		Cols: 80,
	}); err != nil {
		ptmx.Close()
		return nil, err
	}

	session := &terminalSession{
		cmd:  cmd,
		ptmx: ptmx,
		conn: conn,
	}

	return session, nil
}

func (s *terminalSession) readFromPTY() {
	buf := make([]byte, 8192)
	for {
		n, err := s.ptmx.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from PTY: %v", err)
			}
			break
		}

		s.mu.Lock()
		if s.closed {
			s.mu.Unlock()
			break
		}

		// Send data to WebSocket
		if err := s.conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
			log.Printf("Error writing to WebSocket: %v", err)
			s.mu.Unlock()
			break
		}
		s.mu.Unlock()
	}
}

func (s *terminalSession) readFromWebSocket() {
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse message
		var msg TerminalMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			// If not JSON, treat as raw input
			s.ptmx.Write(message)
			continue
		}

		switch msg.Type {
		case "input":
			// Write input to PTY
			s.ptmx.Write([]byte(msg.Data))
		case "resize":
			// Resize PTY
			if msg.Rows > 0 && msg.Cols > 0 {
				pty.Setsize(s.ptmx, &pty.Winsize{
					Rows: msg.Rows,
					Cols: msg.Cols,
				})
			}
		case "ping":
			// Respond to ping
			s.mu.Lock()
			s.conn.WriteJSON(map[string]string{"type": "pong"})
			s.mu.Unlock()
		}
	}
}

func (s *terminalSession) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}
	s.closed = true

	// Close PTY
	if s.ptmx != nil {
		s.ptmx.Close()
	}

	// Kill process
	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Process.Kill()
		s.cmd.Wait()
	}

	// Close WebSocket
	if s.conn != nil {
		s.conn.Close()
	}
}

// ListTerminalSessions returns active terminal sessions (for future multi-session support)
func ListTerminalSessions(c *gin.Context) {
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()

	sessionList := make([]map[string]interface{}, 0)
	for id := range sessions {
		sessionList = append(sessionList, map[string]interface{}{
			"id":        id,
			"createdAt": time.Now(), // In real implementation, track creation time
		})
	}

	c.JSON(http.StatusOK, sessionList)
}
