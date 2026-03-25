package res

import "project-go/internal/models"

type ResponseUser struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Role     models.Role `json:"role"`
	Grade    int         `json:"grade"`
	School   string      `json:"school"`
	Avatar   string      `json:"avatar"`
	Language string      `json:"language"`
}
