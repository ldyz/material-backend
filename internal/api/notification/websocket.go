package notification

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// Voice handler for processing voice messages
	voiceHandler *VoiceHandler
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	userID      uint
	clientID    string    // 唯一客户端ID (UUID)
	deviceType  string    // 设备类型: web, android, ios
	connectedAt time.Time // 连接时间
}

// generateClientID 生成唯一客户端ID
func generateClientID() string {
	return uuid.New().String()
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

// SetVoiceHandler sets the voice handler for processing voice messages
func (h *Hub) SetVoiceHandler(handler *VoiceHandler) {
	h.voiceHandler = handler
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = append(h.clients[client.userID], client)
			totalClients := h.countTotalClients()
			h.mu.Unlock()
			log.Printf("[WebSocket] Client registered: user=%d, clientID=%s, device=%s (users=%d, clients=%d)",
				client.userID, client.clientID, client.deviceType, len(h.clients), totalClients)

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
			totalClients := h.countTotalClients()
			h.mu.Unlock()
			close(client.send)
			log.Printf("[WebSocket] Client unregistered: user=%d, clientID=%s, device=%s (users=%d, clients=%d)",
				client.userID, client.clientID, client.deviceType, len(h.clients), totalClients)

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, clients := range h.clients {
				for _, client := range clients {
					select {
					case client.send <- message:
						// 发送成功
					default:
						// 缓冲区满，跳过但不要关闭channel
						log.Printf("[WebSocket] Client buffer full for user %d device %s, skipping broadcast",
							client.userID, client.deviceType)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// countTotalClients 统计总客户端数量（调用时需持有锁）
func (h *Hub) countTotalClients() int {
	total := 0
	for _, clients := range h.clients {
		total += len(clients)
	}
	return total
}

// BroadcastToUser sends a message to all clients for a specific user
func (h *Hub) BroadcastToUser(userID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[userID]; ok {
		for _, client := range clients {
			select {
			case client.send <- message:
				// Message sent successfully
			default:
				// Client buffer is full, log but don't close the channel
				// The writePump will handle the stuck client
				log.Printf("[WebSocket] Client buffer full for user %d, skipping message", userID)
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
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle client messages
		if len(message) > 0 {
			log.Printf("Received message: %s", string(message))
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err == nil {
				log.Printf("Parsed message: %+v", msg)
				if msgType, ok := msg["type"].(string); ok {
					c.handleClientMessage(msgType, msg)
				}
			} else {
				log.Printf("Failed to parse message: %v", err)
			}
		}
	}
}

// handleClientMessage routes client messages to appropriate handlers
func (c *Client) handleClientMessage(msgType string, msg map[string]interface{}) {
	log.Printf("[WebSocket] handleClientMessage: type=%s, userID=%d", msgType, c.userID)
	switch msgType {
	case "ping":
		c.sendPong()
	case "voice":
		log.Printf("[WebSocket] Routing to voice handler")
		go c.handleVoiceMessage(msg)
	case "voice_stream_chunk":
		// 流式语音片段
		go c.handleVoiceStreamChunk(msg)
	case "voice_stream_end":
		// 结束流式语音
		go c.handleVoiceStreamEnd(msg)
	case "chat":
		log.Printf("[WebSocket] Routing to chat handler")
		go c.handleChatMessage(msg)
	default:
		log.Printf("Unknown message type: %s", msgType)
	}
}

// sendPong sends a pong response
func (c *Client) sendPong() {
	pongMsg := map[string]interface{}{"type": "pong"}
	if data, err := json.Marshal(pongMsg); err == nil {
		c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Failed to send pong: %v", err)
		}
	}
}

// handleVoiceMessage handles voice messages from the client
func (c *Client) handleVoiceMessage(msg map[string]interface{}) {
	if c.hub.voiceHandler == nil {
		c.sendError("语音处理服务未配置")
		return
	}

	log.Printf("[WebSocket] handleVoiceMessage: user=%d, starting processing", c.userID)

	// Get audio data
	audioDataBase64, ok := msg["data"].(string)
	if !ok {
		c.sendError("音频数据格式错误")
		return
	}

	mimeType, _ := msg["mimeType"].(string)
	if mimeType == "" {
		mimeType = "audio/webm"
	}

	// Get conversation history if provided
	var history []map[string]interface{}
	if h, ok := msg["history"].([]interface{}); ok {
		history = make([]map[string]interface{}, 0, len(h))
		for _, item := range h {
			if m, ok := item.(map[string]interface{}); ok {
				history = append(history, m)
			}
		}
	}

	// Decode base64 audio data
	audioData, err := base64.StdEncoding.DecodeString(audioDataBase64)
	if err != nil {
		c.sendError("音频数据解码失败: " + err.Error())
		return
	}

	log.Printf("[WebSocket] handleVoiceMessage: user=%d, audio size=%d bytes, history=%d", c.userID, len(audioData), len(history))

	// Process voice message in a separate goroutine with panic recovery
	// Context is created inside goroutine to prevent premature cancellation
	go func() {
		// Create context for this operation with longer timeout for ASR processing
		// ASR service may take up to 120 seconds for long audio
		ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()

		defer func() {
			if r := recover(); r != nil {
				log.Printf("[WebSocket] Panic in voice handler: %v", r)
				c.sendError("语音处理发生错误")
			}
		}()
		c.hub.voiceHandler.HandleVoiceMessageWithHistory(ctx, c, audioData, mimeType, int(c.userID), history)
	}()
}

// handleChatMessage handles chat messages from the client
func (c *Client) handleChatMessage(msg map[string]interface{}) {
	if c.hub.voiceHandler == nil {
		c.sendError("AI 服务未配置")
		return
	}

	log.Printf("[WebSocket] handleChatMessage: user=%d, starting processing", c.userID)

	// Get message
	message, ok := msg["message"].(string)
	if !ok || message == "" {
		c.sendError("消息内容不能为空")
		return
	}

	// Get conversation history if provided
	var history []map[string]interface{}
	if h, ok := msg["history"].([]interface{}); ok {
		history = make([]map[string]interface{}, 0, len(h))
		for _, item := range h {
			if m, ok := item.(map[string]interface{}); ok {
				history = append(history, m)
			}
		}
	}

	// Process chat message in a separate goroutine with panic recovery
	// Context is created inside goroutine to prevent premature cancellation
	go func() {
		// Create context for this operation
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		defer func() {
			if r := recover(); r != nil {
				log.Printf("[WebSocket] Panic in chat handler: %v", r)
				c.sendError("聊天处理发生错误")
			}
		}()
		c.hub.voiceHandler.HandleChatMessage(ctx, c, message, history, int(c.userID))
	}()
}

// handleVoiceStreamChunk handles streaming voice chunk from the client
func (c *Client) handleVoiceStreamChunk(msg map[string]interface{}) {
	if c.hub.voiceHandler == nil {
		c.sendError("语音处理服务未配置")
		return
	}

	// Get audio data
	audioDataBase64, ok := msg["data"].(string)
	if !ok {
		return // 忽略无效数据
	}

	mimeType, _ := msg["mimeType"].(string)
	if mimeType == "" {
		mimeType = "audio/webm"
	}

	// Decode base64 audio data
	audioData, err := base64.StdEncoding.DecodeString(audioDataBase64)
	if err != nil {
		return // 忽略解码错误
	}

	// Process in goroutine
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		c.hub.voiceHandler.HandleVoiceStreamChunk(ctx, c, audioData, mimeType, int(c.userID))
	}()
}

// handleVoiceStreamEnd handles end of streaming voice from the client
func (c *Client) handleVoiceStreamEnd(msg map[string]interface{}) {
	if c.hub.voiceHandler == nil {
		c.sendError("语音处理服务未配置")
		return
	}

	// Process in goroutine
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		c.hub.voiceHandler.HandleVoiceStreamEnd(ctx, c, int(c.userID))
	}()
}

// sendError sends an error message to the client
func (c *Client) sendError(message string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[WebSocket] Recovered from panic when sending error to user %d: %v", c.userID, r)
		}
	}()

	errMsg := map[string]interface{}{
		"type":    "error",
		"message": message,
	}
	if data, err := json.Marshal(errMsg); err == nil {
		select {
		case c.send <- data:
		default:
			log.Printf("Client send buffer full, cannot send error")
		}
	}
}

// Send sends a message to the client
func (c *Client) Send(msg map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[WebSocket] Recovered from panic when sending to user %d: %v", c.userID, r)
		}
	}()

	if data, err := json.Marshal(msg); err == nil {
		log.Printf("[WebSocket] Sending to user %d: type=%s", c.userID, msg["type"])
		select {
		case c.send <- data:
			log.Printf("[WebSocket] Message sent to channel for user %d", c.userID)
		default:
			log.Printf("[WebSocket] Client send buffer full for user %d, cannot send message type=%s", c.userID, msg["type"])
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

		// 获取设备类型（优先从 Query 参数，其次从 Header）
		deviceType := c.Query("device_type")
		if deviceType == "" {
			deviceType = c.GetHeader("X-Device-Type")
		}
		if deviceType == "" {
			deviceType = "web" // 默认
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			return
		}

		// Create client and register
		client := &Client{
			hub:         hub,
			conn:        conn,
			send:        make(chan []byte, 1024), // Increased buffer for streaming responses
			userID:      userID,
			clientID:    generateClientID(),
			deviceType:  deviceType,
			connectedAt: time.Now(),
		}
		client.hub.register <- client

		// Start pumps
		go client.writePump()
		go client.readPump()

		log.Printf("[WebSocket] Connection established: user=%d, clientID=%s, device=%s", userID, client.clientID, deviceType)
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
	deviceStats := make(map[string]int)
	for _, clients := range h.clients {
		totalClients += len(clients)
		for _, c := range clients {
			deviceStats[c.deviceType]++
		}
	}

	return map[string]interface{}{
		"total_users":   len(h.clients),
		"total_clients": totalClients,
		"by_device":     deviceStats,
		"online_users":  h.GetOnlineUsers(),
	}
}

// DebugHub logs the current state of the hub
func (h *Hub) DebugHub() {
	stats := h.GetStats()
	log.Printf("Hub Stats: Users=%d, Clients=%d",
		stats["total_users"], stats["total_clients"])
}
