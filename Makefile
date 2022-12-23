.PHONY: run

run:build start

start:
	./bin/hasher

build:
	go build -o bin/hasher -v cmd/cli/main.go

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := run