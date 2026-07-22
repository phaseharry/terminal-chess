.DEFAULT_GOAL := help

BINARY := bin/terminal-chess

.PHONY: dev build run test test-e2e fmt vet tidy clean help

dev: ## Hot-reload the TUI on .go changes (requires watchexec)
	watchexec -r -c -e go -- "go run main.go"

build: ## Compile the binary into bin/
	go build -o $(BINARY) .

run: ## Build and run the app once
	go run .

test: ## Run unit tests
	go test ./...

test-e2e: ## Run end-to-end tests (build-tagged, ./e2e)
	go test -tags e2e ./e2e/...

fmt: ## Format all Go source
	go fmt ./...

vet: ## Report suspicious constructs
	go vet ./...

tidy: ## Add missing and remove unused modules
	go mod tidy
