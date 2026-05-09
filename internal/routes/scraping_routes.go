package routes

import (
	"net/http"

	"ufc-backend/internal/auth"
	"ufc-backend/internal/scraping"
)

func RegisterScrapingRoutes(
	mux *http.ServeMux,
) {

	handler := scraping.NewHandler()

	mux.Handle(
		"/scrape/events",
		auth.RequireRoles(
			"admin",
			"manager",
		)(
			http.HandlerFunc(
				handler.ScrapeEvents,
			),
		),
	)
}
