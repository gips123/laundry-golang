.PHONY: help build build-server build-seed test-build clean run run-seed tidy vet fmt

# Default target
help:
	@echo "Available commands:"
	@echo "  make build          - Build server binary"
	@echo "  make build-server   - Build server binary only"
	@echo "  make build-seed     - Build seed binary only"
	@echo "  make test-build     - Test build without creating binary"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make run            - Run server"
	@echo "  make run-seed       - Run seed script"
	@echo "  make tidy           - Run go mod tidy"
	@echo "  make vet            - Run go vet"
	@echo "  make fmt            - Format code"

# Build server binary
build-server:
	@echo "Building server..."
	@go build -v -o bin/server ./cmd/server
	@echo "✓ Server built successfully: bin/server"

# Build seed binary
build-seed:
	@echo "Building seed..."
	@go build -v -o bin/seed ./cmd/seed
	@echo "✓ Seed built successfully: bin/seed"

# Build both
build: build-server build-seed
	@echo "✓ All binaries built successfully"

# Test build (compile without creating binary)
test-build:
	@echo "Testing build..."
	@go build -o /dev/null ./cmd/server && echo "✓ Server: Build OK"
	@go build -o /dev/null ./cmd/seed && echo "✓ Seed: Build OK"
	@go build ./... && echo "✓ All packages: Build OK"
	@echo "✓ Build test passed!"

# Test build with verbose output
test-build-verbose:
	@echo "Testing build (verbose)..."
	@go build -v -o /dev/null ./cmd/server
	@go build -v -o /dev/null ./cmd/seed
	@go build -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean
	@echo "✓ Cleaned"

# Run server
run:
	@go run cmd/server/main.go

# Run seed
run-seed:
	@go run cmd/seed/main.go

# Go mod tidy
tidy:
	@echo "Running go mod tidy..."
	@go mod tidy
	@echo "✓ Done"

# Run go vet (static analysis)
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ No issues found"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Done"

