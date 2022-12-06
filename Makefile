DATE=`date +'%Y/%m/%d %_k:%M:%S'`
COMMIT=`git log -n 1 --pretty=format:"%H"`

.PHONY: build
build:
	go build -o shortener -ldflags "-X main.buildVersion=v1.0.2 -X 'main.buildDate=$(DATE)' -X main.buildCommit=$(COMMIT)" cmd/shortener/main.go
	go build -v ./cmd/staticlint
	./staticlint ./...

.PHONY: test
test:
	#./staticlint ./...
	go test -v ./... -count 1

.PHONY: load
load:
	hey -n 10000 -c 5 -m POST -d 'https://ozon.ru' http://localhost:8080
	hey -n 10000 -c 5 -m GET http://localhost:8080/abdksad_urls
	hey -n 10000 -c 5 -m GET http://localhost:8080/user/urls
	hey -n 10000 -c 5 -m POST -d '{"url": "bench.ru"}' http://localhost:8080/
	hey -n 10000 -c 5 -m POST -d '[{"correlation_id": "228","original_url": "gg.ru"},{"correlation_id": "777","original_url": "marlon.ru"}]' http://localhost:8080/api/shorten/batch

.PHONY: cover
cover:
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func coverage.out



.DEFAULT_GOAL := build
