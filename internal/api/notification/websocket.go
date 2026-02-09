package notification

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients by user ID
	clients map[uint][]*Client

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	// User ID for this client
	userID uint
}

// NotificationMessage represents a notification message sent over WebSocket
type NotificationMessage struct {
	Type string      `json:"type"`
	Data Notification `json:"data"`
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uint][]*Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = append(h.clients[client.userID], client)
			h.mu.Unlock()
			log.Printf("Client registered for user %d (total clients: %d)", client.userID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userID]; ok {
				// Remove this client from the user's client list
				for i, c := range clients {
					if c == client {
						h.clients[client.userID] = append(clients[:i], clients[i+1:]...)
						break
					}
				}
				// If no more clients for this user, remove the entry
				if len(h.clients[client.userID]) == 0 {
					delete(h.clients, client.userID)
				}
			}
			h.mu.Unlock()
			close(client.send)
			log.Printf("Client unregistered for user %d (total clients: %d)", client.userID, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, clients := range h.clients {
				for _, client := range clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToUser sends a message to all clients for a specific user
func (h *Hub) BroadcastToUser(userID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[userID]; ok {
		for _, client := range clients {
			select {
			case client.send <- message:
			default:
				// Client is stuck, close it
				close(client.send)
			}
		}
	}
}

// BroadcastToAll sends a message to all connected clients
func (h *Hub) BroadcastToAll(message []byte) {
	h.broadcast <- message
}

// GetOnlineUserCount returns the count of unique online users
func (h *Hub) GetOnlineUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

// Upgrader specifies the parameters for upgrading an HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, you should check the origin properly
		return true
	},
}

// ServeWS handles websocket requests from clients
func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID, ok := userIDVal.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			return
		}

		// Create client and register
		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan []byte, 256),
			userID: userID,
		}
		client.hub.register <- client

		// Start pumps
		go client.writePump()
		go client.readPump()

		log.Printf("WebSocket connection established for user %d", userID)
	}
}

// BroadcastNotification broadcasts a notification to a specific user
func (h *Hub) BroadcastNotification(userID uint, notification Notification) {
	message := NotificationMessage{
		Type: "notification",
		Data: notification,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return
	}

	h.BroadcastToUser(userID, data)
}

// BroadcastUnreadCount broadcasts updated unread count to a user
func (h *Hub) BroadcastUnreadCount(userID uint, count int64) {
	message := map[string]interface{}{
		"type":  "unread_count",
		"count": count,
		"userID": userID,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal unread count: %v", err)
		return
	}

	h.BroadcastToUser(userID, data)
}

// Global hub instance
var globalHub *Hub

// InitHub initializes the global WebSocket hub
func InitHub() {
	globalHub = NewHub()
	go globalHub.Run()
	log.Println("WebSocket hub initialized")
}

// GetHub returns the global WebSocket hub
func GetHub() *Hub {
	if globalHub == nil {
		InitHub()
	}
	return globalHub
}

// CreateNotificationWithBroadcast creates a notification and broadcasts it via WebSocket
func CreateNotificationWithBroadcast(db interface{}, userID uint, notificationType, title, content string, data map[string]interface{}) error {
	// Create the notification in database (this will be done by the existing CreateNotification function)
	// This is a placeholder for the actual database operation
	if hub := GetHub(); hub != nil {
		// Broadcast the notification
		notification := Notification{
			UserID:  userID,
			Type:    notificationType,
			Title:   title,
			Content: content,
			IsRead:  false,
		}
		hub.BroadcastNotification(userID, notification)
	}
	return nil
}

// BroadcastToUsers sends a notification to multiple users
func BroadcastToUsers(userIDs []uint, notificationType, title, content string, data map[string]interface{}) {
	hub := GetHub()
	if hub == nil {
		return
	}

	for _, userID := range userIDs {
		notification := Notification{
			UserID:  userID,
			Type:    notificationType,
			Title:   title,
			Content: content,
			IsRead:  false,
		}
		hub.BroadcastNotification(userID, notification)
	}
}

// GetOnlineUsers returns a list of online user IDs
func (h *Hub) GetOnlineUsers() []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]uint, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

// IsUserOnline checks if a user is currently connected via WebSocket
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.clients[userID]
	return ok
}

// DisconnectUser disconnects all WebSocket connections for a specific user
func (h *Hub) DisconnectUser(userID uint) {
	h.mu.RLock()
	clients := h.clients[userID]
	h.mu.RUnlock()

	for _, client := range clients {
		h.unregister <- client
	}
}

// GetStats returns hub statistics
func (h *Hub) GetStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	totalClients := 0
	for _, clients := range h.clients {
		totalClients += len(clients)
	}

	return map[string]interface{}{
		"total_users":  len(h.clients),
		"total_clients": totalClients,
		"online_users": h.GetOnlineUsers(),
	}
}

// DebugHub logs the current state of the hub
func (h *Hub) DebugHub() {
	stats := h.GetStats()
	log.Printf("Hub Stats: Users=%d, Clients=%d",
		stats["total_users"], stats["total_clients"])
}
