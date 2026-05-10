package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "ufc-backend/docs"

	"ufc-backend/internal/auth"
	"ufc-backend/internal/database"
	"ufc-backend/internal/routes"
	"ufc-backend/internal/scraping"
	"ufc-backend/internal/shared/http_response"
	"ufc-backend/internal/users"
)

func tryLoadScrapingBrowserEnvFromInfoTxt() {
	if strings.TrimSpace(os.Getenv("SCRAPING_BROWSER_WS_URL")) != "" {
		return
	}

	file, err := os.Open("info.txt")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "wss://") && strings.Contains(line, ":9222") {
			_ = os.Setenv("SCRAPING_BROWSER_WS_URL", line)
			return
		}
	}
}

// @title UFC Backend API
// @version 1.0
// @description UFC scraping and AI platform
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env not found")
	}

	tryLoadScrapingBrowserEnvFromInfoTxt()

	db := database.Connect()

	mux := http.NewServeMux()

	usersRepository := users.NewRepository(
		db,
	)

	authService := auth.NewService(
		usersRepository,
	)

	authHandler := auth.NewHandler(
		authService,
	)

	usersService := users.NewService(
		usersRepository,
	)

	usersHandler := users.NewHandler(
		usersService,
		usersRepository,
	)

	scrapingRepository := scraping.NewRepository(
		db,
	)

	scrapingService := scraping.NewService(
		scrapingRepository,
	)

	scrapingHandler := scraping.NewHandler(
		scrapingService,
	)

	routes.RegisterAuthRoutes(
		mux,
		authHandler,
	)

	routes.RegisterUsersRoutes(
		mux,
		usersHandler,
	)

	routes.RegisterScrapingRoutes(
		mux,
		scrapingHandler,
	)

	mux.Handle(
		"/swagger/",
		httpSwagger.Handler(),
	)

	port := os.Getenv(
		"SERVER_PORT",
	)

	if port == "" {
		port = "8080"
	}

	address := ":" + port

	log.Printf(
		"server running on %s",
		address,
	)

	log.Printf(
		"swagger running on http://localhost%s/swagger/index.html",
		address,
	)

	handler := httpresponse.RecoverMiddleware(mux)

	err = http.ListenAndServe(
		address,
		handler,
	)

	if err != nil {
		log.Fatal(err)
	}
}
