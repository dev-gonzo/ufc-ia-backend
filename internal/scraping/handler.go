package scraping

import (
	"net/http"

	"ufc-backend/internal/scraping/ufcstats"
	"ufc-backend/internal/shared/httpresponse"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// ScrapeEvents godoc
//
//	@Summary		Scrape UFC events
//	@Description	Get all completed UFC events from ufcstats.com
//	@Tags			Scraping
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	ufcstats.Event
//	@Failure		500	{object}	httpresponse.ErrorResponse
//	@Router			/scrape/events [get]
func (h *Handler) ScrapeEvents(
	w http.ResponseWriter,
	r *http.Request,
) {

	events, err := ufcstats.ScrapeEvents()

	if err != nil {

		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			"SCRAPE_EVENTS_FAILED",
			err.Error(),
		)

		return
	}

	httpresponse.JSON(
		w,
		http.StatusOK,
		events,
	)
}
