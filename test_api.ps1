# Script de teste para Digital Wallet API
Write-Host "=== Testando Digital Wallet API ===" -ForegroundColor Green
Write-Host ""

# 1. Health Check
Write-Host "1. Health Check:" -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get
    $health | ConvertTo-Json
} catch {
    Write-Host "Falhou: $_" -ForegroundColor Red
}
Write-Host ""

# 2. Registrar usuário
Write-Host "2. Registrando usuário:" -ForegroundColor Yellow
$registerBody = @{
    name = "Teste User"
    cpf = "12345678909"
    email = "teste@example.com"
    phone = "11987654321"
    password = "senha123"
} | ConvertTo-Json

try {
    $register = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/register" `
        -Method Post `
        -ContentType "application/json" `
        -Body $registerBody
    $register | ConvertTo-Json
} catch {
    Write-Host "Usuário já existe ou erro: $_" -ForegroundColor Yellow
}
Write-Host ""

# 3. Login
Write-Host "3. Fazendo login:" -ForegroundColor Yellow
$loginBody = @{
    email = "teste@example.com"
    password = "senha123"
} | ConvertTo-Json

try {
    $login = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" `
        -Method Post `
        -ContentType "application/json" `
        -Body $loginBody
    $login | ConvertTo-Json
    $token = $login.data.access_token
    Write-Host "Token obtido com sucesso!" -ForegroundColor Green
} catch {
    Write-Host "Falhou: $_" -ForegroundColor Red
    exit
}
Write-Host ""

# 4. Ver perfil
Write-Host "4. Verificando perfil:" -ForegroundColor Yellow
try {
    $headers = @{
        Authorization = "Bearer $token"
    }
    $profile = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/me" `
        -Method Get `
        -Headers $headers
    $profile | ConvertTo-Json
} catch {
    Write-Host "Falhou: $_" -ForegroundColor Red
}
Write-Host ""

# 5. Ver saldo
Write-Host "5. Consultando saldo:" -ForegroundColor Yellow
try {
    $balance = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/wallet/balance" `
        -Method Get `
        -Headers $headers
    $balance | ConvertTo-Json
} catch {
    Write-Host "Falhou: $_" -ForegroundColor Red
}
Write-Host ""

# 6. Creditar carteira
Write-Host "6. Creditando R$ 100,00:" -ForegroundColor Yellow
$creditBody = @{
    amount = 100.00
    description = "Teste de crédito"
} | ConvertTo-Json

try {
    $credit = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/wallet/credit" `
        -Method Post `
        -ContentType "application/json" `
        -Headers $headers `
        -Body $creditBody
    $credit | ConvertTo-Json
} catch {
    Write-Host "Falhou: $_" -ForegroundColor Red
}
Write-Host ""

# 7. Ver saldo atualizado
Write-Host "7. Saldo atualizado:" -ForegroundColor Yellow
try {
    $balance = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/wallet/balance" `
        -Method Get `
        -Headers $headers
    $balance | ConvertTo-Json
} catch {
    Write-Host "Falhou: $_" -ForegroundColor Red
}
Write-Host ""

Write-Host "=== Testes concluídos ===" -ForegroundColor Green
