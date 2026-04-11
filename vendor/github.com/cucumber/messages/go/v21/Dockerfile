#syntax=docker/dockerfile:1.4

# Base image
ARG GO_VERSION=1.19
FROM golang:${GO_VERSION}-alpine AS golang
WORKDIR /cucumber


# Dummy stage for generated code, overriden in main build
FROM scratch AS schema-codegen


FROM golang AS with-dependencies

COPY --link go.mod go.sum .
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download && go mod verify


FROM golang AS tested

RUN apk add gcc libc-dev

COPY --link . .
COPY --link --from=with-dependencies /go/pkg /go/pkg
COPY --link --from=schema-codegen /*.go .

RUN gofmt -w .
RUN go test --v
