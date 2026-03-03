package testcategoryhandler

import (
	"log/slog"
	"net/http"

	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/response"
	testcategoryservice "project-go/internal/service/testcategory"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type categoryResponse struct {
	response.Response
	Category res.TestCategory `json:"category"`
}

func Create(log *slog.Logger, svc *testcategoryservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.testcategory.Create"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req request.TestCategoryReq
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		if req.Name == "" {
			log.Error("name is required")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("name is required"))
			return
		}

		category, err := svc.CreateCategory(req)
		if err != nil {
			log.Error("failed to create category", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create category"))
			return
		}

		render.JSON(w, r, categoryResponse{
			Response: response.OK(),
			Category: res.TestCategory{
				ID:   category.ID,
				Name: category.Name,
			},
		})
	}
}
