package websocket

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	chatservice "project-go/internal/service/chat"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type TokenParser interface {
	ParseToken(token string) (uint, error)
}

type Handler struct {
	hub         *Hub
	auth        TokenParser
	chatService *chatservice.Service
}

func NewHandler(hub *Hub, auth TokenParser, chatService *chatservice.Service) *Handler {
	return &Handler{hub: hub, auth: auth, chatService: chatService}
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

	isSummaryActive, _ := strconv.ParseInt(isSummaryStr, 10, 64)

	userID, err := h.auth.ParseToken(token)
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
	h.hub.AddClient(userID, uint(sessionID64), conn)

	defer func() {
		h.hub.RemoveClient(userID, uint(sessionID64))
		conn.Close()
		log.Println("ws client disconnected:", userID)
	}()

	summary := 0
	if isSummaryActive == 1 {
		summary = 1
	}

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

		conn.WriteJSON("bot")
	}
}
