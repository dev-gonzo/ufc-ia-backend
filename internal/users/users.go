package users

import (
	"net/http"

	"ufc-backend/internal/auth"
	"ufc-backend/internal/shared/http_response"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserContextKey).(*auth.Claims)

	if !ok {
		httpresponse.Error(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"unauthorized",
		)
		return
	}

	httpresponse.Success(
		w,
		http.StatusOK,
		ProfileResponse{
			UserID: claims.UserID,
			Email:  claims.Email,
			Role:   claims.Role,
		},
	)
}
