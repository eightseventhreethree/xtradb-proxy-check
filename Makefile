# Go parameters
GOCMD=go
BINARY_NAME=xtradb-proxy-check
BINARY_LINUX=$(BINARY_NAME)_linux
BUILD_BASE=$(GOCMD) build -o out/

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

container: build-linux
	
run:
	$(GOCMD) run main.go

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD_BASE)$(BINARY_LINUX) -v
