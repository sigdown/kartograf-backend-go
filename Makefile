run:
	go run ./cmd/api

fmt:
	go fmt ./...

tidy:
	go mod tidy

up:
	docker compose up -d

down:
	docker compose down

migrate-up:
	migrate -path migrations -database "postgres://maps:maps@localhost:5432/maps?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://maps:maps@localhost:5432/maps?sslmode=disable" down 1