generate:
	go generate ./...

build:
	go build -o ./bin/ab_force ./cmd/antibruteforce
	go build -o ./bin/migrate ./cmd/migrate

test:
	go test -race ./internal/...

lint:
	golangci-lint run ./...

run: build
	docker-compose -f deployments/docker-compose.yaml -p abf up

down:
	docker-compose -f deployments/docker-compose.yaml -p abf down

.PHONY: build test lint generate run down
