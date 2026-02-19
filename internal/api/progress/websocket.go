package progress

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Client represents a WebSocket client
type Client struct {
	ID        string
	ProjectID uint
	UserID    uint
	Conn      *websocket.Conn
	Send      chan []byte
	Hub       *Hub
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	// Registered clients by project
	Projects map[uint]*ProjectRoom

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Broadcast requests
	Broadcast chan *BroadcastMessage

	// Mutex for thread safety
	mu sync.RWMutex
}

// ProjectRoom represents a room for a specific project
type ProjectRoom struct {
	ProjectID uint
	Clients   map[*Client]bool
	Cursors   map[string]*CursorPosition // userID -> cursor
	mu        sync.RWMutex
}

// CursorPosition represents a user's cursor position
type CursorPosition struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	TaskID *uint   `json:"task_id,omitempty"`
}

// BroadcastMessage represents a message to broadcast
type BroadcastMessage struct {
	ProjectID uint
	Message   []byte
	Exclude   *Client // Optional: exclude this client from broadcast
}

// WSEvent represents a WebSocket event
type WSEvent struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	UserID    uint                   `json:"user_id"`
	ProjectID uint                   `json:"project_id"`
	Timestamp string                 `json:"timestamp"`
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		Projects:  make(map[uint]*ProjectRoom),
		Register:  make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast: make(chan *BroadcastMessage),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastToProject(message)
		}
	}
}

// registerClient adds a client to the appropriate project room
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.Projects[client.ProjectID]
	if !exists {
		room = &ProjectRoom{
			ProjectID: client.ProjectID,
			Clients:   make(map[*Client]bool),
			Cursors:   make(map[string]*CursorPosition),
		}
		h.Projects[client.ProjectID] = room
	}

	room.mu.Lock()
	room.Clients[client] = true
	room.mu.Unlock()

	log.Printf("[WS] Client %s joined project %d", client.ID, client.ProjectID)

	// Notify others in the room
	h.notifyUserJoined(client)
}

// unregisterClient removes a client from the project room
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.Projects[client.ProjectID]
	if !exists {
		return
	}

	room.mu.Lock()
	if _, ok := room.Clients[client]; ok {
		delete(room.Clients, client)
		close(client.Send)
	}
	room.mu.Unlock()

	// Clean up empty rooms
	if len(room.Clients) == 0 {
		delete(h.Projects, client.ProjectID)
	}

	log.Printf("[WS] Client %s left project %d", client.ID, client.ProjectID)

	// Notify others in the room
	h.notifyUserLeft(client)
}

// broadcastToProject sends a message to all clients in a project room
func (h *Hub) broadcastToProject(message *BroadcastMessage) {
	h.mu.RLock()
	room, exists := h.Projects[message.ProjectID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	for client := range room.Clients {
		// Skip excluded client (sender)
		if message.Exclude != nil && client == message.Exclude {
			continue
		}

		select {
		case client.Send <- message.Message:
		default:
			// Client's send channel is full, close connection
			log.Printf("[WS] Client %s send channel full, closing", client.ID)
			h.Unregister <- client
		}
	}
}

// notifyUserJoined broadcasts user joined event
func (h *Hub) notifyUserJoined(client *Client) {
	event := WSEvent{
		Type:      "user:joined",
		UserID:    client.UserID,
		ProjectID: client.ProjectID,
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"user": map[string]interface{}{
				"id":   client.UserID,
				"name": client.ID, // In real app, fetch user name from DB
			},
		},
	}

	messageBytes, _ := json.Marshal(event)

	h.Broadcast <- &BroadcastMessage{
		ProjectID: client.ProjectID,
		Message:   messageBytes,
		Exclude:   client,
	}
}

// notifyUserLeft broadcasts user left event
func (h *Hub) notifyUserLeft(client *Client) {
	event := WSEvent{
		Type:      "user:left",
		UserID:    client.UserID,
		ProjectID: client.ProjectID,
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"userId": client.UserID,
		},
	}

	messageBytes, _ := json.Marshal(event)

	h.Broadcast <- &BroadcastMessage{
		ProjectID: client.ProjectID,
		Message:   messageBytes,
	}
}

// readPump handles incoming messages from client
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS] Error reading from client %s: %v", c.ID, err)
			}
			break
		}

		// Parse and handle message
		c.handleMessage(message)
	}
}

