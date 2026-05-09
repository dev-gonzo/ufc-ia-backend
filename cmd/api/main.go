package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "ufc-backend/docs"

	"ufc-backend/internal/auth"
	"ufc-backend/internal/database"
	"ufc-backend/internal/routes"
	"ufc-backend/internal/users"
)

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

	err = http.ListenAndServe(
		address,
		mux,
	)

	if err != nil {
		log.Fatal(err)
	}
}
