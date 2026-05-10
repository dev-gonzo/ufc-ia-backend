package scraping

const (
	CodeMissingID = "MISSING_ID"
	MsgMissingID  = "The id query parameter is required"

	CodeMissingIDOrURL = "MISSING_ID_OR_URL"
	MsgMissingIDOrURL  = "The id or url query parameter is required"

	CodeEventNotFound = "EVENT_NOT_FOUND"
	MsgEventNotFound  = "event not found"

	CodeScrapeEventsFailed = "SCRAPE_EVENTS_FAILED"
	MsgScrapeEventsFailed  = "failed to scrape events"

	CodeScrapeEventFailed = "SCRAPE_EVENT_FAILED"
	MsgScrapeEventFailed  = "failed to scrape event"

	CodeScrapeEventFightsFailed = "SCRAPE_EVENT_FIGHTS_FAILED"
	MsgScrapeEventFightsFailed  = "failed to scrape event fights"

	CodeMissingFightID = "MISSING_FIGHT_ID"
	MsgMissingFightID  = "The fight id query parameter is required"

	CodeFightNotFound = "FIGHT_NOT_FOUND"
	MsgFightNotFound  = "fight not found"

	CodeScrapeFightDetailsFailed = "SCRAPE_FIGHT_DETAILS_FAILED"
	MsgScrapeFightDetailsFailed  = "failed to scrape fight details"

	CodeMissingFighterIDOrURL = "MISSING_FIGHTER_ID_OR_URL"
	MsgMissingFighterIDOrURL  = "The fighter id or url query parameter is required"

	CodeScrapeFighterFailed = "SCRAPE_FIGHTER_FAILED"
	MsgScrapeFighterFailed  = "failed to scrape fighter"

	CodeScrapingBrowserWSURLMissing = "SCRAPING_BROWSER_WS_URL_MISSING"
	MsgScrapingBrowserWSURLMissing  = "missing scraping browser configuration"

	CodeScrapeTapologyFailed = "SCRAPE_TAPOLOGY_FAILED"
	MsgScrapeTapologyFailed  = "failed to scrape tapology events"
)
