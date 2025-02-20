FROM golang:1 AS builder

COPY . /app
WORKDIR /app
RUN make build

FROM alpine:latest

RUN apk add --no-cache git bash python3 py3-pip && \
    pip3 install --no-cache-dir --break-system-packages jinjanator

COPY --from=builder /app/bin/edgefig /bin/edgefig
