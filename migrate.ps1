$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/ufc?sslmode=disable"

goose -dir migrations postgres $env:DATABASE_URL up