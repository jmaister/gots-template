.PHONY: jmeter

# Project name and executable name
PROJECT_NAME := gots-template
BINARY_NAME := gots
ifeq ($(OS),Windows_NT)
	BINARY_NAME := gots.exe
endif

# Build target
build: api-codegen
	@echo "Building $(PROJECT_NAME)..."
	cd webapp && npm run build
	CGO_ENABLED=0 go build -tags=purego -o $(BINARY_NAME) .

# Run target - starts both the application and taronja gateway
run: stop build
	@echo "Starting $(PROJECT_NAME) and Taronja Gateway..."
	@./$(BINARY_NAME) run & echo $$! > .app.pid
	@tg run --config taronja-gateway.yaml & echo $$! > .gateway.pid
	@echo "Both services started:"
	@echo "  Application: http://localhost:8081 (PID: $$(cat .app.pid))"
	@echo "  Gateway: http://localhost:8080 (PID: $$(cat .gateway.pid))"
	@echo "Press Ctrl+C to stop both services"
	@trap 'make stop' INT TERM; sleep infinity


# Stop running services
stop:
	@echo "Stopping services..."
	@pkill -f "$(BINARY_NAME) run" || true
	@pkill -f "tg run" || true
	@rm -f .app.pid .gateway.pid
	@echo "Services stopped."

# Development target with file watching (requires modd)
dev:
	@echo "Starting development mode with file watching..."
	@echo "Using modd from go.mod tools..."
	go run github.com/cortesi/modd/cmd/modd

# Test target
test:
	@echo "Running tests..."
	go test -cover ./...

# Generate coverage and treemap SVG
cover:
	@echo "Generating coverage report..."
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o coverage.html

# Release targets
release-check:
	@echo "Checking GoReleaser config..."
	goreleaser check

release-local:
	@echo "Building release locally (no publish)..."
	goreleaser release --snapshot --clean

release-docker:
	@echo "Building Docker image locally..."
	goreleaser release --snapshot --clean --skip-publish

setup-goreleaser:
	@echo "Setting up GoReleaser..."
	@if [ -f ./scripts/setup_goreleaser.sh ]; then \
		bash ./scripts/setup_goreleaser.sh; \
	else \
		echo "setup_goreleaser.sh not found!"; \
		exit 1; \
	fi

# Clean target
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -f .app.pid .gateway.pid

# Update dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

api-codegen:
	@echo "Generating OpenAPI code..."
	@go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config api/cfg.yaml api/openapi-spec.yaml
	@echo "Installing npm dependencies..."
	@cd webapp && npm install
	@echo "Generating TypeScript SDK..."
	@cd webapp && npm run generate-api

install: build
	cp $(BINARY_NAME) ~/bin/$(BINARY_NAME)

# Default target
.PHONY: all build build-windows run stop dev test cover clean fmt tidy api-codegen
all: build
