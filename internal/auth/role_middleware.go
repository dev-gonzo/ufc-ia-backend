package auth

import (
	"net/http"

	"ufc-backend/internal/shared/httpresponse"
)

func RequireRoles(
	roles ...string,
) func(http.Handler) http.Handler {

	allowedRoles := map[string]bool{}

	for _, role := range roles {
		allowedRoles[role] = true
	}

	return func(
		next http.Handler,
	) http.Handler {

		return AuthMiddleware(

			http.HandlerFunc(func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				claims, ok := r.Context().Value(
					UserContextKey,
				).(*Claims)

				if !ok {

					httpresponse.Error(
						w,
						http.StatusUnauthorized,
						"UNAUTHORIZED",
						"unauthorized",
					)

					return
				}

				if !allowedRoles[claims.Role] {

					httpresponse.Error(
						w,
						http.StatusForbidden,
						"FORBIDDEN",
						"insufficient permissions",
					)

					return
				}

				next.ServeHTTP(
					w,
					r,
				)
			}),
		)
	}
}
