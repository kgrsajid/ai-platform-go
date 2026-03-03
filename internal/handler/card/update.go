package cardhandler

import (
	"log/slog"
	"net/http"
	"strconv"

	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	cardservice "project-go/internal/service/card"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type cardUpdateResponse struct {
	response.Response
	CardHolder *res.CardHolderDetailResponse `json:"data"`
}

func Update(log *slog.Logger, svc *cardservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.card.Update"
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

		idStr := chi.URLParam(r, "cardId")
		if idStr == "" {
			log.Error("missing card id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("missing card id"))
			return
		}

		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Error("invalid card id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid card id"))
			return
		}

		updatedCardholder, err := svc.UpdateCardHolder(req, uint(id))
		if err != nil {
			log.Error("failed to update card holder", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		render.JSON(w, r, cardUpdateResponse{
			Response:   response.OK(),
			CardHolder: updatedCardholder,
		})
	}
}
