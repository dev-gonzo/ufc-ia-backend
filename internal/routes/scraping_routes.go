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
		"/scrape/event-fights",
		auth.RequireRoles(
			"admin",
			"manager",
		)(
			http.HandlerFunc(
				handler.ScrapeEventFights,
			),
		),
	)

	mux.Handle(
		"/scrape/fighter",
		auth.RequireRoles(
			"admin",
			"manager",
		)(
			http.HandlerFunc(
				handler.ScrapeFighter,
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
