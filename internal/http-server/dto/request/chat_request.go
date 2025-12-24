package req

type AddChatRequest struct {
	Message   string `json:"message"`
	SessionId uint   `json:"session_id"`
}

type AddMessageByCreatingSession struct {
	Message string `json:"message"`
}
