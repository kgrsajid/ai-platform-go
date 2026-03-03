package testcategoryhandler

import (
	"log/slog"
	"net/http"

	res "project-go/internal/dto/response"
	"project-go/internal/lib/response"
	"project-go/internal/models"
	testcategoryservice "project-go/internal/service/testcategory"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type categoryListResponse struct {
	response.Response
	Categories []res.TestCategory `json:"categories"`
}

func GetAll(log *slog.Logger, svc *testcategoryservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.testcategory.GetAll"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		categories, err := svc.GetAllCategories()
		if err != nil {
			log.Error("failed to get categories", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get all categories"))
			return
		}

		render.JSON(w, r, categoryListResponse{
			Response:   response.OK(),
			Categories: toTestCategoryList(categories),
		})
	}
}

func toTestCategoryList(categories []models.Category) []res.TestCategory {
	result := make([]res.TestCategory, 0, len(categories))
	for _, t := range categories {
		result = append(result, res.TestCategory{
			ID:   t.ID,
			Name: t.Name,
		})
	}
	return result
}
