package res

import "project-go/internal/models"

type TestCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ToCategoryResponse(categories []models.Category) []TestCategory {
	var newCategories []TestCategory
	if categories == nil {
		return newCategories
	}
	for _, value := range categories {
		newCategories = append(newCategories, TestCategory{
			ID:   value.ID,
			Name: value.Name,
		})
	}
	return newCategories
}
