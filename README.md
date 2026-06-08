# Digital Wallet API

A complete REST API for a Digital Wallet platform built with Golang, following Clean Architecture principles.

## 🚀 Features

### User Management
- User registration with CPF, email, phone validation
- JWT-based authentication with refresh tokens
- Account blocking/unblocking

### Wallet Operations
- Get wallet balance
- Credit wallet
- Debit wallet

### Transactions
- Transfer money between users
- Transaction history with pagination
- Transaction status tracking (pending, confirmed, cancelled)

### Payments
- Generate Boleto
- PIX payment support
- QR Code payment generation
- Payment processing

### Notifications
- Transaction notifications (console-based, ready for email/SMS integration)

## 🛠️ Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT with refresh tokens
- **Password Hashing**: bcrypt
- **API Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose
- **Migrations**: golang-migrate compatible

## 📁 Project Structure

```
Wallet/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── domain/                  # Domain models/entities
│   │   ├── user.go
│   │   ├── wallet.go
│   │   ├── transaction.go
│   │   └── payment.go
│   ├── handler/                 # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── wallet_handler.go
│   │   ├── transaction_handler.go
│   │   └── payment_handler.go
│   ├── service/                 # Business logic
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── wallet_service.go
│   │   ├── transaction_service.go
│   │   ├── payment_service.go
│   │   └── notification_service.go
│   ├── repository/              # Data access layer
│   │   ├── user_repository.go
│   │   ├── wallet_repository.go
│   │   ├── transaction_repository.go
│   │   ├── payment_repository.go
│   │   └── refresh_token_repository.go
│   ├── middleware/              # HTTP middleware
│   │   ├── auth.go
│   │   └── logging.go
│   └── utils/
│       └── response.go          # Response helpers
├── pkg/
│   ├── jwt/
│   │   └── jwt.go               # JWT utilities
│   └── validator/
│       └── validator.go         # CPF, email, phone validators
├── migrations/                  # Database migrations
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── docs/                        # Swagger documentation (generated)
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
└── README.md
```

## 🔧 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL 15 (if running locally without Docker)
- golang-migrate (optional, for manual migrations)

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone <repository-url>
cd Wallet
```

### 2. Set up environment variables

```bash
cp .env.example .env
```

Edit `.env` with your configuration if needed.

### 3. Run with Docker Compose (Recommended)

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- API server on port 8080

### 4. Access the API

- API Base URL: `http://localhost:8080`
- Swagger Documentation: `http://localhost:8080/swagger/index.html`
- Health Check: `http://localhost:8080/health`

## 🔨 Development Setup

### Install dependencies

```bash
go mod download
```

### Run PostgreSQL locally

```bash
docker run --name wallet_postgres \
  -e POSTGRES_USER=wallet_user \
  -e POSTGRES_PASSWORD=wallet_pass \
  -e POSTGRES_DB=wallet_db \
  -p 5432:5432 \
  -d postgres:15-alpine
```

### Run migrations (optional - auto-migrated on startup)

```bash
make migrate-up
```

### Run the application

```bash
make run
# or
go run cmd/api/main.go
```

### Generate Swagger documentation

```bash
make swagger
```

### Run tests

```bash
make test
```

## 📚 API Documentation

### Authentication

All protected endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <access_token>
```

### API Endpoints

#### Auth

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/me` - Get current user profile (protected)

#### Wallet

- `GET /api/v1/wallet/balance` - Get wallet balance (protected)
- `POST /api/v1/wallet/credit` - Credit wallet (protected)
- `POST /api/v1/wallet/debit` - Debit wallet (protected)

#### Transactions

- `POST /api/v1/transactions/transfer` - Transfer money (protected)
- `GET /api/v1/transactions` - Get transaction history with pagination (protected) — query params: `limit`, `offset`

#### Payments

- `POST /api/v1/boleto` - Generate boleto (protected) — body: `{ "amount": 500.00 }`
- `POST /api/v1/pix` - Create PIX payment (protected) — body: `{ "pix_key": "...", "amount": 100.00 }`. **Nota:** requer que a chave PIX do destinatário já exista na tabela `pix_keys`; o CRUD de chaves PIX via API é trabalho futuro.

### Example Requests

#### Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "cpf": "12345678909",
    "email": "john@example.com",
    "phone": "11987654321",
    "password": "password123"
  }'
```

#### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

#### Get Balance

```bash
curl -X GET http://localhost:8080/api/v1/wallet/balance \
  -H "Authorization: Bearer <access_token>"
```

#### Transfer Money

```bash
curl -X POST http://localhost:8080/api/v1/transactions/transfer \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "to_user_id": "uuid-here",
    "amount": 100.50,
    "description": "Payment for services"
  }'
```

## 🧪 Testing

Run all tests:

```bash
make test
```

Run tests with coverage:

```bash
go test -v -cover ./...
```

## 🐳 Docker Commands

Build and start containers:

```bash
make docker-up
```

Stop containers:

```bash
make docker-down
```

View logs:

```bash
docker-compose logs -f api
```

## 📝 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| APP_ENV | Application environment | development |
| APP_PORT | Server port | 8080 |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | wallet_user |
| DB_PASSWORD | Database password | wallet_pass |
| DB_NAME | Database name | wallet_db |
| DB_SSLMODE | Database SSL mode | disable |
| JWT_SECRET | JWT secret key | your-secret-key |
| JWT_ACCESS_EXPIRY | Access token expiry | 15m |
| JWT_REFRESH_EXPIRY | Refresh token expiry | 7d |
| BCRYPT_COST | Bcrypt hashing cost | 10 |

## 🏗️ Architecture

This project follows **Clean Architecture** principles:

- **Domain Layer**: Core business entities and rules
- **Repository Layer**: Data access abstraction
- **Service Layer**: Business logic implementation
- **Handler Layer**: HTTP request/response handling
- **Middleware**: Cross-cutting concerns (auth, logging)

### Key Design Patterns

- Repository Pattern
- Dependency Injection
- Interface-based design
- Separation of Concerns

## 🔒 Security Features

- Password hashing with bcrypt
- JWT-based authentication
- Refresh token rotation
- CPF, email, and phone validation
- Account blocking mechanism
- SQL injection prevention (GORM)

## 🚧 Future Enhancements

- [ ] Email/SMS notification integration
- [ ] Rate limiting
- [ ] API key authentication for webhooks
- [ ] Transaction rollback mechanism
- [ ] Audit logging
- [ ] Multi-currency support
- [ ] KYC integration
- [ ] 2FA authentication

## 📄 License

MIT License

## 👥 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

For support, email support@wallet.com or open an issue in the repository.

---

Built with ❤️ using Go and Clean Architecture
