SHELL = /usr/bin/env bash
BINDIR ?= release
GOBUILD ?= CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-w -s"
YARN ?= yarn

all: deps fmt go.build js.build

deps:
	$(info Checking development deps...)
	@command -v go > /dev/null || { printf 'Please install go (follow the steps in DEVELOPMENT.md)\n'; exit 1; }
	@command -v yarn > /dev/null || { printf 'Please install yarn (follow the steps in DEVELOPMENT.md)\n'; exit 1; }
	@command -v dep > /dev/null || go get -u -v github.com/golang/dep/cmd/dep
	$(info Syncing go deps...)
	@dep ensure -v
	@$(YARN) install
	@$(YARN) --cwd front install

vendor: deps

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

go.build: soccer-robot-remote

js.fmt:
	$(info Formatting JS code...)
	@$(YARN) run js.fmt

js.build:
	@$(YARN) build

md.fmt:
	$(info Formatting MD code...)
	@$(YARN) run md.fmt

fmt: go.fmt js.fmt md.fmt

clean:
	rm -rf vendor $(BINDIR)/soccer-robot-remote-linux-amd64

.PHONY: all soccer-robot-remote deps fmt test go.build go.fmt go.test js.build js.fmt md.fmt clean
