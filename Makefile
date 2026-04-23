POSTGRES_DSN=postgres://maps:maps@localhost:5432/maps?sslmode=disable

run:
	go run ./cmd/api

tidy:
	go mod tidy

build:
	go build ./...

migrate-up:
	migrate -path migrations -database "$(POSTGRES_DSN)" up

migrate-down:
	migrate -path migrations -database "$(POSTGRES_DSN)" down 1

migrate-version:
	migrate -path migrations -database "$(POSTGRES_DSN)" version

db-up:
	docker compose up -d

db-down:
	docker compose down

db-logs:
	docker compose logs -f postgres