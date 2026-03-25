package userhandler

import (
	"log/slog"
	"net/http"
	"strings"

	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/lib/password"
	"project-go/internal/lib/response"
	"project-go/internal/models"
	userservice "project-go/internal/service/user"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type createResponse struct {
	response.Response
	User res.ResponseUser `json:"user"`
}

func Create(log *slog.Logger, svc *userservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.Create"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req request.CreateUserRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

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
			Grade:    req.Grade,
			School:   req.School,
			Language: req.Language,
		}

		newUser, err := svc.CreateUser(&user)
		if err != nil {
			// Check for duplicate email
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "23505") {
				log.Warn("duplicate email", slog.String("email", req.Email))
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, response.Error("email already exists"))
				return
			}
			log.Error("failed to create user", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create user"))
			return
		}

		render.JSON(w, r, createResponse{
			Response: response.OK(),
			User: res.ResponseUser{
				ID:       newUser.ID,
				Email:    newUser.Email,
				Name:     newUser.Name,
				Role:     newUser.Role,
				Grade:    newUser.Grade,
				School:   newUser.School,
				Avatar:   newUser.Avatar,
				Language: newUser.Language,
			},
		})
	}
}
