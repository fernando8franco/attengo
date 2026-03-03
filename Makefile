DB_URL=./assistance.db

migrate-up:
	goose -dir db/migrations sqlite3 $(DB_URL) up

migrate-down:
	goose -dir db/migrations sqlite3 $(DB_URL) down

generate:
	sqlc generate

run:
	go run ./cmd/attengo/main.go