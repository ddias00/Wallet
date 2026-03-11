.PHONY: run build test test-coverage clean docker-up docker-down docker-dev docker-logs swagger deps

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./test/...

test-coverage:
	go test -v -coverprofile=coverage.out ./test/...
	go tool cover -html=coverage.out -o coverage.html

test-report:
	go test -v -coverprofile=coverage.out ./test/...
	go tool cover -func=coverage.out

swagger:
	swag init -g cmd/api/main.go -o docs

clean:
	rm -rf bin/ coverage.out coverage.html tmp/

# Docker Production
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

docker-rebuild:
	docker-compose up -d --build

# Docker Development (Hot Reload)
docker-dev:
	docker-compose -f docker-compose.dev.yml up

docker-dev-down:
	docker-compose -f docker-compose.dev.yml down

# Dependencies
deps:
	go mod download
	go mod tidy
