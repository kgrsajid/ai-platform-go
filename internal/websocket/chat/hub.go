package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID    uint
	SessionID uint
	Conn      *websocket.Conn
}

type Hub struct {
	clients map[uint]map[uint]*Client // userID -> sessionID -> Client
	lock    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint]map[uint]*Client),
	}
}

func (h *Hub) AddClient(userID, sessionID uint, conn *websocket.Conn) {
	if conn == nil {
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.clients[userID]; !ok {
		h.clients[userID] = make(map[uint]*Client)
	}

	h.clients[userID][sessionID] = &Client{
		UserID:    userID,
		SessionID: sessionID,
		Conn:      conn,
	}
}

func (h *Hub) RemoveClient(userID, sessionID uint) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if sessions, ok := h.clients[userID]; ok {
		if client, exists := sessions[sessionID]; exists {
			client.Conn.Close()
			delete(sessions, sessionID)
		}
		if len(sessions) == 0 {
			delete(h.clients, userID)
		}
	}
}

func (h *Hub) SendMessage(userID uint, role, message string) {
	h.lock.RLock()
	sessions, ok := h.clients[userID]
	if !ok || len(sessions) == 0 {
		h.lock.RUnlock()
		return
	}

	clients := make([]*Client, 0, len(sessions))
	for _, c := range sessions {
		clients = append(clients, c)
	}
	h.lock.RUnlock()

	for _, client := range clients {
		if err := client.Conn.WriteJSON(map[string]string{
			"role":    role,
			"message": message,
		}); err != nil {
			h.RemoveClient(userID, client.SessionID)
		}
	}
}

func (h *Hub) Broadcast(role, message string) {
	h.lock.RLock()
	allClients := make([]*Client, 0)
	for _, sessions := range h.clients {
		for _, client := range sessions {
			allClients = append(allClients, client)
		}
	}
	h.lock.RUnlock()

	for _, client := range allClients {
		if err := client.Conn.WriteJSON(map[string]string{
			"role":    role,
			"message": message,
		}); err != nil {
			h.RemoveClient(client.UserID, client.SessionID)
		}
	}
}
