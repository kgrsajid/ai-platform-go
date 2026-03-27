package request

type AiRequest struct {
	UserID   string `json:"user_id"`
	Message  string `json:"message"`
	Language string `json:"language"`
	Grade    int    `json:"grade,omitempty"`
}

type SummaryRequest struct {
	UserID   string `json:"user_id"`
	Topic    string `json:"topic"`
	Language string `json:"language"`
}

type TitleRequest struct {
	Message  string `json:"message"`
	Language string `json:"language"`
}
