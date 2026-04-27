.PHONY: run test build docker-build docker-run

run:
	go run ./cmd/api

test:
	go test -race -count=1 ./...

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/api ./cmd/api

docker-build:
	docker build -t app-go:latest .

docker-run:
	docker run --rm -p 8080:8080 app-go:latest
