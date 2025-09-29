# URLSHORTENER

![urlshortner logo](assets/logo.png)

This is a simple Golang server written in Go for url shortener.

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

## Build and Run (Locally)

Build:

```bash
make build
```

Run:

```bash
PORT=8080 BASE_URL=http://localhost:8080 ./urlshortener
```

Or using the Makefile target:

```bash
make run
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

