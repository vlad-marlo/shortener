.PHONY: build
build:
	go build -v ./cmd/shortener

.PHONY: test
test:
	go test -v ./...

.DEFAULT_GOAL := build
