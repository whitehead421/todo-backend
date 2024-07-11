all: swag run

swag:
	swag init --parseDependency -g main.go -d ./cmd/api,./internal/handlers

run:
	go run cmd/api/main.go

