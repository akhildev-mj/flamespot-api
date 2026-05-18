# ====================================================================================
# Development Tools & Scripts
# ====================================================================================

.PHONY: dev fmt build start clean

# Run the app locally in development mode
dev:
	@echo "🚀 Starting Fiber server in dev mode..."
	go run main.go

# Format all files using the official Go tool
fmt:
	@echo "🧹 Formatting Go files..."
	go fmt ./...

# Build a tiny, stripped, production-ready binary for Render (Linux target)
build:
	@echo "🏗️ Building optimized production binary..."
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/server main.go

# Run the locally compiled production binary
start:
	@echo "🔥 Running production server..."
	./bin/server

# Clean up build artifacts
clean:
	@echo "🗑️ Cleaning build directory..."
	rm -rf bin/