generate:
	go generate ./...

build:
	go build -o ./bin/ab_force ./cmd/antibruteforce

test:
	go test -race ./internal/...

lint:
	golangci-lint run ./...

up: build
	docker-compose -f deployments/docker-compose.yaml -p ab-force up

down:
	docker-compose -f deployments/docker-compose.yaml -p ab-force down

.PHONY: build test lint generate up down
