# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
# Binary name
BINARY_NAME=sitp
# Entry
ENTRYPOINT=cmd/sitp/main.go

_all: test build
_build: 
	$(GOBUILD) -o $(BINARY_NAME) -v $(ENTRYPOINT)
_build_alpine:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME)_alpine -v $(ENTRYPOINT)
_test: 
	$(GOTEST) -v ./...
_clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
_run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)