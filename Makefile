up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build --no-cache

reset:
	docker compose down -v
	docker volume prune -f
	docker compose up -d --build

restart:
	docker compose restart

clean:
	docker compose down -v
	docker volume prune -f

migrate-db:


test:
	go test -json -v ./internal/... 2>&1 -cover | gotestfmt

.PHONY: up down build restart test