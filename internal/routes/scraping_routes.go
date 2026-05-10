package routes

import (
	"net/http"

	"ufc-backend/internal/auth"
	"ufc-backend/internal/scraping"
)

func RegisterScrapingRoutes(
	mux *http.ServeMux,
	handler *scraping.Handler,
) {

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

	mux.Handle(
		"/scrape/event",
		auth.RequireRoles(
			"admin",
			"manager",
		)(
			http.HandlerFunc(
				handler.ScrapeSingleEvent,
			),
		),
	)

	mux.Handle(
		"/scrape/tapology/events",
		auth.RequireRoles(
			"admin",
			"manager",
		)(
			http.HandlerFunc(
				handler.ScrapeTapologyUFCEvents,
			),
		),
	)

}
