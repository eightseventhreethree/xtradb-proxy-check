# Go parameters
GOCMD=go
BINARY_NAME=xtradb-proxy-check
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build: 
	$(GOCMD) -o $(BINARY_NAME) -v

test: 
	$(GOCMD) test -v ./...

clean: 
	$(GOCMD) clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

vendor:
	$(GOCMD) mod vendor

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
