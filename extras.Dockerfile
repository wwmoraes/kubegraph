# hadolint ignore=DL3007
FROM wwmoraes/kubegraph:latest AS kubegraph

FROM alpine:3.16

RUN apk add --no-cache --update \
  graphviz=~3.0 \
  librsvg=~2.54 \
  ;

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

COPY --from=kubegraph /kubegraph /usr/local/bin/kubegraph

USER kubegraph

ENTRYPOINT [ "kubegraph" ]
