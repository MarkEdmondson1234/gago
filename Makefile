# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
LIB=gago
CLI=gagocli
BINARY_NAME=github.com/MarkEdmondson1234/gago/$(LIB)
BINARY_NAME_CLI=github.com/MarkEdmondson1234/gago/$(CLI)
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_UNIX_CLI=$(BINARY_NAME_CLI)_unix
VERSION ?=vlatest

all: deps test install
install: 
		$(GOINSTALL) $(BINARY_NAME)
		$(GOINSTALL) $(BINARY_NAME_CLI)
test: 
		$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...
clean: 
		$(GOCLEAN)
		rm -f $(GOPATH)/bin/$(LIB)
		rm -f $(GOPATH)/bin/$(CLI)

run:
		$(GOINSTALL) $(BINARY_NAME)
		$(GOINSTALL) $(BINARY_NAME_CLI)
		$(CLI)
deps:
		$(GOGET) github.com/olebedev/config
		$(GOGET) google.golang.org/api/analytics/v3
		$(GOGET) google.golang.org/api/analyticsreporting/v4
		$(GOGET) google.golang.org/api/googleapi
		$(GOGET) google.golang.org/api/option
		$(GOGET) golang.org/x/sys/unix

# Cross compilation
release:
		mkdir -p release
		cd gagocli
		env GOOS=linux GOARCH=amd64 $(GOBUILD) -o release/$(CLI)-$(VERSION)-linux-amd64 -v $(BINARY_NAME_CLI)
		env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o release/$(CLI)-$(VERSION)-darwin-amd64 -v $(BINARY_NAME_CLI)
		env GOOS=windows GOARCH=amd64 $(GOBUILD) -o release/$(CLI)-$(VERSION)-windows-amd64.exe -v $(BINARY_NAME_CLI)
