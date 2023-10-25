# syntax=docker/dockerfile:1

FROM golang:1.21

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
RUN go mod download
COPY . .

# Build
RUN go build -a -o /app/bin/server /app/cmd/server/main.go
EXPOSE 8000

# Run
CMD ["/app/bin/server"]