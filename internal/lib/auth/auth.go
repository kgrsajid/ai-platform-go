package auth

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

const (
	CtxUserID = "user_id"
	CtxRole   = "role"
)

type Service struct {
	secret string
}

func GetUserID(r *http.Request) (uint, bool) {
	val := r.Context().Value(CtxUserID)
	if val == nil {
		return 0, false
	}

	id, ok := val.(uint)
	if !ok {
		return 0, false
	}

	return id, true
}

func GetRole(r *http.Request) (string, bool) {
	val := r.Context().Value(CtxRole)
	if val == nil {
		return "", false
	}
	role, ok := val.(string)
	return role, ok
}

func New(secret string) *Service {
	return &Service{secret: secret}
}

func (s *Service) ParseToken(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found")
	}

	return uint(userID), nil
}
