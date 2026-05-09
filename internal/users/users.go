package users

import (
	"encoding/json"
	"net/http"

	"ufc-backend/internal/auth"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(auth.UserContextKey).(*auth.Claims)

	json.NewEncoder(w).Encode(map[string]string{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
	})
}
