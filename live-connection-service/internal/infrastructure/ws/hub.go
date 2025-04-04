package ws

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type BroadcastMessage struct {
	UserId  string
	Message []byte
}

type Hub struct {
	clients   map[string]*websocket.Conn
	mu        sync.RWMutex
	broadcast chan BroadcastMessage
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[string]*websocket.Conn),
		broadcast: make(chan BroadcastMessage, 100),
	}
}

func (h *Hub) WebSocketHandler(userId string, c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
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

func (h *Hub) Send(userId string, message interface{}) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if conn, ok := h.clients[userId]; ok {
		return conn.WriteJSON(message) == nil
	}
	return false
}
