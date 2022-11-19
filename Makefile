.PHONY: build
build:
	go build -v ./cmd/shortener

.PHONY: test
test:
	go test -v ./... -count 1

.PHONY: load
load:
	hey -n 10000 -c 5 -m POST -d 'http://ozon.ru' http://localhost:8080
	hey -n 10000 -c 5 -m GET http://localhost:8080/abdksad_urls
	hey -n 10000 -c 5 -m GET http://localhost:8080/user/urls
	hey -n 10000 -c 5 -m POST -d '{"url": "bench.ru"}' http://localhost:8080/
	hey -n 10000 -c 5 -m POST -d '[{"correlation_id": "228","original_url": "gg.ru"},{"correlation_id": "777","original_url": "marlon.ru"}]' http://localhost:8080/api/shorten/batch

.DEFAULT_GOAL := build
