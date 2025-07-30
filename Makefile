# Development Bot Makefile
# Provides common development tasks

.PHONY: help run test format lint clean validate deps check all

# Default target
.DEFAULT_GOAL := help

# Variables
GO_VERSION := 1.24
OUTPUT_DIR := output
RSS_FILE := $(OUTPUT_DIR)/killarney-development.xml

## help: Show available commands
help:
	@echo "Development Bot - Available Commands:"
	@echo ""
	@echo "  make run        - Execute the development bot"
	@echo "  make test       - Run all tests"
	@echo "  make format     - Format Go code with gofmt"
	@echo "  make lint       - Run go vet for static analysis"
	@echo "  make check      - Run format check, lint, and tests"
	@echo "  make validate   - Validate generated RSS feed"
	@echo "  make deps       - Download and verify dependencies"
	@echo "  make clean      - Clean generated files"
	@echo "  make all        - Run full development workflow"
	@echo "  make help       - Show this help message"
	@echo ""

## run: Execute the development bot
run:
	@echo "🤖 Running development bot..."
	go run main.go
	@echo "✅ Development bot completed"

## test: Run all tests with verbose output
test:
	@echo "🧪 Running tests..."
	go test -v ./...
	@echo "✅ All tests completed"

## format: Format Go code and show what was changed
format:
	@echo "🎨 Formatting Go code..."
	@files=$$(gofmt -s -l .); \
	if [ -n "$$files" ]; then \
		echo "Formatting files: $$files"; \
		gofmt -s -w .; \
		echo "✅ Code formatted"; \
	else \
		echo "✅ Code already properly formatted"; \
	fi

## lint: Run go vet for static analysis
lint:
	@echo "🔍 Running static analysis..."
	go vet ./...
	@echo "✅ Static analysis completed"

## check: Run format check, lint, and tests
check: format lint test
	@echo "✅ All checks completed successfully"

## validate: Validate the generated RSS feed
validate: $(RSS_FILE)
	@echo "📡 Validating RSS feed..."
	@if [ ! -f "$(RSS_FILE)" ]; then \
		echo "❌ RSS feed not found at $(RSS_FILE)"; \
		echo "   Run 'make run' first to generate the feed"; \
		exit 1; \
	fi
	@if command -v xmllint >/dev/null 2>&1; then \
		if xmllint --noout $(RSS_FILE) 2>/dev/null; then \
			echo "✅ RSS feed is valid XML"; \
		else \
			echo "❌ RSS feed is not valid XML"; \
			exit 1; \
		fi; \
	else \
		echo "⚠️  xmllint not found, skipping XML validation"; \
		echo "   Install libxml2-utils to enable XML validation"; \
	fi
	@echo "📊 RSS feed stats:"
	@echo "   File size: $$(du -h $(RSS_FILE) | cut -f1)"
	@echo "   Items: $$(grep -c '<item>' $(RSS_FILE) 2>/dev/null || echo '0')"

## deps: Download and verify dependencies
deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	@echo "🔒 Verifying dependencies..."
	go mod verify
	@echo "✅ Dependencies ready"

## clean: Remove generated files and caches
clean:
	@echo "🧹 Cleaning up..."
	@if [ -d "$(OUTPUT_DIR)" ]; then \
		rm -rf $(OUTPUT_DIR)/*; \
		echo "✅ Cleaned output directory"; \
	fi
	go clean -cache
	go clean -testcache
	@echo "✅ Cleanup completed"

## all: Run the complete development workflow
all: clean deps check run validate
	@echo ""
	@echo "🎉 Complete development workflow finished successfully!"
	@echo ""
	@echo "📡 RSS feed available at: $(RSS_FILE)"
	@echo "🌐 GitHub Pages URL: https://jeffadavidson.github.io/development-bot/killarney-development.xml"

# Ensure RSS file exists for validation
$(RSS_FILE):
	@echo "📡 RSS feed not found, generating..."
	$(MAKE) run 