FROM balenalib/raspberry-pi-golang:1.11 as builder
WORKDIR /go/src/corviator
COPY periph.go ./
RUN go install -a -ldflags '-extldflags "-static"'

FROM balena-os/scratch

ENV INITSYSTEM on

COPY --from=builder /go/bin/corviator /corviator
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
CMD ["/corviator"]