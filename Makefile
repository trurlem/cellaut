# Variables
BINARY_NAME=cellaut
GOFMT=go fmt
GOBUILD=go build
GOCLEAN=go clean
GORUN=go run
GOLINT=golangci-lint run

# Default target
all: format build

# Format Go source code
format:
	$(GOFMT) ./...

# Compile the Go program
build: format
	$(GOBUILD) -o ./bin/$(BINARY_NAME) && GOOS=darwin GOARCH=arm64 $(GOBUILD) -o ./bin/$(BINARY_NAME)-AppleSilicon

# Run the Go program
run: build
	./bin/$(BINARY_NAME)

# Lint the Go source code
lint:
	$(GOLINT) ./...

# Clean up the compiled binary
clean:
	$(GOCLEAN)
	rm -f ./bin/*

.PHONY: all format build run lint clean
