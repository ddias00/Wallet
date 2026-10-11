# Digital Wallet API

API REST completa para uma plataforma de carteira digital, construída em Go seguindo os princípios da Clean Architecture.

## Features

### Gestão de Usuários
- Cadastro com validação de CPF, e-mail e telefone
- Autenticação JWT com refresh tokens
- Bloqueio/desbloqueio de conta

### Operações de Carteira
- Consulta de saldo
- Crédito
- Débito

### Transações
- Transferência entre usuários
- Histórico paginado
- Status (pending, confirmed, cancelled)

### Pagamentos
- Geração de boleto
- Pagamento via PIX
- QR Code

### Notificações
- Notificações de transação (console-based, pronto para integração com e-mail/SMS)

## Tech Stack

- **Linguagem**: Go 1.27.2
- **Framework HTTP**: Gin
- **Banco de dados**: PostgreSQL 15
- **ORM**: GORM
- **Autenticação**: JWT com refresh tokens
- **Hash de senhas**: bcrypt
- **Documentação**: Swagger/OpenAPI
- **Containerização**: Docker & Docker Compose
- **Hot reload (dev)**: Air
- **Migrations**: compatível com golang-migrate

## Estrutura do Projeto

```
Wallet/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── config/                  # Configuração
│   ├── domain/                  # Entidades de domínio
│   ├── handler/                 # Handlers HTTP
│   ├── service/                 # Regras de negócio
│   ├── repository/              # Camada de acesso a dados
│   ├── middleware/              # Auth, logging
│   └── utils/                   # Helpers
├── pkg/
│   ├── jwt/                     # Utilitários JWT
│   └── validator/               # Validadores (CPF, e-mail, telefone)
├── migrations/                  # Migrations SQL
├── docs/                        # Documentação Swagger (gerada)
├── scripts/                     # Scripts utilitários, smoke tests, demo e Postman collection
├── test/                        # Testes de integração
├── docker-compose.yml
├── docker-compose.dev.yml       # Stack com hot reload
├── Dockerfile
├── Dockerfile.dev
├── Makefile
├── go.mod
├── .env.example
└── README.md
```

## Pré-requisitos

- Go 1.27.2 ou superior
- Docker e Docker Compose
- PostgreSQL 15 (se rodar localmente sem Docker)
- golang-migrate (opcional, para migrations manuais)

## Quick Start

### 1. Clone o repositório

```bash
git clone <repository-url>
cd Wallet
```

### 2. Configure as variáveis de ambiente

```bash
cp .env.example .env
```

Edite o `.env` conforme necessário.

### 3. Suba tudo com Docker Compose (recomendado)

```bash
docker compose up -d
```

Isso vai subir:
- PostgreSQL na porta `5432`
- API na porta `8080`

### 4. Acesse a API

- Base URL: `http://localhost:8080`
- Swagger: `http://localhost:8080/swagger/index.html`
- Health check: `http://localhost:8080/health`

## Desenvolvimento

### Instalar dependências

```bash
go mod download
```

### Rodar PostgreSQL localmente (sem docker-compose)

```bash
docker run --name wallet_postgres \
  -e POSTGRES_USER=wallet_user \
  -e POSTGRES_PASSWORD=wallet_pass \
  -e POSTGRES_DB=wallet_db \
  -p 5432:5432 \
  -d postgres:15-alpine
```

### Rodar a aplicação

```bash
make run
# ou
go run cmd/api/main.go
```

### Hot reload com Air (modo dev)

```bash
make docker-dev
```

### Gerar documentação Swagger

```bash
make swagger
```

### Rodar testes

```bash
make test            # todos os testes
make coverage        # gera coverage.html
make coverage-func   # mostra coverage no terminal
```

## Scripts e ferramentas

Estão em `scripts/` (entre por lá pra ver detalhes). Resumo:

| Arquivo | Plataforma | O que faz |
|---------|-----------|-----------|
| `setup.bat` / `setup.sh` | Windows / Linux+Mac | Checa Go, baixa deps, valida Docker |
| `test_api.ps1` / `test_api.sh` | Windows / Linux+Mac | Smoke test end-to-end da API (register → login → balance → credit) |
| `reset-demo.ps1` | Windows | Reseta banco, cria Alice e Bob, semeia chave PIX (TCC) |
| `postman_collection.json` | Postman / Insomnia / Bruno | Coleção completa dos endpoints — `File → Import` no Postman |

