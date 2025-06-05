FROM golang:1.24-alpine AS builder
RUN apk update && apk add make

WORKDIR /usr/src/app

COPY go.mod go.mod
# COPY go.sum go.sum
RUN go mod download

COPY Makefile Makefile
COPY cmd cmd
COPY internal internal
RUN make

FROM scratch

COPY --from=builder /usr/src/app/dist/server /server

ENTRYPOINT ["/server"]
