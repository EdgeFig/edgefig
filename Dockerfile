FROM golang:1 as builder

COPY . /app
WORKDIR /app
RUN make build

FROM alpine:latest

RUN apk add --no-cache python3 py-pip && \
    pip3 install --break-system-packages j2cli

COPY --from=builder /app/bin/edgefig /bin/edgefig
