FROM golang:1 as builder

COPY . /app
WORKDIR /app
RUN make build

FROM alpine:latest

ENV PYENV_ROOT=/root/.pyenv
ENV PATH=$PYENV_ROOT/shims:$PYENV_ROOT/bin:$PATH

RUN apk add --no-cache git bash build-base libffi-dev openssl-dev bzip2-dev zlib-dev xz-dev readline-dev sqlite-dev tk-dev && \
    git clone https://github.com/pyenv/pyenv.git ~/.pyenv && \
    pyenv install 3.11 && \
    pyenv global 3.11 && \
    pip install --upgrade pip && \
    pip install --no-cache-dir j2cli

COPY --from=builder /app/bin/edgefig /bin/edgefig
