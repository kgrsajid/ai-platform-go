package userhandler

import (
	"log/slog"
	"net/http"

	"project-go/internal/dto/request"
	"project-go/internal/lib/response"
	userservice "project-go/internal/service/user"

	"github.com/go-chi/render"
)

func ResetPassword(log *slog.Logger, svc *userservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.ResetPasswordRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		if req.Token == "" || req.NewPassword == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("token and new_password are required"))
			return
		}

		if err := svc.ResetPassword(req.Token, req.NewPassword); err != nil {
			log.Error("reset password failed", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}
