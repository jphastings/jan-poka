FROM golang:1.17-bullseye as builder
WORKDIR /go/src/github.com/jphastings/jan-poka

RUN apt-get update && apt-get install -y libnova-dev

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN mkdir -p bin/ && cd bin/ && go build -tags 'libnova' github.com/jphastings/jan-poka/cmd/...

FROM debian:bullseye-slim
COPY --from=builder /go/src/github.com/jphastings/jan-poka/bin/controller /usr/bin/

RUN apt-get update && apt-get install -y libnova-0.16-0

ENV JP_PORT 2678
EXPOSE 2678

ENV JP_MQTTPORT 1883
EXPOSE 1883

CMD ["/usr/bin/controller"]