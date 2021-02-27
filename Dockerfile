# syntax = docker/dockerfile:experimental
FROM golang:1.15-alpine AS build

WORKDIR /go/src/kubegraph

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN go mod vendor
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/root/.cache/go-build go build \
  -mod=vendor \
  -o kubegraph cmd/kubegraph/main.go

FROM alpine:3 AS user

RUN echo "kubegraph:x:10001:kubegraph" >> /tmp/group
RUN echo "kubegraph:x:10001:10001::/:/dev/null" >> /tmp/passwd

FROM scratch

COPY --from=user /tmp/passwd /etc/passwd
COPY --from=user /tmp/group /etc/group
COPY --from=build --chown=kubegraph:kubegraph /go/src/kubegraph/kubegraph /

USER kubegraph:kubegraph

ARG OCI_VERSION
ARG OCI_CREATED
ARG OCI_REVISION
ARG OCI_AUTHORS
ARG OCI_VENDOR

LABEL org.opencontainers.image.title=kubegraph
LABEL org.opencontainers.image.description="Kubernetes resource graph generator"
LABEL org.opencontainers.image.url=https://github.com/wwmoraes/kubegraph
LABEL org.opencontainers.image.source=https://github.com/wwmoraes/kubegraph
LABEL org.opencontainers.image.version=${OCI_VERSION}
LABEL org.opencontainers.image.created=${OCI_CREATED}
LABEL org.opencontainers.image.revision=${OCI_REVISION}
LABEL org.opencontainers.image.licenses=MIT
LABEL org.opencontainers.image.authors="${OCI_AUTHORS}"
LABEL org.opencontainers.image.documentation=https://github.com/wwmoraes/kubegraph
LABEL org.opencontainers.image.vendor="${OCI_VENDOR}"

ENTRYPOINT [ "/kubegraph" ]
