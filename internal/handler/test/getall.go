package testhandler

import (
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/auth"
	"project-go/internal/lib/response"
	testservice "project-go/internal/service/test"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type testListResponse struct {
	response.Response
	Data res.PaginationResponse[res.TestResponse] `json:"data"`
}

func GetAll(log *slog.Logger, svc *testservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.test.GetAll"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		search := r.URL.Query().Get("search")
		difficulty := r.URL.Query().Get("difficulty")
		isPrivateStr := r.URL.Query().Get("isPrivate")
		categoriesParams := r.URL.Query()["categories[]"]

		isPrivate, err := strconv.ParseBool(isPrivateStr)
		if err != nil {
			isPrivate = false
		}

		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("unauthorized: missing user id")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("unauthorized"))
			return
		}

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

		var minQ, maxQ *int
		if v, err := strconv.Atoi(r.URL.Query().Get("minQ")); err == nil {
			minQ = &v
		}
		if v, err := strconv.Atoi(r.URL.Query().Get("maxQ")); err == nil {
			maxQ = &v
		}

		filter := request.TestFilter{
			Limit:      limit,
			Offset:     (page - 1) * limit,
			Search:     &search,
			IsPrivate:  &isPrivate,
			UserId:     userId,
			Difficulty: &difficulty,
			Categories: categories,
			MinQ:       minQ,
			MaxQ:       maxQ,
		}

		tests, total, err := svc.GetAllTest(filter)
		if err != nil {
			log.Error("failed to get tests", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get all tests"))
			return
		}

		testResponses := make([]res.TestResponse, 0, len(tests))
		for _, t := range tests {
			testResponses = append(testResponses, res.ToTestResponse(t))
		}

		render.JSON(w, r, testListResponse{
			Response: response.OK(),
			Data: res.PaginationResponse[res.TestResponse]{
				Data:       testResponses,
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: int64(math.Ceil(float64(total) / float64(limit))),
			},
		})
	}
}
