.PHONY: build run clean all

# Default target
all: build

# Build the game
build:
	@echo "Building CivIdleCli..."
	@go build -o cividlecli

# Run the game
run: build
	@echo "Running CivIdleCli..."
	@./cividlecli

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f cividlecli

# Build for multiple platforms
release:
	@echo "Building releases for multiple platforms..."
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/cividlecli-darwin-amd64
	@GOOS=darwin GOARCH=arm64 go build -o ./bin/cividlecli-darwin-arm64
	@GOOS=linux GOARCH=amd64 go build -o ./bin/cividlecli-linux-amd64
	@GOOS=windows GOARCH=amd64 go build -o ./bin/cividlecli-windows-amd64.exe
	@echo "Release builds complete."
