FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

ENV CGO_ENABLED=0 GOOS=linux

RUN go build -o urlshortener ./cmd

############################

FROM scratch

ENV PORT=8080 BASE_URL=http://localhost:8080

COPY --from=builder /app/urlshortener /urlshortener

EXPOSE 8080

ENTRYPOINT ["/urlshortener"]