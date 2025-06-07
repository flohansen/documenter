################################################################################

GO ?= go
GOENV ?= CGO_ENABLED=0
DOCKER ?= docker

IMPORTER_BIN ?= dist/server
IMPORTER_IMAGE ?= documenter-server
IMPORTER_TAG ?= latest

################################################################################

all: build

.PHONY: build
build: generate
	$(GOENV) $(GO) build -o $(IMPORTER_BIN) cmd/server/main.go

.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: test
test: generate
	$(GO) test ./... -cover -race

.PHONY: docker-build
docker-build:
	$(DOCKER) build -t $(IMPORTER_IMAGE):$(IMPORTER_TAG) .
