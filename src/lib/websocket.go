package lib

import (
	"encoding/json"
	"foglio/v2/src/config"
	"foglio/v2/src/models"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn   *websocket.Conn
	send   chan models.Notification
	hub    *Hub
	userID string
}

type Hub struct {
	clients    map[string]map[*Client]bool // userID -> clients map
	broadcast  chan models.Notification
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan models.Notification),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*Client]bool)
			}
			h.clients[client.userID][client] = true
			h.mu.Unlock()
			log.Printf("Client connected for user %s. Total users: %d", client.userID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if userClients, ok := h.clients[client.userID]; ok {
				if _, ok := userClients[client]; ok {
					delete(userClients, client)
					close(client.send)
					if len(userClients) == 0 {
						delete(h.clients, client.userID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client disconnected for user %s. Total users: %d", client.userID, len(h.clients))

		case notification := <-h.broadcast:
			h.mu.RLock()
			if userClients, ok := h.clients[notification.OwnerID.String()]; ok {
				for client := range userClients {
					select {
					case client.send <- notification:
						// Message sent successfully
					default:
						close(client.send)
						delete(userClients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Send notification to specific user
func (h *Hub) SendToUser(userID string, notification models.Notification) {
	notification.OwnerID = uuid.Must(uuid.Parse(userID))
	notification.CreatedAt = time.Now()
	h.broadcast <- notification
}

func (h *Hub) BroadcastToAll(notification models.Notification) {
	h.mu.RLock()
	for userID := range h.clients {
		notification.OwnerID = uuid.Must(uuid.Parse(userID))
		notification.CreatedAt = time.Now()
		for client := range h.clients[userID] {
			select {
			case client.send <- notification:
			default:
				close(client.send)
				delete(h.clients[userID], client)
			}
		}
	}
	h.mu.RUnlock()
}

func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	total := 0
	for _, userClients := range h.clients {
		total += len(userClients)
	}
	return total
}

func (h *Hub) GetUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming messages (e.g., mark as read)
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			if action, ok := msg["action"].(string); ok {
				switch action {
				case "mark_read":
					if notificationID, ok := msg["notification_id"].(string); ok {
						// Handle mark as read logic
						log.Printf("Marking notification %s as read for user %s", notificationID, c.userID)
					}
				}
			}
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for notification := range c.send {
		if err := c.conn.WriteJSON(notification); err != nil {
			log.Printf("Write error: %v", err)
			return
		}
	}
}

type WebSocketHandler struct {
	hub *Hub
}

func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

func (wsh *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not upgrade connection"})
		return
	}

	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get(config.AppConfig.CurrentUserId)
	if !exists {
		userID = "anonymous"
	}

	client := &Client{
		conn:   conn,
		send:   make(chan models.Notification, 256),
		hub:    wsh.hub,
		userID: userID.(string),
	}

	wsh.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (wsh *WebSocketHandler) GetStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"connected_clients": wsh.hub.GetClientCount(),
		"connected_users":   wsh.hub.GetUserCount(),
	})
}

// Send notification to specific user
func (wsh *WebSocketHandler) SendNotification(c *gin.Context) {
	var req struct {
		UserID  string                 `json:"user_id" binding:"required"`
		Type    string                 `json:"type" binding:"required"`
		Title   string                 `json:"title" binding:"required"`
		Message string                 `json:"message" binding:"required"`
		Data    map[string]interface{} `json:"data,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := models.Notification{
		ID:      uuid.Must(uuid.NewUUID()),
		Type:    models.NotificationType(req.Type),
		Title:   req.Title,
		Content: req.Message,
		OwnerID: uuid.Must(uuid.Parse(req.UserID)),
		IsRead:  false,
	}

	wsh.hub.SendToUser(req.UserID, notification)
	c.JSON(http.StatusOK, gin.H{"message": "Notification sent"})
}

func (wsh *WebSocketHandler) Broadcast(c *gin.Context) {
	var req struct {
		Type    string                 `json:"type" binding:"required"`
		Title   string                 `json:"title" binding:"required"`
		Message string                 `json:"message" binding:"required"`
		Data    map[string]interface{} `json:"data,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := models.Notification{
		ID:      uuid.Must(uuid.NewUUID()),
		Type:    models.NotificationType(req.Type),
		Title:   req.Title,
		Content: req.Message,
	}

	wsh.hub.BroadcastToAll(notification)
	c.JSON(http.StatusOK, gin.H{"message": "Broadcast sent to all users"})
}
