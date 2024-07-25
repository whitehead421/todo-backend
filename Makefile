build:
	go build -o bin/$(shell basename $(PWD)) cmd/api/*.go && ./bin/$(shell basename $(PWD))

test:
	go test -json -v ./internal/... 2>&1 -cover | gotestfmt