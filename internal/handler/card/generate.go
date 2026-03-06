package cardhandler

import (
	"log/slog"
	"net/http"
	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/api/response"
	cardservice "project-go/internal/service/card"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Cards res.GeneratedCardHolderDetailResponse `json:"data"`
}

func GenerateCard(log *slog.Logger, cardservice *cardservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.card.generate.GenerateCard"
		log = log.With(
			slog.String("op", op),
			slog.String("ReqId", middleware.GetReqID(r.Context())),
		)
		var req request.GenerateCardReq
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}
		if req.NumberOfCards == 0 || req.Description == "" || len(req.Categories) == 0 {
			log.Error("missing required fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		resp, err := cardservice.GenerateCards(r.Context(), req)
		if err != nil {
			log.Error("failed to generate quiz", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error server"))
			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Cards:    *resp,
		})
	}
}
