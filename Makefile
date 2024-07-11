all: swag run

swag:
	swag init --parseDependency -d ./cmd/api

run:
	go run cmd/api/main.go

