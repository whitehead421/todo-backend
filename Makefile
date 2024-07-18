build:
	go build -o bin/$(shell basename $(PWD)) cmd/api/*.go && ./bin/$(shell basename $(PWD))

