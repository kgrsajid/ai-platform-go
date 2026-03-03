package request

import "github.com/lib/pq"

type CardHolderRequest struct {
	Title       string         `json:"title"`
	Description *string        `json:"description"`
	Tags        pq.StringArray `json:"tags"`
	Categories  []uint         `json:"categories"`
	AuthorID    *uint
	ID          *uint
	Cards       []CardRequest `json:"cards"`
}

type CardRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type CardFilter struct {
	Offset     int
	Limit      int
	Search     *string
	Categories []uint
}
