package testhandler

import (
	"log/slog"
	"net/http"
	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/api/response"
	testservice "project-go/internal/service/test"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Test res.GeneratedTestResponse `json:"data"`
}

func GenerateQuiz(log *slog.Logger, svg *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.test.generate.New"
		log = log.With(
			slog.String("op", op),
			slog.String("Req_id", middleware.GetReqID(r.Context())),
		)
		var req request.GenerateQuizReq
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		if req.Description == "" || req.Difficulty == "" || len(req.Categories) == 0 {
			log.Error("missing required fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("title, difficulty, and categories are required"))
			return
		}

		resp, err := svg.GenerateQuiz(r.Context(), req)
		if err != nil {
			log.Error("failed to generate quiz", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error server"))
			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Test:     *resp,
		})
	}
}
