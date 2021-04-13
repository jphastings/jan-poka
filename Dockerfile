FROM golang:1.16-buster as builder
WORKDIR /go/src/github.com/jphastings/jan-poka

RUN apt-get update && apt-get install -y libnova-dev

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go install -tags 'libnova' github.com/jphastings/jan-poka/cmd/...

#FROM scratch
#COPY --from=builder /go/bin/controller /bin/controller

ENV JP_PORT 2678
EXPOSE 2678

CMD ["/go/bin/controller"]