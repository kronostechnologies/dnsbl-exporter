FROM golang:1.14 AS builder
WORKDIR /go/src/github.com/kronostechnologies/dnsbl-exporter/
COPY * ./
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o dnsbl-exporter .

FROM scratch
COPY --from=builder /go/src/github.com/kronostechnologies/dnsbl-exporter/dnsbl-exporter /bin/
ENTRYPOINT ["/bin/dnsbl-exporter"]
