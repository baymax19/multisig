PACKAGES=$(shell go list ./... | grep -v '/vendor/')

all: get_tools get_vendor_deps build

get_tools:
	go get github.com/golang/dep/cmd/dep

build:
	go build -o bin/sentinelcli cmd/sentinelcli/main.go && go build -o bin/sentineld cmd/sentineld/main.go

get_vendor_deps:
	@rm -rf vendor/
	@dep init
	@dep ensure

test:
	@go test $(PACKAGES)

benchmark:
	@go test -bench=. $(PACKAGES)

.PHONY: all build test benchmark