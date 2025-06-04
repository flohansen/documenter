################################################################################

GO ?= go
GOENV ?= CGO_ENABLED=0

SERVER_BIN ?= dist/server

################################################################################

all: build

build:
	$(GOENV) $(GO) build -o $(SERVER_BIN) cmd/server/main.go
