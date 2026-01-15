package getbyid

import (
	"log/slog"
	"net/http"
	res "project-go/internal/http-server/dto/response"
	cardservice "project-go/internal/http-server/service/card"
	"project-go/internal/lib/api/response"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Data *res.CardHolderDetailResponse `json:"data"`
}

func New(log *slog.Logger, cardService *cardservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.card.getById.New"
		log = log.With(
			slog.String("op", op),
			slog.String("Req_id", middleware.GetReqID(r.Context())),
		)
		idStr := chi.URLParam(r, "cardId")
		if idStr == "" {
			log.Error("id is null")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("id is null"))
			return
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Error("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid id"))
			return
		}
		cardHolder, err := cardService.GetCardsByCardHolderId(uint(id))
		if err != nil {
			log.Error("invalid server error")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("invaid server error"))
			return
		}
		render.JSON(w, r, Response{
			Response: response.OK(),
			Data:     cardHolder,
		})
	}
}
