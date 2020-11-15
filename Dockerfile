# syntax = docker/dockerfile:experimental
FROM golang:1.15-alpine AS build

WORKDIR /go/src/kubegraph

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN go mod vendor
RUN --mount=type=cache,target=/root/.cache/go-build go build -mod=vendor -o kubegraph cmd/kubegraph/main.go

FROM alpine:latest

COPY --from=build /go/src/kubegraph/kubegraph /usr/local/bin

### Prepare user
RUN addgroup --gid 1001 kubegraph \
  && adduser \
  --home /home/kubegraph \
  --gecos "" \
  --shell /bin/ash \
  --ingroup kubegraph \
  --disabled-password \
  --uid 1001 \
  kubegraph

WORKDIR /home/kubegraph

USER kubegraph

ENTRYPOINT [ "kubegraph" ]
