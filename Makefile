APP_NAME=urlshortener

build:
	GO111MODULE=on go build -o $(APP_NAME) ./cmd

run:
	PORT=8080 BASE_URL=http://localhost:8080 ./$(APP_NAME)

tidy:
	go mod tidy

test:
	go test ./...