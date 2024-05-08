FROM golang:1 as builder

COPY . /app
WORKDIR /app
RUN make build

FROM alpine:latest

COPY --from=builder /app/bin/edgefig /edgefig
ENTRYPOINT ["/edgefig"]
