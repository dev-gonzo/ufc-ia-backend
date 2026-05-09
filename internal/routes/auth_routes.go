package routes

import (
	"net/http"

	"ufc-backend/internal/auth"
)

func RegisterAuthRoutes(
	mux *http.ServeMux,
	handler *auth.Handler,
) {

	mux.HandleFunc(
		"/login",
		handler.Login,
	)
}
