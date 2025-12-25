package req

import "time"

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
	Categories []uint
	MinQ       *int
	MaxQ       *int
}

type GetALlTestResultsFilter struct {
	UserID   uint   `json:"userId" validate:"required,gt=0"`
	TestID   uint   `json:"testId" validate:"required,gt=0"`
	Duration string `json:"duration" validate:"required,oneof=day week month"`
}

type TestResultReq struct {
	UserId      uint      `json:"userId"`
	TestId      uint      `json:"testId"`
	Score       int       `json:"score"`
	MaxScore    int       `json:"maxScore"`
	StartedAt   time.Time `json:"startedAt"`
	FinishedAt  time.Time `json:"finishedAt"`
	DurationSec int       `json:"durationSec"`
}
