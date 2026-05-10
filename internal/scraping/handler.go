package scraping

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
	"ufc-backend/internal/scraping/tapology"
	"ufc-backend/internal/shared/http_response"
	"ufc-backend/internal/shared/logger"
)

type Handler struct {
	service *Service
}

func errorMessage(defaultMessage string, err error) string {
	if logger.DebugEnabled() && err != nil {
		return err.Error()
	}
	return defaultMessage
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
		logger.Errorf("scrape_events_failed err=%s", err.Error())

		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			CodeScrapeEventsFailed,
			errorMessage(MsgScrapeEventsFailed, err),
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
//	@Description	Get event details from ufcstats.com/event-details/{id}, save to DB, then scrape fights and fighters for this event and mark event_sync=true
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
			CodeMissingID,
			MsgMissingID,
		)
		return
	}

	event, err := h.service.ScrapeAndSaveEventByID(r.Context(), id)

	if err != nil {
		logger.Errorf("scrape_event_failed id=%s err=%s", id, err.Error())
		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			CodeScrapeEventFailed,
			errorMessage(MsgScrapeEventFailed, err),
		)
		return
	}

	httpresponse.Success(
		w,
		http.StatusOK,
		event,
	)
}

// ScrapeEventFights godoc
//
//	@Summary		Scrape and save UFC event fights
//	@Description	If id is provided, loads event URL from DB. If url is provided and event does not exist, scrapes event and stores it, then scrapes fights. Upserts fights and fighters.
//	@Tags			Scraping
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	query		string	false	"Event UUID (from events table)"
//	@Param			url	query		string	false	"Event URL (ufcstats.com/event-details/...)"
//	@Success		200	{array}	ufcstats.Fight
//	@Failure		400	{object}	httpresponse.ErrorResponse
//	@Failure		404	{object}	httpresponse.ErrorResponse
//	@Failure		500	{object}	httpresponse.ErrorResponse
//	@Router			/scrape/event-fights [get]
func (h *Handler) ScrapeEventFights(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.URL.Query().Get("id")
	url := r.URL.Query().Get("url")

	if strings.TrimSpace(id) == "" && strings.TrimSpace(url) == "" {
		httpresponse.Error(
			w,
			http.StatusBadRequest,
			CodeMissingIDOrURL,
			MsgMissingIDOrURL,
		)
		return
	}

	fights, err := h.service.ScrapeAndSaveEventFights(r.Context(), id, url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Errorf("scrape_event_fights_event_not_found id=%s url=%s", strings.TrimSpace(id), strings.TrimSpace(url))
			httpresponse.Error(
				w,
				http.StatusNotFound,
				CodeEventNotFound,
				MsgEventNotFound,
			)
			return
		}

		logger.Errorf("scrape_event_fights_failed id=%s url=%s err=%s", strings.TrimSpace(id), strings.TrimSpace(url), err.Error())
		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			CodeScrapeEventFightsFailed,
			errorMessage(MsgScrapeEventFightsFailed, err),
		)
		return
	}

	httpresponse.Success(
		w,
		http.StatusOK,
		fights,
	)
}

// ScrapeFighter godoc
//
//	@Summary		Scrape and save a UFC fighter
//	@Description	If id is provided, builds ufcstats.com/fighter-details/{id}. If url is provided, uses it directly. Upserts fighter data.
//	@Tags			Scraping
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	query		string	false	"Fighter ID hash (fighter-details/{id})"
//	@Param			url	query		string	false	"Fighter URL (ufcstats.com/fighter-details/...)"
//	@Success		200	{object}	ufcstats.Fighter
//	@Failure		400	{object}	httpresponse.ErrorResponse
//	@Failure		500	{object}	httpresponse.ErrorResponse
//	@Router			/scrape/fighter [get]
func (h *Handler) ScrapeFighter(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.URL.Query().Get("id")
	url := r.URL.Query().Get("url")

	if strings.TrimSpace(id) == "" && strings.TrimSpace(url) == "" {
		httpresponse.Error(
			w,
			http.StatusBadRequest,
			CodeMissingFighterIDOrURL,
			MsgMissingFighterIDOrURL,
		)
		return
	}

	fighter, err := h.service.ScrapeAndSaveFighter(r.Context(), id, url)
	if err != nil {
		logger.Errorf("scrape_fighter_failed id=%s url=%s err=%s", strings.TrimSpace(id), strings.TrimSpace(url), err.Error())
		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			CodeScrapeFighterFailed,
			errorMessage(MsgScrapeFighterFailed, err),
		)
		return
	}

	httpresponse.Success(
		w,
		http.StatusOK,
		fighter,
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
		logger.Errorf("tapology_handler_failed err=%s", redactTapologyError(err.Error()))

		if errors.Is(err, tapology.ErrMissingScrapingBrowserWS) {
			httpresponse.Error(
				w,
				http.StatusInternalServerError,
				CodeScrapingBrowserWSURLMissing,
				MsgScrapingBrowserWSURLMissing,
			)
			return
		}

		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			CodeScrapeTapologyFailed,
			errorMessage(MsgScrapeTapologyFailed, err),
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
