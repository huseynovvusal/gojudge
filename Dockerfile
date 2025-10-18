FROM golang:1.24.2 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/main.go

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    python3 \
    g++ \
    gcc \
    ca-certificates \
    git \
    make \
    autoconf \
    libtool \
    pkg-config \
    flex \
    bison \
    protobuf-compiler \
    libprotobuf-dev \
    libnl-3-dev \
    libnl-route-3-dev \
    && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/google/nsjail.git /tmp/nsjail && \
    cd /tmp/nsjail && \
    make && \
    cp nsjail /usr/local/bin/nsjail && \
    cd / && rm -rf /tmp/nsjail

COPY --from=builder /app/server /app/server

WORKDIR /app

CMD ["./server"]