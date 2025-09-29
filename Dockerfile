FROM golang:1.24-alpine AS builder

WORKDIR /app

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary, strip symbols for smaller size
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -trimpath -ldflags "-s -w" -o urlshortener ./cmd

# Prepare writable data dir for BadgerDB (if used)
RUN mkdir -p /data

############################

FROM scratch

# Run in release mode by default
ENV GIN_MODE=release \
    PORT=8080 \
    BASE_URL=http://localhost:8080 \
    STORAGE_BACKEND=memory \
    DATA_DIR=/data

# Create non-root user and writable data directory
# Use numeric UID/GID to work in scratch
USER 10001:10001
COPY --from=builder --chown=10001:10001 /data /data

# Copy the binary
COPY --from=builder /app/urlshortener /urlshortener

EXPOSE 8080

ENTRYPOINT ["/urlshortener"]