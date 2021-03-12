FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/kronostechnologies/dnsbl-exporter/
COPY * ./
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o dnsbl-exporter .
RUN echo "nobody:x:65534:65534:nobody:/:" > /tmp/passwd

FROM scratch
COPY --from=builder /go/src/github.com/kronostechnologies/dnsbl-exporter/dnsbl-exporter /bin/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /tmp/passwd /etc/passwd

USER 65534:65534
ENTRYPOINT ["/bin/dnsbl-exporter"]