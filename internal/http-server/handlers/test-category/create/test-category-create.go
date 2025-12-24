package create

import (
	"log/slog"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	testcategory "project-go/internal/http-server/service/test-category"
	"project-go/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Category res.TestCategory `json:"category"`
}

func New(log *slog.Logger, testCategoryService *testcategory.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.testcategory.testcategory-create.New"
		log = log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)
		var req req.TestCategoryReq
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode to json format", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode to json format"))
			return
		}
		if req.Name == "" {
			log.Error("field is empty", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("field is empty"))
			return
		}
		category, err := testCategoryService.CreateCategory(req)
		if err != nil {
			log.Error("failed to create category", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create category"))
			return
		}
		render.JSON(w, r, Response{
			Response: response.OK(),
			Category: res.TestCategory{
				ID:   category.ID,
				Name: category.Name,
			},
		})
	}
}
