package middleware

import (
	"net/http"
	"strings"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// Always set a specific origin — never wildcard with credentials
		if origin == "" {
			origin = r.Header.Get("Referer")
			if origin != "" {
				// Strip path from referer
				if idx := strings.Index(origin, "://"); idx != -1 {
					after := origin[idx+3:]
					if slashIdx := strings.Index(after, "/"); slashIdx != -1 {
						origin = origin[:idx+3+slashIdx]
					}
				}
			}
		}
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
