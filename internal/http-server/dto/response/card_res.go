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
	Tags        pq.StringArray `json:"tags"`
	Categories  []TestCategory `json:"cateogires"`
	Cards       []CardResponse `json:"cards"`
	CreatedAt   time.Time      `json:"createdAt"`
}

type CardHolderResponse struct {
	ID                uint           `json:"id"`
	AuthorID          uint           `json:"authorId"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	Tags              pq.StringArray `json:"tags"`
	Categories        []TestCategory `json:"cateogires"`
	NumberOfQuestions int            `json:"numberOfQuestions"`
	CreatedAt         time.Time      `json:"createdAt"`
}

type CardsWithTotal struct {
	CardHolders []CardHolderResponse
	Total       int
}

type CardResponse struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
