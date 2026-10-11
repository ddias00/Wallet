#!/bin/bash

echo "========================================"
echo "Digital Wallet API - Setup Script"
echo "========================================"
echo ""

echo "[1/5] Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed. Please install Go 1.27.2 or higher."
    exit 1
fi
go version
echo ""

echo "[2/5] Downloading dependencies..."
go mod download
if [ $? -ne 0 ]; then
    echo "ERROR: Failed to download dependencies."
    exit 1
fi
echo ""

echo "[3/5] Tidying go.mod..."
go mod tidy
echo ""

echo "[4/5] Checking Docker..."
if ! command -v docker &> /dev/null; then
    echo "WARNING: Docker is not installed. You'll need to run PostgreSQL manually."
else
    docker --version
    echo "Docker is installed!"
fi
echo ""

echo "[5/5] Setup complete!"
echo ""
echo "========================================"
echo "Next Steps:"
echo "========================================"
echo "1. Start PostgreSQL and API:"
echo "   docker compose up -d"
echo ""
echo "2. Or run API only (requires PostgreSQL running):"
echo "   go run cmd/api/main.go"
echo ""
echo "3. Access Swagger documentation:"
echo "   http://localhost:8080/swagger/index.html"
echo ""
echo "4. Health check:"
echo "   http://localhost:8080/health"
echo "========================================"
