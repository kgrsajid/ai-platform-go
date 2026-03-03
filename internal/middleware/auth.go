package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"project-go/internal/lib/auth"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing auth header", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			idFloat, ok := claims["user_id"].(float64)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}
			userID := uint(idFloat)

			ctx := context.WithValue(r.Context(), auth.CtxUserID, userID)
			ctx = context.WithValue(ctx, auth.CtxRole, claims["role"])
			slog.Info("saved values in context", slog.Any("context", ctx))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
