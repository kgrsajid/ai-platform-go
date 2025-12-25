package test_get_all

import (
	"log/slog"
	"math"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	testservice "project-go/internal/http-server/service/test"
	"project-go/internal/lib/api/response"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
	Data res.PaginationResponse[res.TestResponse] `json:"data"`
}

func New(log *slog.Logger, testService *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.test.test.New"
		log = log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		search := r.URL.Query().Get("search")
		categoriesParams := r.URL.Query()["categories[]"]
		minQStr := r.URL.Query().Get("minQ")
		maxQStr := r.URL.Query().Get("maxQ")
		var minQ, maxQ *int
		if v, err := strconv.Atoi(minQStr); err == nil {
			minQ = &v
		}
		if v, err := strconv.Atoi(maxQStr); err == nil {
			maxQ = &v
		}
		var categories []uint
		for _, p := range categoriesParams {
			id, err := strconv.Atoi(p)
			if err != nil {
				continue
			}
			categories = append(categories, uint(id))
		}
		difficulty := r.URL.Query().Get("difficulty")
		if page <= 0 {
			page = 1
		}
		if limit <= 0 || limit > 100 {
			limit = 10
		}
		offset := (page - 1) * limit

		testFilter := &req.TestFilter{
			Limit:      limit,
			Offset:     offset,
			Search:     &search,
			Difficulty: &difficulty,
			Categories: categories,
			MinQ:       minQ,
			MaxQ:       maxQ,
		}

		tests, total, err := testService.GetAllTest(*testFilter)

		if err != nil {
			log.Error("failed to get all tests", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get all tests"))
			return
		}

		data := res.PaginationResponse[res.TestResponse]{
			Data:       ToTestResponseList(tests),
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

func ToTestResponseList(tests []res.TestWithQuestionsCount) []res.TestResponse {
	result := make([]res.TestResponse, 0, len(tests))

	for _, t := range tests {
		result = append(result, res.ToTestResponse(t))
	}

	return result
}
