package addtestview

import (
	"log/slog"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	testviewservice "project-go/internal/http-server/service/test-view"
	"project-go/internal/lib/api/response"
	"project-go/internal/lib/auth"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func New(log *slog.Logger, testViewService *testviewservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.test.addTestView.New"
		log = log.With(
			slog.String("op", op),
			slog.String("Req_id", middleware.GetReqID(r.Context())),
		)
		var req req.TestViewReq
		error := render.DecodeJSON(r.Body, &req)
		if error != nil {
			log.Error("failed to decode json")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode json"))
			return
		}
		userId, ok := auth.GetUserID(r)
		if !ok {
			log.Error("user id is null")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("user id is null"))
			return
		}
		if req.TestId == 0 {
			log.Error("test id is null")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("test id is null"))
			return
		}
		err := testViewService.AddTestView(req.TestId, userId)
		if err != nil {
			log.Error("internal server error")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}
		render.JSON(w, r, response.OK())
	}
}
