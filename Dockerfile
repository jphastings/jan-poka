FROM golang:1.16-buster as builder
WORKDIR /go/src/github.com/jphastings/jan-poka

#RUN apt-get update && apt-get install -y libnova

COPY . .
RUN go install -a -ldflags '-extldflags "-static"' -tags '' github.com/jphastings/jan-poka/cmd/...

FROM scratch

COPY --from=builder /go/bin/controller /jp-controller

ENV JP_PORT 80
EXPOSE 80

CMD ["/jp-controller"]