// writePump handles outgoing messages to client
func (c *Client) writePump() {
	ticker := time.NewTicker(50 * time.Second)
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

			// Add queued messages to the current message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

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

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(message []byte) {
	var event WSEvent
	if err := json.Unmarshal(message, &event); err != nil {
		log.Printf("[WS] Error parsing message: %v", err)
		return
	}

	// Validate project ID
	if event.ProjectID != c.ProjectID {
		log.Printf("[WS] Invalid project ID from client %s", c.ID)
		return
	}

	// Handle different event types
	switch event.Type {
	case "cursor:move":
		c.handleCursorMove(event.Data)

	case "user:typing":
		c.handleUserTyping(event.Data)

	case "task:update", "task:create", "task:delete":
		// Broadcast task changes to other clients
		c.Hub.mu.RLock()
		room := c.Hub.Projects[c.ProjectID]
		c.Hub.mu.RUnlock()

		if room != nil {
			// Update timestamp and user info
			event.Timestamp = time.Now().Format(time.RFC3339)
			event.UserID = c.UserID

			messageBytes, _ := json.Marshal(event)

			c.Hub.Broadcast <- &BroadcastMessage{
				ProjectID: c.ProjectID,
				Message:   messageBytes,
				Exclude:   c, // Don't echo back to sender
			}
		}

	default:
		log.Printf("[WS] Unknown event type: %s", event.Type)
	}
}

// handleCursorMove updates cursor position
func (c *Client) handleCursorMove(data map[string]interface{}) {
	c.Hub.mu.RLock()
	room := c.Hub.Projects[c.ProjectID]
	c.Hub.mu.RUnlock()

	if room == nil {
		return
	}

	var cursor CursorPosition
	if x, ok := data["x"].(float64); ok {
		cursor.X = x
	}
	if y, ok := data["y"].(float64); ok {
		cursor.Y = y
	}
	if taskID, ok := data["task_id"].(float64); ok && taskID > 0 {
		taskIDUint := uint(taskID)
		cursor.TaskID = &taskIDUint
	}

	room.mu.Lock()
	room.Cursors[c.ID] = &cursor
	room.mu.Unlock()

	// Broadcast to others
	event := WSEvent{
		Type:      "cursor:update",
		UserID:    c.UserID,
		ProjectID: c.ProjectID,
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"userId": c.UserID,
			"cursor": cursor,
		},
	}

	messageBytes, _ := json.Marshal(event)

	c.Hub.Broadcast <- &BroadcastMessage{
		ProjectID: c.ProjectID,
		Message:   messageBytes,
		Exclude:   c,
	}
}

// handleUserTyping handles typing indicators
func (c *Client) handleUserTyping(data map[string]interface{}) {
	isTyping := false
	if val, ok := data["is_typing"].(bool); ok {
		isTyping = val
	}

	var taskID *uint
	if val, ok := data["task_id"].(float64); ok && val > 0 {
		taskIDUint := uint(val)
		taskID = &taskIDUint
	}

	event := WSEvent{
		Type:      "user:typing",
		UserID:    c.UserID,
		ProjectID: c.ProjectID,
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"userId":   c.UserID,
			"isTyping": isTyping,
			"taskId":   taskID,
		},
	}

	messageBytes, _ := json.Marshal(event)

	c.Hub.Broadcast <- &BroadcastMessage{
		ProjectID: c.ProjectID,
		Message:   messageBytes,
		Exclude:   c,
	}
}

// Global hub instance
var GlobalHub *Hub

// InitWebSocket initializes the WebSocket hub
func InitWebSocket() {
	GlobalHub = NewHub()
	go GlobalHub.Run()
}

// HandleWebSocket handles WebSocket connection requests
func HandleWebSocket(c *gin.Context) {
	projectID := c.Query("projectId")
	userID := c.Query("userId")

	if projectID == "" || userID == "" {
		c.JSON(400, gin.H{"error": "projectId and userId are required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] Error upgrading connection: %v", err)
		return
	}

	client := &Client{
		ID:        userID,
		ProjectID: parseUint(projectID),
		UserID:    parseUint(userID),
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Hub:       GlobalHub,
	}

	// Register client
	GlobalHub.Register <- client

	// Start pumps
	go client.writePump()
	go client.readPump()
}

// parseUint converts string to uint
func parseUint(s string) uint {
	var result uint
	fmt.Sscanf(s, "%d", &result)
	return result
}

// BroadcastTaskChange broadcasts task changes to connected clients
func BroadcastTaskChange(projectID uint, changeType string, taskID uint, changes interface{}) {
	if GlobalHub == nil {
		return
	}

	event := WSEvent{
		Type:      changeType,
		ProjectID: projectID,
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"taskId":  taskID,
			"changes": changes,
		},
	}

	messageBytes, _ := json.Marshal(event)

	GlobalHub.Broadcast <- &BroadcastMessage{
		ProjectID: projectID,
		Message:   messageBytes,
	}
}
