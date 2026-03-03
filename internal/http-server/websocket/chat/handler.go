package websocket

import (
	"fmt"
	"log"
	"net/http"
	chatservice "project-go/internal/http-server/service/chat"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // ❗ для dev, в проде ограничь
	},
}

type TokenParser interface {
	ParseToken(token string) (uint, error)
}

type Handler struct {
	Hub         *WebsocketHub
	Auth        TokenParser
	chatservice *chatservice.Service
}

func NewHandler(hub *WebsocketHub, auth TokenParser, chatservice *chatservice.Service) *Handler {
	return &Handler{Hub: hub, Auth: auth, chatservice: chatservice}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("WS: handler entered")

	token := r.URL.Query().Get("token")
	sessionIdStr := r.URL.Query().Get("session_id")
	isSummaryStr := r.URL.Query().Get("summary")
	if token == "" {
		http.Error(w, "no token", http.StatusUnauthorized)
		return
	}

	if sessionIdStr == "" {
		http.Error(w, "no session id", http.StatusUnauthorized)
		return
	}

	sessionID64, err := strconv.ParseUint(sessionIdStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid session_id", http.StatusBadRequest)
		return
	}

	isSummaryActive, err := strconv.ParseInt(isSummaryStr, 10, 64)

	userID, err := h.Auth.ParseToken(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}

	log.Println("WS upgraded, user:", userID)

	h.Hub.AddClient(userID, uint(sessionID64), conn)
	log.Println("ws client connected:", userID)

	// 🔑 Гарантированный cleanup
	defer func() {
		h.Hub.RemoveClient(userID, uint(sessionID64))
		conn.Close()
		log.Println("ws client disconnected:", userID)
	}()

	// Цикл чтения сообщений
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WS read error:", err)
			break
		}

		err = conn.WriteJSON("user")
		if err != nil {
			log.Println("ws write error:", err)
		}

		summary := 0
		if isSummaryActive == 1 {
			summary = 1
		}

		fmt.Println(summary)

		botMsg, err := h.chatservice.AddMessage(
			r.Context(),
			userID,
			uint(sessionID64),
			string(msg),
			summary,
		)
		if err != nil {
			log.Println(err)
			// Всё равно уведомляем фронт — он сделает refetch и покажет статус Error + кнопку retry
			conn.WriteJSON("bot")
			continue
		}
		fmt.Print(botMsg)
		conn.WriteJSON("bot")

	}
}
