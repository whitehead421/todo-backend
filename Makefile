up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build --no-cache

restart:
	docker compose restart

clean:
	docker compose down -v
	docker volume prune -f

migrate-db:
	migrate -path ./migrations -database "postgresql://postgres:1234@localhost:5433/todo-db?sslmode=disable" up

reset: clean build up migrate-db

test:
	go test -json -v ./internal/... 2>&1 -cover | gotestfmt

.PHONY: up down build restart test