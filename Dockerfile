# Build Geth in a stock Go builder container
FROM golang:1.10.3-alpine3.8 as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

COPY . /go/src/ipninja
WORKDIR /go/src/ipninja

RUN go get github.com/tools/godep && godep restore
RUN go build -o testapp

# Pull binaries into a second stage deploy alpine container
FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY ipn.db /ipninja/
COPY ./static /ipninja/static
COPY --from=builder /go/src/ipninja/testapp /ipninja/

WORKDIR /ipninja

ENV DBPATH /ipninja/ipn.db

ENTRYPOINT ["./testapp"]