SHELL = /usr/bin/env bash
BINDIR ?= release
GOBUILD ?= CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-w -s"
YARN ?= yarn

GO_FILES = `find cmd pkg -name '*.go'`

all: deps go.build js.build

deps:
	$(info Checking development deps...)
	@command -v yarn > /dev/null || { printf 'Please install yarn (follow the steps in DEVELOPMENT.md)\n'; exit 1; }
	@command -v go > /dev/null || { printf 'Please install go (follow the steps in DEVELOPMENT.md)\n'; exit 1; }
	@command -v dep > /dev/null || go get -u -v github.com/golang/dep/cmd/dep
	@command -v gometalinter > /dev/null || go get -u -v github.com/alecthomas/gometalinter
	@command -v unconvert > /dev/null || gometalinter -i
	@command -v misspell > /dev/null || gometalinter -i
	$(info Syncing go deps...)
	@dep ensure -v
	@$(YARN) install
	@$(YARN) --cwd front install

vendor: deps

go.fmt: deps
	$(info Formatting Go code...)
	@gofmt -w -s $(GO_FILES)
	@unconvert -safe -apply ./...
	@misspell -w ./...

go.lint: deps
	$(info Linting Go code...)
	@gometalinter --fast $(GO_FILES)

test: go.test

go.test: deps
	$(info Formatting Go code...)
	@go test -cover -v ./...

$(BINDIR)/soccer-robot-remote-linux-amd64: vendor
	$(info Compiling $@...)
	@$(GOBUILD) -o $@ ./cmd/soccer-robot-remote

soccer-robot-remote: $(BINDIR)/soccer-robot-remote-linux-amd64

go.build: soccer-robot-remote

js.fmt: deps
	$(info Formatting JS code...)
	@$(YARN) run js.fmt

js.build: deps
	@$(YARN) build
	@rm -rf $(BINDIR)/front
	@mv ./front/build $(BINDIR)/front

md.fmt: deps
	$(info Formatting MD code...)
	@$(YARN) run md.fmt

lint: go.lint

fmt: go.fmt js.fmt md.fmt

clean:
	rm -rf node_modules front/node_modules vendor $(BINDIR)/soccer-robot-remote-linux-amd64

.PHONY: all soccer-robot-remote deps fmt test go.build go.fmt go.test go.lint js.build js.fmt md.fmt clean
