FROM golang:1.24-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd cmd
COPY internal internal
RUN CGO_ENABLED=0 go build -o importer cmd/importer/main.go

FROM scratch
COPY --from=builder /usr/src/app/importer /importer
ENTRYPOINT ["/importer"]
