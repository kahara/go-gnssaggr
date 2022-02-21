# syntax=docker/dockerfile:1.1.7-experimental

FROM golang:1.17.7-bullseye as build

RUN mkdir /workdir
COPY go.* /workdir/
COPY src /workdir/src

WORKDIR /workdir
RUN go build -o gnssaggr ./src

FROM gcr.io/distroless/base-debian11 as production

COPY --from=build /workdir/gnssaggr /

CMD ["/gnssaggr"]
