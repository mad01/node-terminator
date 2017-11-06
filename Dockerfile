FROM golang:1.8.4-jessie as builder
ENV buildpath=/go/src/github.com/mad01/k8s-node-terminator
RUN mkdir -p $buildpath
WORKDIR $buildpath
COPY . .

RUN make test
RUN make build/release

FROM debian:8
RUN  apt-get update && apt-get install -y ca-certificates
COPY --from=builder /go/src/github.com/mad01/k8s-node-terminator/_release/k8s-node-terminator /k8s-node-terminator

ENTRYPOINT ["/k8s-node-terminator"]
