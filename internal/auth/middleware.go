package auth

import (
	"context"
	"net/http"
	"strings"

	httpresponse "ufc-backend/internal/shared/http_response"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			httpresponse.Error(
				w,
				http.StatusUnauthorized,
				"MISSING_TOKEN",
				"missing token",
			)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := ValidateToken(tokenString)

		if err != nil {
			httpresponse.Error(
				w,
				http.StatusUnauthorized,
				"INVALID_TOKEN",
				"invalid token",
			)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
