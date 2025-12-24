package req

type TestRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Tags        []string          `json:"tags"`
	Categories  []uint            `json:"categories"`
	Difficulty  string            `json:"difficulty"`
	Questions   []QuestionRequest `json:"questions"`
}

type QuestionRequest struct {
	Question string          `json:"question"`
	Options  []OptionRequest `json:"options"`
}

type OptionRequest struct {
	OptionText string `json:"optionText"`
	IsCorrect  bool   `json:"isCorrect"`
}

type TestFilter struct {
	Offset     int
	Limit      int
	Search     *string
	Difficulty *string
	Category   *string
}
