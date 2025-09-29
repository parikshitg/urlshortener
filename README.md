# URLSHORTENER

![urlshortner logo](assets/logo.png)

A production-ready URL shortener service built in Go with comprehensive features for modern web applications.

## Features

### Core Functionality
- **URL Shortening**: Convert long URLs into short, shareable links
- **URL Resolution**: Redirect short codes to original URLs with 302 redirects
- **Custom Code Length**: Configurable short code length (default: 7 characters)
- **TTL Support**: Automatic expiration of shortened URLs with configurable duration

### Storage Options
- **In-Memory Storage**: Fast, ephemeral storage for development/testing
- **BadgerDB**: Persistent, embedded key-value store with native TTL support
- **Configurable Backend**: Switch between storage backends via environment variables

### Security & Performance
- **Rate Limiting**: Per-IP fixed-window rate limiting to prevent abuse
- **CORS Support**: Configurable Cross-Origin Resource Sharing headers
- **Request Validation**: URL format validation and sanitization
- **Concurrent Safe**: Thread-safe operations for high-traffic scenarios

### Analytics & Monitoring
- **Top Domains**: Track and rank most frequently shortened domains
- **Health Checks**: Built-in health and readiness endpoints
- **Structured Logging**: JSON/text logging with configurable levels
- **Metrics Collection**: Domain-based analytics and usage statistics

### QR Code Generation
- **QR Code API**: Generate QR codes for any URL
- **Configurable Size**: Customizable QR code dimensions
- **PNG Output**: High-quality PNG image generation

### Developer Experience
- **Docker Support**: Containerized deployment with multi-stage builds
- **Comprehensive Testing**: Unit tests, integration tests, and mocks
- **Makefile Targets**: Easy build, run, and test commands
- **Environment Configuration**: Extensive configuration via environment variables

### API Features
- **RESTful Design**: Clean, intuitive API endpoints
- **JSON Responses**: Consistent JSON API responses
- **Error Handling**: Proper HTTP status codes and error messages
- **Middleware Support**: CORS, rate limiting, and logging middleware

Note: This project uses Gin instead of net/http or Gorilla Mux because Gin provides:

- Better routing performance and logging.
- Powerful request context (`*gin.Context`) for passing data.
- Easy JSON encode/decode and binding helpers out of the box.
- Built-in middleware (logger, recovery) and a clean developer experience.

## Configuration

Environment variables:

- `PORT` – HTTP port (default: `8080`)
- `BASE_URL` – Base URL used to construct returned short URLs (default: `http://localhost:8080`)
- `CODE_LENGTH` – Length of generated short code (default: `7`)
- `TOP_N` – Default number of top domains to return (default: `3`)
- `EXPIRY` – TTL for shortened URLs, Go duration (default: `1h`)
- `LOG_LEVEL` – `debug|info|warn|error|fatal` (default: `info`)
- `LOG_FORMAT` – `text|json` (default: `text`)

CORS:

- `CORS_ALLOWED_ORIGINS` – CSV of origins or `*` (default: `*`)
- `CORS_ALLOWED_METHODS` – CSV methods (default: `GET,POST,PUT,DELETE,OPTIONS`)
- `CORS_ALLOWED_HEADERS` – CSV headers or `*` (default: `*`)
- `CORS_MAX_AGE` – Seconds to cache preflight (default: `43200`)
- `CORS_ALLOW_CREDENTIALS` – `true|false` (default: `false`)

Rate Limiter (per-IP, fixed window):

- `RATE_LIMIT_MAX_TOKENS` – Requests allowed per window (default: `1000`)
- `RATE_LIMIT_EXPIRY` – Window duration, Go duration (default: `1h`)
- `RATE_LIMIT_PURGE_INTERVAL` – Cleanup interval, Go duration (default: `10m`)

Storage Backend:

- `STORAGE_BACKEND` – `memory` or `badger` (default: `memory`)
- `DATA_DIR` – Database directory for BadgerDB (default: `./data`)

## Build and Run (Locally)

Build:

```bash
make build
```

Run:

```bash
PORT=8080 BASE_URL=http://localhost:8080 ./urlshortener
```

Or using the Makefile targets:

```bash
# Run with in-memory storage
make run

# Run with BadgerDB storage
make run-badger
```

Health check:

```bash
curl -i http://localhost:8080/health
```

## Docker

Build image:

```bash
docker build -t urlshortener .
```

Run container:

```bash
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e BASE_URL=http://localhost:8080 \
  urlshortener
```

## API

### Shorten URL

`POST /v1/shorten`

Request body:

```
{ "url": "https://www.example.com/very/long/path" }
```

Curl example:

```
curl -i -X POST http://localhost:8080/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.example.com/very/long/path"}'
```

Successful response (200):

```json
{ "shortUrl": "http://localhost:8080/abc1234" }
```

### Metrics (Top Domains)

`POST /v1/metrics`

Request body:

```json
{ "topN": 5 }
```

Curl example:

```bash
curl -i -X POST http://localhost:8080/v1/metrics \
  -H "Content-Type: application/json" \
  -d '{"topN":5}'
```

Successful response (200):

```json
[
  { "rank": 0, "domain": "example.com", "shortened": 15 },
  { "rank": 1, "domain": "google.com",  "shortened": 8  },
  { "rank": 2, "domain": "github.com",  "shortened": 3  }
]
```

### Resolve Short URL

`GET /{code}` – Redirects to the original URL.

Example:

```bash
curl -i http://localhost:8080/abc1234
```

Response: `302 Found` with `Location` header pointing to the original URL.

### QR Code Generation

`POST /v1/qr`

Request body:

```json
{ "url": "https://www.example.com", "size": 256 }
```

Curl example:

```bash
curl -i -X POST http://localhost:8080/v1/qr \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.example.com","size":256}' \
  --output qr.png
```

Response: `200 OK` with `image/png` content type and QR code image.

### Health Check

`GET /health` – Basic health check

`GET /health/ready` – Readiness check (includes storage connectivity)

Example:

```bash
curl -i http://localhost:8080/health
```

Response: `200 OK` with service status information.

## Testing

Run all tests:

```bash
make test
```

Run specific test suites:

```bash
# Unit tests
go test ./...

# Integration tests for BadgerDB
go test ./internal/storage/badgerdb -v

# API tests
go test ./api/v1 -v
```

## Development

Generate mocks:

```bash
make generate-mocks
```

Build for production:

```bash
make build
```

