################################################################################

GO ?= go
GOENV ?= CGO_ENABLED=0
DOCKER ?= docker

SERVER_BIN ?= dist/server
SERVER_IMAGE ?= documenter-server
SERVER_TAG ?= latest

################################################################################

all: build

.PHONY: build
build: generate
	$(GOENV) $(GO) build -o $(SERVER_BIN) cmd/server/main.go

.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: test
test: generate
	$(GO) test ./... -cover -race

.PHONY: docker-build
docker-build:
	$(DOCKER) build -t $(SERVER_IMAGE):$(SERVER_TAG) .
