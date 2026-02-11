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

type WebsocketHub struct {
	clients map[uint]map[uint]*Client // userID -> sessionID -> Client
	lock    sync.RWMutex
}

// Создает новый хаб
func NewHub() *WebsocketHub {
	return &WebsocketHub{
		clients: make(map[uint]map[uint]*Client),
	}
}

// Добавляет клиента в хаб
func (h *WebsocketHub) AddClient(userID, sessionID uint, conn *websocket.Conn) {
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

// Удаляет клиента по sessionID
func (h *WebsocketHub) RemoveClient(userID, sessionID uint) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if sessions, ok := h.clients[userID]; ok {
		if client, exists := sessions[sessionID]; exists {
			client.Conn.Close() // безопасно закрываем соединение
			delete(sessions, sessionID)
		}

		if len(sessions) == 0 {
			delete(h.clients, userID)
		}
	}
}

// Отправляет сообщение конкретному пользователю всем его сессиям
func (h *WebsocketHub) SendMessage(userID uint, role, message string) {
	h.lock.RLock()
	sessions, ok := h.clients[userID]
	if !ok || len(sessions) == 0 {
		h.lock.RUnlock()
		return
	}

	// Копируем соединения, чтобы не держать блокировку при WriteJSON
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
			// Если ошибка — закрываем соединение и убираем клиента
			h.RemoveClient(userID, client.SessionID)
		}
	}
}

// Отправка сообщения всем пользователям
func (h *WebsocketHub) Broadcast(role, message string) {
	h.lock.RLock()
	allClients := []*Client{}
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
