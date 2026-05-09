package routes

import (
	"net/http"

	"ufc-backend/internal/auth"
	"ufc-backend/internal/users"
)

func RegisterUsersRoutes(
	mux *http.ServeMux,
	handler *users.Handler,
) {

	mux.HandleFunc(
		"/users",
		handler.Create,
	)

	mux.Handle(
		"/users/list",
		auth.AuthMiddleware(
			http.HandlerFunc(handler.List),
		),
	)

	mux.Handle(
		"/users/change-password",
		auth.AuthMiddleware(
			http.HandlerFunc(handler.ChangePassword),
		),
	)

	mux.Handle(
		"/users/role",

		auth.AuthMiddleware(

			auth.RequireRoles(
				auth.RoleAdmin,
			)(
				http.HandlerFunc(
					handler.ChangeRole,
				),
			),
		),
	)
}
