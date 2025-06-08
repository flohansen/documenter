################################################################################

# Binaries
GO ?= go
GOENV ?= CGO_ENABLED=0
DOCKER ?= docker

# Tools
LOCALBIN ?= $(shell pwd)/bin
SQLC ?= $(LOCALBIN)/sqlc
MOCKGEN ?= $(LOCALBIN)/mockgen

SQLC_VERSION ?= v1.29.0
MOCKGEN_VERSION ?= v0.5.2

# Build Outputs
IMPORTER_BIN ?= dist/importer
IMPORTER_IMAGE ?= documenter-importer
IMPORTER_TAG ?= latest

################################################################################

all: build

.PHONY: build
build: generate
	$(GOENV) $(GO) build -o $(IMPORTER_BIN) cmd/importer/main.go

.PHONY: generate
generate: mockgen sqlc
	$(SQLC) generate
	PATH=$(LOCALBIN):$(PATH) $(GO) generate ./...

.PHONY: test
test: generate
	$(GO) test ./... -cover -race

.PHONY: docker-build
docker-build: generate
	$(DOCKER) build -t $(IMPORTER_IMAGE):$(IMPORTER_TAG) .

$(LOCALBIN):
	mkdir -p $@

.PHONY: sqlc
sqlc: $(SQLC)
$(SQLC): $(LOCALBIN)
	GOBIN=$(LOCALBIN) $(GO) install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)

.PHONY: mockgen
mockgen: $(MOCKGEN)
$(MOCKGEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) $(GO) install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)
