.PHONY: build
build:
	go build -v ./cmd/shortener

.PHONY: test
test:
	go test -v ./... -count 1

.DEFAULT_GOAL := build
