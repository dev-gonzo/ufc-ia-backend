# UFC Backend Project Context (AGENTS.md)

This file contains essential architectural and domain knowledge for AI agents and developers working on this project.

## 1. Project Overview
This is a Go-based backend API for scraping, storing, and serving UFC event and fight data. It uses PostgreSQL for persistence and relies on the `ufcstats.com` website as its primary data source.

## 2. Architecture & Patterns

The project follows a standard layered architecture:
- **Handler (Controllers):** Processes HTTP requests, validates inputs, and calls the Service layer. (e.g., `internal/scraping/handler.go`)
- **Service (Business Logic):** Orchestrates operations, such as coordinating the scraper and the database. (e.g., `internal/scraping/service.go`)
- **Repository (Data Access):** Handles all database interactions using `pgx/v5`. (e.g., `internal/scraping/repository.go`)
- **Scraper Client (`ufcstats`):** Dedicated package for parsing HTML using `goquery`. (e.g., `internal/scraping/ufcstats/parser.go`)

### Key Data Patterns
- **Pointer Fields:** Structs (like `ufcstats.Event`) use pointers for optional fields (e.g., `*time.Time`, `*bool`) combined with `json:",omitempty"`. This prevents zero-values (like "0001-01-01") from polluting JSON responses when metadata fields are missing or not yet populated.
- **Idempotent Upserts:** Database insertions use `ON CONFLICT` clauses (e.g., based on the unique `url` of an event). This allows the scraping endpoints to be called repeatedly without duplicating data, updating existing records instead.
- **Synchronization State:** Tables include a `event_sync` boolean flag to track whether deep processing (like scraping fights for an event) has been completed.

## 3. Standardization Rules

### HTTP Responses
**CRITICAL:** ALL HTTP responses (success or error) MUST use the standardized `http_response` package. 
- **Success:** Use `http_response.JSON(w, status, payload)`
- **Error:** Use `http_response.Error(w, status, code, message)`
*Do not use standard `http.Error` or `json.NewEncoder` directly in handlers or middleware.*

### Import Paths
- Use absolute module paths (e.g., `"ufc-backend/internal/shared/http_response"`).
- **Naming Gotcha:** The shared response utility folder is named `http_response`. Ensure imports match this exact path to avoid compilation errors (`ufc-backend/internal/shared/http_response`).

## 4. Authentication
The API uses JWT-based authentication with Role-Based Access Control (RBAC).
- Core middlewares are located in `internal/auth`.
- Protected routes use `auth.RequireRoles("admin", "manager")` to restrict access.

## 5. Next Steps / Future Scope
The current implementation successfully scrapes and persists high-level **Events**.
The next logical step for the system is to process individual fights within those events:
1. Fetch events where `event_sync = false`.
2. Scrape the event URL to extract individual fight details (fighters, weight class, method, round, time).
3. Persist the fight data and mark the event as `event_sync = true`.
