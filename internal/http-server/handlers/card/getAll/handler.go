package getall

import (
	"log/slog"
	"math"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	cardservice "project-go/internal/http-server/service/card"
	"project-go/internal/lib/api/response"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Data res.PaginationResponse[res.CardHolderResponse] `json:"data"`
}

func New(log *slog.Logger, cardService *cardservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.card.getAll.New"
		log = log.With(
			slog.String("op", op),
			slog.String("Req_id", middleware.GetReqID(r.Context())),
		)
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		search := r.URL.Query().Get("search")
		categoriesParams := r.URL.Query()["categories[]"]
		var categories []uint
		for _, p := range categoriesParams {
			id, err := strconv.Atoi(p)
			if err != nil {
				continue
			}
			categories = append(categories, uint(id))
		}
		if page <= 0 {
			page = 1
		}
		if limit <= 0 || limit > 100 {
			limit = 10
		}
		offset := (page - 1) * limit

		cardFilter := &req.CardFilter{
			Limit:      limit,
			Offset:     offset,
			Search:     &search,
			Categories: categories,
		}
		cardHolders, total, err := cardService.GetAll(*cardFilter)
		if err != nil {
			log.Error("internal server error")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}
		data := res.PaginationResponse[res.CardHolderResponse]{
			Data:       cardHolders,
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: int64(math.Ceil(float64(total) / float64(limit))),
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Data:     data,
		})
	}
}
