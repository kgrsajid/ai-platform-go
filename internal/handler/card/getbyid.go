package cardhandler

import (
	"log/slog"
	"net/http"
	"strconv"

	res "project-go/internal/dto/response"
	"project-go/internal/lib/response"
	cardservice "project-go/internal/service/card"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type cardByIDResponse struct {
	response.Response
	Data *res.CardHolderDetailResponse `json:"data"`
}

func GetByID(log *slog.Logger, svc *cardservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.card.GetByID"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

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

		cardHolder, err := svc.GetCardsByCardHolderId(uint(id))
		if err != nil {
			log.Error("failed to get card", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		render.JSON(w, r, cardByIDResponse{
			Response: response.OK(),
			Data:     cardHolder,
		})
	}
}
