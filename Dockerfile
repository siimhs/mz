# Stage 1: Build the binary
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy Go modules and source code
# COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

COPY . ./

# Build the binary
RUN go build -o mz .

# Stage 2: Final stage
FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates
# Set the working directory
WORKDIR /output

# Copy the binary from the builder stage
COPY --from=builder /app/mz .

# Entrypoint (optional, for testing purposes)
ENTRYPOINT ["/output/mz"]
