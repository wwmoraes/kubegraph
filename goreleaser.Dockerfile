FROM scratch

COPY kubegraph /usr/bin/local

ENTRYPOINT ["kubegraph"]