## Documentação da API

### Autenticação

Endpoints protegidos exigem o header:

```
Authorization: Bearer <access_token>
```

### Endpoints

#### Auth

- `POST /api/v1/auth/register` — cadastra novo usuário
- `POST /api/v1/auth/login` — autentica
- `POST /api/v1/auth/refresh` — renova o access token
- `GET /api/v1/auth/me` — perfil do usuário (protected)

#### Wallet

- `GET /api/v1/wallet/balance` — consulta saldo (protected)
- `POST /api/v1/wallet/credit` — credita (protected)
- `POST /api/v1/wallet/debit` — debita (protected)

#### Transactions

- `POST /api/v1/transactions/transfer` — transferência (protected)
- `GET /api/v1/transactions` — histórico paginado (protected) — query params: `limit`, `offset`

#### Payments

- `POST /api/v1/boleto` — gera boleto (protected) — body: `{ "amount": 500.00 }`
- `POST /api/v1/pix` — paga via PIX (protected) — body: `{ "pix_key": "...", "amount": 100.00 }`. **Nota:** requer que a chave PIX do destinatário já exista na tabela `pix_keys`; o CRUD de chaves PIX via API é trabalho futuro.

### Exemplos

#### Register

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

#### Balance

```bash
curl -X GET http://localhost:8080/api/v1/wallet/balance \
  -H "Authorization: Bearer <access_token>"
```

#### Transferência

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

## Docker

```bash
make docker-up        # sobe API + Postgres
make docker-down      # derruba
make docker-logs      # logs da API
make docker-rebuild   # rebuild + restart
make docker-dev       # stack com hot reload
```

## Variáveis de ambiente

| Variável | Descrição | Default |
|----------|-----------|---------|
| `APP_ENV` | Ambiente | `development` |
| `APP_PORT` | Porta do servidor | `8080` |
| `DB_HOST` | Host do banco | `localhost` |
| `DB_PORT` | Porta do banco | `5432` |
| `DB_USER` | Usuário do banco | `wallet_user` |
| `DB_PASSWORD` | Senha do banco | `wallet_pass` |
| `DB_NAME` | Nome do banco | `wallet_db` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `JWT_SECRET` | Segredo do JWT | `your-secret-key` |
| `JWT_ACCESS_EXPIRY` | Expiração do access token | `15m` |
| `JWT_REFRESH_EXPIRY` | Expiração do refresh token | `7d` |
| `BCRYPT_COST` | Custo do bcrypt | `10` |

## Arquitetura

Projeto segue **Clean Architecture**:

- **Domain Layer** — entidades e regras de negócio centrais
- **Repository Layer** — abstração de acesso a dados
- **Service Layer** — implementação das regras de negócio
- **Handler Layer** — recepção/resposta HTTP
- **Middleware** — autenticação, logging e demais cross-cutting concerns

### Padrões aplicados

- Repository Pattern
- Dependency Injection
- Programação por interfaces
- Separação de responsabilidades

## Segurança

- Hash de senha com bcrypt
- Autenticação JWT
- Rotação de refresh token
- Validação de CPF, e-mail e telefone
- Mecanismo de bloqueio de conta
- Prevenção de SQL injection (GORM)

## Roadmap

- [ ] Integração com e-mail/SMS para notificações
- [ ] Rate limiting
- [ ] API key para webhooks
- [ ] Mecanismo de rollback de transações
- [ ] Audit logging
- [ ] Multi-currency
- [ ] Integração com KYC
- [ ] 2FA

## Licença

MIT License

## Contribuindo

1. Fork
2. Crie sua branch (`git checkout -b feature/amazing-feature`)
3. Commit (`git commit -m 'feat: amazing feature'`)
4. Push (`git push origin feature/amazing-feature`)
5. Abra um Pull Request

---

Construído com Go e Clean Architecture.
