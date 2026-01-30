.PHONY: build test test-unit test-e2e clean install help

BINARY_NAME=basecamp
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

help:
	@echo "Basecamp CLI"
	@echo ""
	@echo "Usage:"
	@echo "  make build       Build the CLI"
	@echo "  make test-unit   Run unit tests"
	@echo "  make test-e2e    Run e2e tests (requires credentials)"
	@echo "  make test        Run unit tests (alias)"
	@echo "  make clean       Remove build artifacts"
	@echo "  make install     Install to GOPATH/bin"
	@echo ""
	@echo "Environment variables for e2e tests:"
	@echo "  BASECAMP_TEST_PROJECT_ID   Project ID to test against"
	@echo "  BASECAMP_TEST_BOARD_ID     Board ID for card tests"
	@echo "  BASECAMP_TEST_CARD_ID      Card ID for detail/move tests"
	@echo "  BASECAMP_TEST_TODOLIST_ID  Todolist ID for todo tests"
	@echo ""
	@echo "Example:"
	@echo "  export BASECAMP_TEST_PROJECT_ID=12345678"
	@echo "  export BASECAMP_TEST_BOARD_ID=87654321"
	@echo "  export BASECAMP_TEST_CARD_ID=44444444"
	@echo "  export BASECAMP_TEST_TODOLIST_ID=99999999"
	@echo "  make test-e2e"

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) ./cmd/basecamp

# Run unit tests (no API required)
test-unit:
	go test -v ./internal/...

# Run e2e tests (requires API credentials and config)
test-e2e: build
	@if [ -z "$$BASECAMP_TEST_PROJECT_ID" ]; then \
		echo "Error: BASECAMP_TEST_PROJECT_ID not set"; \
		echo "Run 'make help' for required environment variables"; \
		exit 1; \
	fi
	BASECAMP_TEST_BINARY=$(CURDIR)/$(BINARY_NAME) go test -v ./e2e/tests/...

# Default test target runs unit tests
test: test-unit

clean:
	rm -f $(BINARY_NAME)

install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/
