FROM balenalib/raspberry-pi-golang:1.11 as builder
WORKDIR /go/src/github.com/jphastings/corviator

# getting "/usr/bin/ld: cannot find -lasound" here :(
# RUN apt-get update && apt-get install -y libasound-dev

# No libnova in raspbian :(
# RUN apt-get update && apt-get install -y libnova

COPY . .
RUN go install -a -ldflags '-extldflags "-static"' -tags 'rpi'

FROM balenalib/raspberry-pi

ENV INITSYSTEM on

COPY --from=builder /go/bin/corviator /corviator
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENV CORVIATOR_PORT 80
EXPOSE 80

CMD ["/corviator"]