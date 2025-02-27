.PHONY: build test clean install list validate

# Binary name
BINARY_NAME=envgen

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install

# Build directory
BUILD_DIR=build

# Build flags
LDFLAGS=-ldflags "-s -w"

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(LDFLAGS)

build-all:
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(LDFLAGS)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(LDFLAGS)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(LDFLAGS)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(LDFLAGS)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(LDFLAGS)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)

install:
	$(GOINSTALL) -v ./...

deps:
	$(GOGET) -v -t ./...

list:
	./$(BINARY_NAME) list

validate:
	./$(BINARY_NAME) validate
