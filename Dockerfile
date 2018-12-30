FROM balenalib/raspberry-pi-golang:1.11 as builder
WORKDIR /go/src/corviator
COPY . .
RUN go install -a -ldflags '-extldflags "-static"'

FROM balenalib/scratch

ENV INITSYSTEM on

COPY --from=builder /go/bin/corviator /corviator
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
CMD ["/corviator"]