# Go parameters
GOCMD=go
BINARY_NAME=xtradb-proxy-check
BINARY_LINUX=$(BINARY_NAME)_linux
BUILD_BASE=$(GOCMD) build -o out/
TAG_NAME:=$(shell git describe --abbrev=0 --tags)

all: test build
build:
	$(BUILD_BASE)$(BINARY_NAME) -v

test:
	$(GOCMD) test -v ./...

clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_LINUX)

vendor:
	$(GOCMD) mod vendor

run:
	$(GOCMD) run main.go

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD_BASE)$(BINARY_LINUX) -v

build-container:
	docker build --no-cache -t xtradb-proxy-check:$(TAG_NAME) .
