FROM wwmoraes/kubegraph:latest

FROM alpine:latest

RUN apk add --no-cache --update graphviz librsvg

RUN addgroup --gid 10001 kubegraph \
  && adduser \
  --home /home/kubegraph \
  --gecos "" \
  --shell /bin/ash \
  --ingroup kubegraph \
  --disabled-password \
  --uid 10001 \
  kubegraph

WORKDIR /home/kubegraph

COPY --from=0 /usr/local/bin/kubegraph /usr/local/bin/kubegraph

USER kubegraph

ENTRYPOINT [ "kubegraph" ]
