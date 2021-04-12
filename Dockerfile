FROM balenalib/raspberry-pi-golang:1.13 as builder
WORKDIR /go/src/github.com/jphastings/jan-poka

# getting "/usr/bin/ld: cannot find -lasound" here :(
# RUN apt-get update && apt-get install -y libasound-dev

# No libnova in raspbian :(
# RUN apt-get update && apt-get install -y libnova

COPY . .
RUN go install -a -ldflags '-extldflags "-static"' -tags 'rpi' -o jp-controller cmd/controller

FROM balenalib/raspberry-pi

ENV INITSYSTEM on

COPY --from=builder /go/bin/jp-controller /jp-controller
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENV JP_PORT 80
EXPOSE 80

CMD ["/jp-controller"]