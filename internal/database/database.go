package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Connect() *pgx.Conn {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("error loading .env")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(
		context.Background(),
		databaseURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("postgres connected")

	return conn
}
