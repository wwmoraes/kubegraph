FROM golang:1.15-alpine AS build

WORKDIR /go/src/kubegraph

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal
COPY icons/icons.go icons/icons.go
RUN go build -o kubegraph cmd/kubegraph/main.go

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