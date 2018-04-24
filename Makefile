SHELL = /usr/bin/env bash
BINDIR ?= release
GOBUILD ?= CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-w -s"

all: soccer-robot-remote

deps:
	$(info Checking development deps...)
	@command -v dep > /dev/null || go get -u -v github.com/golang/dep/cmd/dep
	$(info Syncing go deps...)
	@dep ensure -v

vendor: deps

fmt: go.fmt

go.fmt:
	$(info Formatting Go code...)
	@go fmt ./...

test: go.test

go.test:
	$(info Formatting Go code...)
	@go test -cover -v ./...

$(BINDIR)/soccer-robot-remote-linux-amd64: vendor
	$(info Compiling $@...)
	@$(GOBUILD) -o $@ ./cmd/soccer-robot-remote

soccer-robot-remote: $(BINDIR)/soccer-robot-remote-linux-amd64

clean:
	rm -rf vendor $(BINDIR)/soccer-robot-remote-linux-amd64

.PHONY: all soccer-robot-remote deps fmt test go.fmt go.test clean
