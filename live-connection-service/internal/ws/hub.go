package ws

import (
	"net/http"
	"sync"
	"time"

	"github.com/chakchat/chakchat-backend/live-connection-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	lastPing time.Time
}

type BroadcastMessage struct {
	UserId  uuid.UUID
	Message any
}

type Hub struct {
	clients    map[uuid.UUID]*Client
	mu         sync.RWMutex
	broadcast  chan BroadcastMessage
	pingTicker *time.Ticker
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		broadcast:  make(chan BroadcastMessage, 100),
		pingTicker: time.NewTicker(5 * time.Second),
	}
}

func (h *Hub) WebSocketHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimId, ok := auth.GetClaims(c.Request.Context())[auth.ClaimId]
		if !ok {
			restapi.SendUnauthorizedError(c, nil)
			return
		}

		id := claimId.(string)

		userId, err := uuid.Parse(id)
		if err != nil {
			restapi.SendUnauthorizedError(c, nil)
			return
		}

		upgrade := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := &Client{
			conn:     conn,
			lastPing: time.Now(),
		}

		conn.SetPongHandler(func(string) error {
			h.mu.Lock()
			client.lastPing = time.Now()
			h.mu.Unlock()
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			return nil
		})

		h.mu.Lock()
		defer h.mu.Unlock()
		h.clients[userId] = client

		defer func() {
			h.mu.Lock()
			delete(h.clients, userId)
			h.mu.Unlock()
			conn.Close()
		}()

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}
}

func (h *Hub) Send(userId uuid.UUID, message any) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.clients[userId]; ok {
		client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		return client.conn.WriteJSON(message) == nil
	}
	return false
}

func (h *Hub) GetOnlineStatus(userIds []uuid.UUID) map[uuid.UUID]bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	status := make(map[uuid.UUID]bool)
	now := time.Now()

	for _, userId := range userIds {
		if client, ok := h.clients[userId]; ok {
			status[userId] = now.Sub(client.lastPing) < 10*time.Second
		} else {
			status[userId] = false
		}
	}
	return status
}

func (h *Hub) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"clients": len(h.clients),
		})
	}
}
