SHELL = /usr/bin/env bash
BINDIR ?= release
GOARCH ?= $(shell go env GOARCH)
GOOS ?= $(shell go env GOOS)
GOBUILD ?= CGO_ENABLED=0 go build -ldflags="-w -s"
YARN ?= yarn
DOCKER_IMAGE_VERSION ?= "latest"

GO_FILES = $(shell find cmd pkg -name '*.go')

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
	@unconvert -safe -apply ./pkg/... ./cmd/...
	@misspell -w ./pkg/... ./cmd/...

go.lint: deps
	$(info Linting Go code...)
	@gometalinter --fast ./pkg/... ./cmd/...

test: go.test js.test

go.test: deps
	$(info Testing Go code...)
	@go test -cover -coverprofile=coverage.txt -covermode=atomic -v ./...

$(BINDIR)/trcd-$(GOOS)-$(GOARCH): vendor $(GO_FILES)
	$(info Compiling $@...)
	@$(GOBUILD) -o $@ ./cmd/trcd

trcd: $(BINDIR)/trcd-$(GOOS)-$(GOARCH)

$(BINDIR)/srrs-$(GOOS)-$(GOARCH): vendor $(GO_FILES)
	$(info Compiling $@...)
	@$(GOBUILD) -o $@ ./cmd/srrs

srrs: $(BINDIR)/srrs-$(GOOS)-$(GOARCH)

$(BINDIR)/srrs-$(GOOS)-$(GOARCH)-noauth: vendor $(GO_FILES)
	$(info Compiling $@...)
	@$(GOBUILD) -tags noauth -o $@ ./cmd/srrs

srrs-noauth: $(BINDIR)/srrs-$(GOOS)-$(GOARCH)-noauth

$(BINDIR)/relay-$(GOOS)-$(GOARCH): vendor $(GO_FILES)
	$(info Compiling $@...)
	@$(GOBUILD) -o $@ ./cmd/relay

relay: $(BINDIR)/relay-$(GOOS)-$(GOARCH)

go.build: srrs

js.fmt: deps
	$(info Formatting JS code...)
	@$(YARN) run js.fmt

js.build: deps
	@$(YARN) build
	@rm -rf $(BINDIR)/front
	@mv ./front/build $(BINDIR)/front

$(BINDIR)/front: js.build

js.test: deps
	@$(info Testing JS code...)
	@$(YARN) test --coverage

md.fmt: deps
	$(info Formatting MD code...)
	@$(YARN) run md.fmt

lint: go.lint

build: go.build js.build

fmt: go.fmt js.fmt md.fmt

docker: $(BINDIR)/front $(BINDIR)/srrs-linux-amd64
	docker build -t rvolosatovs/srr:$(DOCKER_IMAGE_VERSION) .

clean:
	rm -rf node_modules front/node_modules vendor $(BINDIR)/srrs-* $(BINDIR)/front*

.PHONY: all srrs srrs-noauth relay trcd deps fmt test go.build go.fmt go.test go.lint js.build js.fmt md.fmt clean
