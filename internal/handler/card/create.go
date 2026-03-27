package cardhandler

import (
	"log/slog"
	"net/http"

	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	cardservice "project-go/internal/service/card"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type cardDetailResponse struct {
	response.Response
	CardHolder *res.CardHolderDetailResponse `json:"data"`
}

func Create(log *slog.Logger, svc *cardservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.card.Create"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req request.CardHolderRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		if req.Title == "" {
			log.Error("title is required")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("title is required"))
			return
		}

		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("unauthorized: missing user id")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}
		req.AuthorID = &userId

		createdCardholder, err := svc.CreateCardHolder(req)
		if err != nil {
			log.Error("failed to create card holder", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create card set: "+err.Error()))
			return
		}

		render.JSON(w, r, cardDetailResponse{
			Response:   response.OK(),
			CardHolder: createdCardholder,
		})
	}
}
