package userhandler

import (
	"log/slog"
	"net/http"

	res "project-go/internal/dto/response"
	"project-go/internal/lib/jwt"
	"project-go/internal/lib/password"
	"project-go/internal/lib/response"
	userservice "project-go/internal/service/user"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	drequest "project-go/internal/dto/request"
)

type loginResponse struct {
	response.Response
	Token string           `json:"token"`
	User  res.ResponseUser `json:"user"`
}

func Login(log *slog.Logger, svc *userservice.Service, jwtService *jwt.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.Login"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req drequest.LoginRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		if req.Password == "" || req.Email == "" {
			log.Error("missing fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("missing fields"))
			return
		}

		foundUser, err := svc.FindUserByEmail(req.Email)
		if err != nil {
			log.Error("email not found", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid email or password"))
			return
		}

		if !password.CheckPasswordHash(req.Password, foundUser.Password) {
			log.Error("invalid password")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid email or password"))
			return
		}

		token, err := jwtService.GenerateJWT(foundUser)
		if err != nil {
			log.Error("failed to generate token", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to generate token"))
			return
		}

		render.JSON(w, r, loginResponse{
			Response: response.OK(),
			Token:    token,
			User: res.ResponseUser{
				ID:    foundUser.ID,
				Email: foundUser.Email,
				Name:  foundUser.Name,
				Role:  foundUser.Role,
			},
		})
	}
}
