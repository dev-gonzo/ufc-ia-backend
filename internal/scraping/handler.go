package scraping

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"ufc-backend/internal/scraping/tapology"
	"ufc-backend/internal/shared/http_response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ScrapeEvents godoc
//
//	@Summary		Scrape and save UFC events
//	@Description	Get all completed UFC events from ufcstats.com and save to DB
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

	events, err := h.service.ScrapeAndSaveEvents(r.Context())

	if err != nil {

		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			"SCRAPE_EVENTS_FAILED",
			"failed to scrape events",
		)

		return
	}

	httpresponse.Success(
		w,
		http.StatusOK,
		events,
	)
}

// ScrapeSingleEvent godoc
//
//	@Summary		Scrape and save a single UFC event
//	@Description	Get event details from ufcstats.com/event-details/{id} and save to DB
//	@Tags			Scraping
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	query		string	true	"Event ID hash"
//	@Success		200	{object}	ufcstats.Event
//	@Failure		500	{object}	httpresponse.ErrorResponse
//	@Router			/scrape/event [get]
func (h *Handler) ScrapeSingleEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.URL.Query().Get("id")

	if id == "" {
		httpresponse.Error(
			w,
			http.StatusBadRequest,
			"MISSING_ID",
			"The id query parameter is required",
		)
		return
	}

	event, err := h.service.ScrapeAndSaveEventByID(r.Context(), id)

	if err != nil {
		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			"SCRAPE_EVENT_FAILED",
			"failed to scrape event",
		)
		return
	}

	httpresponse.Success(
		w,
		http.StatusOK,
		event,
	)
}

// ScrapeTapologyUFCEvents godoc
//
//	@Summary		Scrape and save UFC events from Tapology
//	@Description	Fetch UFC promotion events list from Tapology via remote scraping browser and save to DB
//	@Tags			Scraping
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	tapology.Event
//	@Failure		500	{object}	httpresponse.ErrorResponse
//	@Router			/scrape/tapology/events [get]
func (h *Handler) ScrapeTapologyUFCEvents(
	w http.ResponseWriter,
	r *http.Request,
) {
	events, err := h.service.ScrapeAndSaveTapologyUFCEvents(r.Context())
	if err != nil {
		log.Printf("tapology: handler failed err=%s", redactTapologyError(err.Error()))

		if errors.Is(err, tapology.ErrMissingScrapingBrowserWS) {
			httpresponse.Error(
				w,
				http.StatusInternalServerError,
				"SCRAPING_BROWSER_WS_URL_MISSING",
				"missing scraping browser configuration",
			)
			return
		}

		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			"SCRAPE_TAPOLOGY_FAILED",
			"failed to scrape tapology events",
		)
		return
	}

	httpresponse.Success(
		w,
		http.StatusOK,
		events,
	)
}

var tapologyCredentialsInURL = regexp.MustCompile(`(?i)(wss?|https?)://[^/@\s]+@`)

func redactTapologyError(message string) string {
	if message == "" {
		return message
	}

	return tapologyCredentialsInURL.ReplaceAllStringFunc(message, func(m string) string {
		idx := strings.Index(m, "://")
		if idx == -1 {
			return m
		}
		scheme := m[:idx+3]
		return scheme + "***@"
	})
}
