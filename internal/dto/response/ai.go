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

type TitleResponse struct {
	Title string `json:"title"`
}

type AssignmentEvalResponse struct {
	Score       int      `json:"score"`        // 0-100
	MaxScore    int      `json:"max_score"`    // always 100
	Feedback    string   `json:"feedback"`     // LLM explanation
	Strengths   []string `json:"strengths"`    // what student did well
	Improvements []string `json:"improvements"` // what to improve
	GradeLevel  string   `json:"grade_level"`  // how well it matches grade expectations
}
