.PHONY: help run build test vet fmt docker-up docker-down

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-12s %s\n", $$1, $$2}'

run: ## Run with live reload (air), loading .env
	bash -c 'set -a && source .env && set +a && go tool air .'

build: ## Build the binary into bin/
	go build -o bin/pulse-api .

test: ## Run unit tests with the race detector
	go test -race ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Format the code
	gofmt -w .

docker-up: ## Start the API and Postgres with Docker Compose
	docker compose up -d --build

docker-down: ## Stop the Docker Compose stack
	docker compose down
