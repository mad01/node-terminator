FROM golang:1.8.5-alpine3.6 as builder
ENV buildpath=/go/src/github.com/mad01/node-terminator
RUN mkdir -p $buildpath
RUN apk add --update make bash
WORKDIR $buildpath
COPY . .

RUN make test
RUN make build/release

FROM alpine:3.6
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/mad01/node-terminator/_release/node-terminator /node-terminator

ENTRYPOINT ["/node-terminator"]
