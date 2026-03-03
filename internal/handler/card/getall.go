package cardhandler

import (
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/response"
	cardservice "project-go/internal/service/card"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type cardListResponse struct {
	response.Response
	Data res.PaginationResponse[res.CardHolderResponse] `json:"data"`
}

func GetAll(log *slog.Logger, svc *cardservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.card.GetAll"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		search := r.URL.Query().Get("search")
		categoriesParams := r.URL.Query()["categories[]"]

		if page <= 0 {
			page = 1
		}
		if limit <= 0 || limit > 100 {
			limit = 10
		}

		var categories []uint
		for _, p := range categoriesParams {
			id, err := strconv.Atoi(p)
			if err != nil {
				continue
			}
			categories = append(categories, uint(id))
		}

		filter := request.CardFilter{
			Limit:      limit,
			Offset:     (page - 1) * limit,
			Search:     &search,
			Categories: categories,
		}

		cardHolders, total, err := svc.GetAll(filter)
		if err != nil {
			log.Error("failed to get cards", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		render.JSON(w, r, cardListResponse{
			Response: response.OK(),
			Data: res.PaginationResponse[res.CardHolderResponse]{
				Data:       cardHolders,
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: int64(math.Ceil(float64(total) / float64(limit))),
			},
		})
	}
}
