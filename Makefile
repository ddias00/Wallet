.PHONY: run build deps swagger clean \
        test test-verbose test-short test-service test-handler test-repository \
        coverage coverage-func \
        docker-up docker-down docker-logs docker-rebuild docker-dev docker-dev-down

# ----- App -----
run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

deps:
	go mod download
	go mod tidy

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs

clean:
	rm -rf bin/ tmp/ coverage coverage.out coverage.html coverage_*

# ----- Tests -----
test:
	go test ./... -v

test-verbose:
	go test ./... -v -cover

test-short:
	go test ./... -short -v

test-service:
	go test ./internal/service/... -v -cover

test-handler:
	go test ./internal/handler/... -v -cover

test-repository:
	go test ./internal/repository/... -v -cover

# ----- Coverage -----
coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

coverage-func:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

# ----- Docker (prod) -----
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f api

docker-rebuild:
	docker compose up -d --build

# ----- Docker (dev, hot reload via Air) -----
docker-dev:
	docker compose -f docker-compose.dev.yml up

docker-dev-down:
	docker compose -f docker-compose.dev.yml down
