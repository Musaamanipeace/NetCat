# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=net-cat
CMD_DIR=./cmd/server

.PHONY: all build clean test run deps

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(CMD_DIR)/main.go

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test:
	$(GOTEST) -v -race ./...

run:
	$(GOCMD) run $(CMD_DIR)/main.go $(port)

deps:
	$(GOMOD) tidy