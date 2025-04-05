package ws

import (
	"net/http"
	"sync"

	"github.com/chakchat/chakchat-backend/live-connection-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type BroadcastMessage struct {
	UserId  string
	Message []byte
}

type Hub struct {
	clients map[string]*websocket.Conn
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*websocket.Conn),
	}
}

func (h *Hub) WebSocketHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimId, ok := auth.GetClaims(c.Request.Context())[auth.ClaimId]
		if !ok {
			restapi.SendUnauthorizedError(c, nil)
			return
		}

		userId := claimId.(string)

		upgrade := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		h.mu.Lock()
		h.clients[userId] = conn
		h.mu.Unlock()

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

func (h *Hub) Send(userId string, message any) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if conn, ok := h.clients[userId]; ok {
		return conn.WriteJSON(message) == nil
	}
	return false
}

func (h *Hub) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"clients": len(h.clients),
		})
	}
}
