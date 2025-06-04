################################################################################

GO ?= go
GOENV ?= CGO_ENABLED=0

SERVER_BIN ?= dist/server

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
