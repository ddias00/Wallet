# Tests

Todos os testes do projeto estão organizados aqui.

## Estrutura

```
test/
├── unit/
│   └── unit_test.go           # Testes unitários (service, repository, auth)
└── integration/
    └── integration_test.go    # Testes de integração (API E2E)
```

## Executar

```bash
# Todos os testes
make test

# Apenas unitários
go test -v ./test/unit

# Apenas integração
go test -v ./test/integration

# Com coverage
make test-coverage
```

## Coverage

```bash
make test-coverage
open coverage.html
```

Veja [TESTING.md](../TESTING.md) para documentação completa.
