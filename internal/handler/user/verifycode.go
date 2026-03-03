package userhandler

import (
	"log/slog"
	"net/http"

	"project-go/internal/dto/request"
	"project-go/internal/lib/response"
	userservice "project-go/internal/service/user"

	"github.com/go-chi/render"
)

type verifyCodeResponse struct {
	response.Response
	Token string `json:"token"`
}

func VerifyCode(log *slog.Logger, svc *userservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.VerifyCodeRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		if req.Email == "" || req.Code == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("email and code are required"))
			return
		}

		token, err := svc.VerifyCode(req.Email, req.Code)
		if err != nil {
			log.Error("verify code failed", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, verifyCodeResponse{
			Response: response.OK(),
			Token:    token,
		})
	}
}
