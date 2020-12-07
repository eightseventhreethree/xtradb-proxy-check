# build executable binary
FROM golang:1.15.6-buster AS builder

RUN apt update && apt install make

WORKDIR $GOPATH/src/eightseventhreethree/xtradb-proxy-check

ADD . $GOPATH/src/eightseventhreethree/xtradb-proxy-check

RUN make build-linux

# build a small image
FROM scratch

COPY --from=builder /go/src/eightseventhreethree/xtradb-proxy-check/out/xtradb-proxy-check_linux /go/bin/xtradb-proxy-check

ENV CLUSTERCHECK_API_PORT=9200

ENTRYPOINT ["/go/bin/xtradb-proxy-check"]

EXPOSE 9200
