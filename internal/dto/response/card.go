package res

import (
	"time"

	"github.com/lib/pq"
)

type CardHolderDetailResponse struct {
	ID          uint           `json:"id"`
	AuthorID    uint           `json:"authorId"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	IsPrivate   bool           `json:"isPrivate"`
	Tags        pq.StringArray `json:"tags"`
	Categories  []TestCategory `json:"categories"`
	Cards       []CardResponse `json:"cards"`
	CreatedAt   time.Time      `json:"createdAt"`
}

type CardHolderResponse struct {
	ID                uint           `json:"id"`
	AuthorID          uint           `json:"authorId"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	IsPrivate         bool           `json:"isPrivate"`
	Tags              pq.StringArray `json:"tags"`
	Categories        []TestCategory `json:"categories"`
	NumberOfQuestions int            `json:"numberOfQuestions"`
	CreatedAt         time.Time      `json:"createdAt"`
}

type CardResponse struct {
	ID       uint   `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type GeneratedCardHolderDetailResponse struct {
	ID          uint           `json:"id"`
	AuthorID    uint           `json:"authorId"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	IsPrivate   bool           `json:"isPrivate"`
	Tags        pq.StringArray `json:"tags"`
	Categories  []uint         `json:"categories"`
	Cards       []CardResponse `json:"cards"`
	CreatedAt   time.Time      `json:"createdAt"`
}
