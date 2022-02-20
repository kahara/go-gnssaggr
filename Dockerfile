# syntax=docker/dockerfile:1.1.7-experimental

FROM golang:1.17.7-alpine3.15 as base

RUN mkdir /workdir
COPY go.* /workdir/
COPY src /workdir/src

WORKDIR /workdir
RUN go build -o gnssaggr ./src

FROM base as production

COPY --from=base /workdir/gnssaggr /usr/local/bin/gnssaggr
COPY docker/entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
