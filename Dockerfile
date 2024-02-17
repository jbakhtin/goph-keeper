# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder

# Set destination for COPY
WORKDIR /app

RUN apt-get update && \
    apt-get -y install gcc

# Download Go appmodules
COPY go.mod .
RUN go mod download
COPY . .

# Build
RUN go build -a -o /app/bin/server /app/cmd/server/main.go

# Run
CMD ["/app/bin/server"]