package create

import (
	"log/slog"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	cardservice "project-go/internal/http-server/service/card"
	"project-go/internal/lib/api/response"
	"project-go/internal/lib/auth"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	CardHolder *res.CardHolderDetailResponse `json:"data"`
}

func New(log *slog.Logger, cardService *cardservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.card.create.New"
		log = log.With(
			slog.String("op", op),
			slog.String("Req_id", middleware.GetReqID(r.Context())),
		)
		var req req.CardHolderRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("decode is failed")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("decode is failed"))
			return
		}
		if req.Title == "" {
			log.Error("some field is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("some field is empty"))
			return
		}
		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("user id is null")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("user id is null"))
			return
		}
		req.AuthorID = &userId
		createdCardholder, err := cardService.CreateCardHolder(req)
		if err != nil {
			log.Error("internal server error")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		render.JSON(w, r, Response{
			Response:   response.OK(),
			CardHolder: createdCardholder,
		})
	}
}
