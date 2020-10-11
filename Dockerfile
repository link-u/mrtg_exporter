FROM golang:1.15-alpine as build

WORKDIR /go/src/github.com/lnk-u/mrtg_exporter
COPY . .

RUN apk add git gcc g++ musl-dev bash make &&\
    make clean &&\
    make test &&\
    make mrtg_exporter

FROM alpine:3.12

COPY --from=build /go/src/github.com/link-u/mrtg_exporter/mrtg_exporter mrtg_exporter

RUN ["chmod", "a+x", "/mrtg_exporter"]
ENTRYPOINT "/mrtg_exporter"
EXPOSE 8575
