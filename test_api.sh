#!/bin/bash

echo "=== Testando Digital Wallet API ==="
echo ""

# 1. Health Check
echo "1. Health Check:"
curl -s http://localhost:8080/health | jq '.' || echo "Falhou"
echo ""

# 2. Registrar usuário
echo "2. Registrando usuário:"
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Teste User",
    "cpf": "12345678909",
    "email": "teste@example.com",
    "phone": "11987654321",
    "password": "senha123"
  }')
echo $REGISTER_RESPONSE | jq '.'
echo ""

# 3. Login
echo "3. Fazendo login:"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teste@example.com",
    "password": "senha123"
  }')
echo $LOGIN_RESPONSE | jq '.'

# Extrair token
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')
echo "Token: $TOKEN"
echo ""

# 4. Ver perfil
echo "4. Verificando perfil:"
curl -s -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# 5. Ver saldo
echo "5. Consultando saldo:"
curl -s -X GET http://localhost:8080/api/v1/wallet/balance \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# 6. Creditar carteira
echo "6. Creditando R$ 100,00:"
curl -s -X POST http://localhost:8080/api/v1/wallet/credit \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.00,
    "description": "Teste de crédito"
  }' | jq '.'
echo ""

# 7. Ver saldo atualizado
echo "7. Saldo atualizado:"
curl -s -X GET http://localhost:8080/api/v1/wallet/balance \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

echo "=== Testes concluídos ==="
