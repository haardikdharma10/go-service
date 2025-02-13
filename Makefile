# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=go-service
MAIN_PATH=./main.go

# Make all
all: test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	./$(BINARY_NAME)

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	$(GOGET) -v -d ./...

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Build Docker image
docker-build:
	docker build -t $(BINARY_NAME) .

# Run Docker container
docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

# Help command
help:
	@echo "make - Compile the project"
	@echo "make build - Build the binary"
	@echo "make run - Run the application"
	@echo "make clean - Remove binary and coverage files"
	@echo "make test - Run tests"
	@echo "make test-coverage - Run tests with coverage"
	@echo "make deps - Download dependencies"
	@echo "make fmt - Format code"
	@echo "make vet - Vet code"
	@echo "make docker-build - Build Docker image"
	@echo "make docker-run - Run Docker container"

.PHONY: all build run clean test test-coverage deps fmt vet docker-build docker-run help