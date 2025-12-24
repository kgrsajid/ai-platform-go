package getall

import (
	"log/slog"
	"net/http"
	res "project-go/internal/http-server/dto/response"
	testcategory "project-go/internal/http-server/service/test-category"
	"project-go/internal/lib/api/response"
	"project-go/internal/models"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Categories []res.TestCategory `json:"categories"`
}

func New(log *slog.Logger, testCategoryService testcategory.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.test-category.getAll.New"
		log = log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)
		categories, err := testCategoryService.GetAllCategories()
		if err != nil {
			log.Error("failed to get all categories", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get all categories"))
			return
		}
		render.JSON(w, r, Response{
			Response:   response.OK(),
			Categories: ToTestCategoryList(categories),
		})
	}
}

func ToTestCategoryList(categories []models.Category) []res.TestCategory {
	result := make([]res.TestCategory, 0, len(categories))
	for _, t := range categories {
		result = append(result, res.TestCategory{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	return result
}
