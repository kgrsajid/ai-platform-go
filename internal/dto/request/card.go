package request

import "github.com/lib/pq"

type CardHolderRequest struct {
	Title       string         `json:"title"`
	Description *string        `json:"description"`
	IsPrivate   bool           `json:"isPrivate"`
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
	IsPrivate  *bool
	UserId     uint
	Categories []uint
	MinQ       *int
	MaxQ       *int
}

type GenerateCardReq struct {
	Title         string  `json:"title"`
	Description   string  `json:"context"`
	Categories    []uint  `json:"categories"`
	IsPrivate     bool    `json:"is_private"`
	NumberOfCards int     `json:"num_cards"`
	Language      *string `json:"language"`
}
