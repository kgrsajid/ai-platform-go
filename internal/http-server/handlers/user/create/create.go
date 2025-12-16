package userCreate

import (
	"log/slog"
	"net/http"
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	userservice "project-go/internal/http-server/service/user"
	"project-go/internal/lib/api/response"
	"project-go/internal/lib/password"
	"project-go/internal/models"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UserCreate interface {
	CreateUser(u *models.User) (*models.User, error)
}

type Response struct {
	response.Response
	User res.ResponseUser
}

func New(log *slog.Logger, userservice *userservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/user/create/New"
		log = log.With(
			slog.String("op", op),
			slog.String("requiest_id", middleware.GetReqID(r.Context())),
		)
		var req req.CreateUserRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalid request body"))
			return
		}
		log.Info("request body decoded", slog.Any("req body:", req))

		if req.Email == "" || req.Name == "" || req.Password == "" {
			log.Error("missing fields")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("missing fields"))
			return
		}

		if !models.IsValidRole(req.Role) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid role"))
			return
		}

		hash, err := password.HashPassword(req.Password)
		if err != nil {
			log.Error("failed to hash password", slog.Any("err", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to hash password"))
			return
		}
		user := models.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: hash,
			Role:     models.Role(req.Role),
		}

		newUser, err := userservice.CreateUser(&user)
		if err != nil {
			log.Error("failed to create user", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create user"))
			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			User: res.ResponseUser{
				ID:    newUser.ID,
				Email: newUser.Email,
				Name:  newUser.Name,
				Role:  models.Role(newUser.Role),
			},
		})
	}
}
