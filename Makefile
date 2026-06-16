BINARY     = dauth
CMD        = ./cmd/server
BUILD_DIR  = ./bin

.PHONY: build run test lint tidy migrate-up migrate-down clean

build:
	go build -o $(BUILD_DIR)/$(BINARY) $(CMD)

run:
	go run $(CMD)

test:
	go test ./... -v -race

test/cover:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

migrate-up:
	migrate -path ./migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path ./migrations -database "$$DATABASE_URL" down

clean:
	rm -rf $(BUILD_DIR)
