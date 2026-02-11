package res

type AIResponse struct {
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
	Response  string `json:"response"`
	Timestamp string `json:"timestamp"`
}

type SummaryResponse struct {
	Summary string `json:"summary"`
}
