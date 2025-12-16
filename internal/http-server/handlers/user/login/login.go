package login

import (
	"log/slog"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	userservice "project-go/internal/http-server/service/user"
	"project-go/internal/lib/api/response"
	"project-go/internal/lib/jwt"
	"project-go/internal/lib/password"
	"project-go/internal/models"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UserFindByEmail interface {
	FindUserByEmail(email string) (*models.User, error)
}

type Response struct {
	response.Response
	Token string `json:"token"`
}

func New(log *slog.Logger, userservice *userservice.Service, jwtService *jwt.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.login.New"
		log = log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)
		var req req.LoginRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("req", req))
		if req.Password == "" || req.Email == "" {
			log.Error("missing fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("missing fields"))
			return
		}
		email := req.Email
		FoundUser, err := userservice.FindUserByEmail(email)
		if err != nil {
			log.Error("The email is not correct", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("The email is not correct"))
			return
		}
		if !password.CheckPasswordHash(req.Password, FoundUser.Password) {
			log.Error("Invalid password")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "Invalid Password")
			return
		}
		token, err := jwtService.GenerateJWT(FoundUser)
		if err != nil {
			log.Error("failed to generate token", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to generate token"))
			return
		}
		render.JSON(w, r, Response{
			Response: response.OK(),
			Token:    token,
		})
	}
}
