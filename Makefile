APP_NAME=urlshortener

build:
	GO111MODULE=on go build -o $(APP_NAME) ./cmd

run:
	PORT=8080 BASE_URL=http://localhost:8080 EXPIRY=1m ./$(APP_NAME)

tidy:
	go mod tidy

test:
	go test ./...

run-badger:
	STORAGE_BACKEND=badger DATA_DIR=./data PORT=8080 BASE_URL=http://localhost:8080 EXPIRY=1h ./$(APP_NAME)

generate-mocks:
	mockgen -source=internal/storage/storage.go -destination=internal/storage/mocks/mock_storage.go -package=mocks

docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run --rm -p 8080:8080 \
	  -e PORT=8080 \
	  -e BASE_URL=http://localhost:8080 \
	  -e STORAGE_BACKEND=badger \
	  -e DATA_DIR=/data \
	  -v $(PWD)/data:/data \
	  $(APP_NAME)