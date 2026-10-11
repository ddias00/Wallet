<#
.SYNOPSIS
    Reseta o ambiente da demo: limpa o banco, cadastra Alice e Bob,
    semeia a chave PIX do Bob e imprime tudo pronto pra usar no Postman.

.DESCRIPTION
    Pre-condicoes:
      - Docker rodando (wallet_api + wallet_postgres)
      - API respondendo em http://localhost:8080

    Por padrao a Alice comeca com saldo ZERO (pra demo do credito ao vivo).
    Use -CreditAlice <valor> pra ja entrar com saldo (acelera a demo).

.PARAMETER CreditAlice
    Valor (em R$) para creditar a Alice automaticamente. 0 = nao credita.

.PARAMETER NoClean
    Pula a etapa de TRUNCATE (caso voce queira reaproveitar usuarios existentes).

.EXAMPLE
    .\reset-demo.ps1
    # Banco limpo, Alice e Bob criados, saldo da Alice = R$ 0

.EXAMPLE
    .\reset-demo.ps1 -CreditAlice 1000
    # Banco limpo, Alice e Bob criados, Alice ja com R$ 1000 de saldo

.NOTES
    Diego Alves Dias - TCC II - Engenharia de Software - UNIGRAN 2026
#>
[CmdletBinding()]
param(
    [double]$CreditAlice = 0,
    [switch]$NoClean
)

$ErrorActionPreference = 'Stop'
$base = 'http://localhost:8080'
$api  = "$base/api/v1"

# Encoding UTF-8 no terminal pra acentos nao virarem '?'
try { [Console]::OutputEncoding = [System.Text.Encoding]::UTF8 } catch {}

function Write-Step($msg) { Write-Host "`n>>> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "[OK] $msg" -ForegroundColor Green }
function Write-Warn($msg) { Write-Host "[!]  $msg" -ForegroundColor Yellow }
function Write-Err($msg)  { Write-Host "[X]  $msg" -ForegroundColor Red }

# --------------------------------------------------------------
# 0. Pre-checks
# --------------------------------------------------------------
Write-Step 'Pre-checks (Docker + API)'

$containers = docker ps --format '{{.Names}}' 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Err 'Docker nao esta rodando. Inicie o Docker Desktop e tente de novo.'
    exit 1
}
if ($containers -notcontains 'wallet_postgres') {
    Write-Err "Container 'wallet_postgres' nao encontrado. Rode: docker compose up -d"
    exit 1
}
if ($containers -notcontains 'wallet_api') {
    Write-Err "Container 'wallet_api' nao encontrado. Rode: docker compose up -d"
    exit 1
}

try {
    $h = Invoke-RestMethod -Method GET -Uri "$base/health" -TimeoutSec 5
    if ($h.status -ne 'ok') { throw 'health nao retornou ok' }
    Write-Ok 'API respondendo em /health'
} catch {
    Write-Err "API nao respondeu em $base/health . Detalhe: $($_.Exception.Message)"
    exit 1
}

# --------------------------------------------------------------
# 1. Limpar banco (opcional)
# --------------------------------------------------------------
if (-not $NoClean) {
    Write-Step 'Limpando banco (TRUNCATE em todas as tabelas)'
    $sql = @'
TRUNCATE TABLE refresh_tokens, notifications, boletos, pix_keys, transactions, wallets, users RESTART IDENTITY CASCADE;
'@
    docker exec wallet_postgres psql -U wallet_user -d wallet_db -c $sql | Out-Null
    if ($LASTEXITCODE -ne 0) { Write-Err 'Falha ao truncar tabelas'; exit 1 }
    Write-Ok 'Banco zerado'
} else {
    Write-Warn 'Pulando limpeza do banco (-NoClean)'
}

# --------------------------------------------------------------
# 2. Funcao auxiliar: gerar CPF valido aleatorio
# --------------------------------------------------------------
function New-ValidCPF {
    do {
        $first9 = -join (1..9 | ForEach-Object { Get-Random -Maximum 10 })
        $sum = 0
        for ($i=0; $i -lt 9; $i++) { $sum += [int]$first9[$i].ToString() * (10 - $i) }
        $d1 = ($sum * 10) % 11; if ($d1 -eq 10) { $d1 = 0 }
        $ten = $first9 + $d1
        $sum = 0
        for ($i=0; $i -lt 10; $i++) { $sum += [int]$ten[$i].ToString() * (11 - $i) }
        $d2 = ($sum * 10) % 11; if ($d2 -eq 10) { $d2 = 0 }
        $cpf = $ten + $d2
    } while ($cpf -match '^(\d)\1{10}$')
    return $cpf
}

# --------------------------------------------------------------
# 3. Cadastrar Alice e Bob
# --------------------------------------------------------------
Write-Step 'Cadastrando Alice'
$aliceCpf = New-ValidCPF
$aliceBody = @{
    name = 'Alice Silva'; cpf = $aliceCpf; email = 'alice@wallet.com'
    phone = '11987654321'; password = 'Senha@123'
} | ConvertTo-Json
$aliceResp = Invoke-RestMethod -Method POST -Uri "$api/auth/register" -ContentType 'application/json' -Body $aliceBody
$aliceId = $aliceResp.data.id
Write-Ok "Alice cadastrada (UUID: $aliceId, CPF: $aliceCpf)"

Write-Step 'Cadastrando Bob'
$bobCpf = New-ValidCPF
$bobBody = @{
    name = 'Bob Costa'; cpf = $bobCpf; email = 'bob@wallet.com'
    phone = '11912345678'; password = 'Senha@123'
} | ConvertTo-Json
$bobResp = Invoke-RestMethod -Method POST -Uri "$api/auth/register" -ContentType 'application/json' -Body $bobBody
$bobId = $bobResp.data.id
Write-Ok "Bob cadastrado (UUID: $bobId, CPF: $bobCpf)"

# --------------------------------------------------------------
# 4. Seed da chave PIX do Bob
# --------------------------------------------------------------
Write-Step 'Semeando chave PIX do Bob (email)'
$seedSql = "INSERT INTO pix_keys (id, user_id, key, type, status, created_at, updated_at) VALUES (gen_random_uuid(), '$bobId', 'bob@wallet.com', 'email', 'active', NOW(), NOW());"
docker exec wallet_postgres psql -U wallet_user -d wallet_db -c $seedSql | Out-Null
if ($LASTEXITCODE -ne 0) { Write-Err 'Falha ao semear chave PIX'; exit 1 }
Write-Ok "Chave PIX 'bob@wallet.com' cadastrada para o Bob"

# --------------------------------------------------------------
# 5. Login Alice + Bob (capturar tokens)
# --------------------------------------------------------------
Write-Step 'Fazendo login da Alice'
$loginAlice = Invoke-RestMethod -Method POST -Uri "$api/auth/login" -ContentType 'application/json' -Body (@{ email = 'alice@wallet.com'; password = 'Senha@123' } | ConvertTo-Json)
$tokenAlice  = $loginAlice.data.access_token
$refreshAlice= $loginAlice.data.refresh_token
Write-Ok 'Token da Alice obtido'

Write-Step 'Fazendo login do Bob'
$loginBob = Invoke-RestMethod -Method POST -Uri "$api/auth/login" -ContentType 'application/json' -Body (@{ email = 'bob@wallet.com'; password = 'Senha@123' } | ConvertTo-Json)
$tokenBob  = $loginBob.data.access_token
$refreshBob= $loginBob.data.refresh_token
Write-Ok 'Token do Bob obtido'

# --------------------------------------------------------------
# 6. Credito opcional na Alice
# --------------------------------------------------------------
if ($CreditAlice -gt 0) {
    Write-Step "Creditando R$ $CreditAlice na Alice"
    Invoke-RestMethod -Method POST -Uri "$api/wallet/credit" `
        -Headers @{ Authorization = "Bearer $tokenAlice" } `
        -ContentType 'application/json' `
        -Body (@{ amount = $CreditAlice; description = 'Carga inicial (reset-demo.ps1)' } | ConvertTo-Json) | Out-Null
    Write-Ok "Alice creditada com R$ $CreditAlice"
}

# --------------------------------------------------------------
# 7. Verificar saldos finais
# --------------------------------------------------------------
$saldoAlice = (Invoke-RestMethod -Method GET -Uri "$api/wallet/balance" -Headers @{ Authorization = "Bearer $tokenAlice" }).data.balance
$saldoBob   = (Invoke-RestMethod -Method GET -Uri "$api/wallet/balance" -Headers @{ Authorization = "Bearer $tokenBob" }).data.balance

# --------------------------------------------------------------
# 8. Sumario final (copiavel)
# --------------------------------------------------------------
Write-Host ""
Write-Host ("=" * 70) -ForegroundColor DarkGray
Write-Host "  AMBIENTE DA DEMO PRONTO" -ForegroundColor Green
Write-Host ("=" * 70) -ForegroundColor DarkGray
Write-Host ""
Write-Host "ALICE" -ForegroundColor Cyan
Write-Host "  UUID         : $aliceId"
Write-Host "  Email        : alice@wallet.com"
Write-Host "  CPF          : $aliceCpf"
Write-Host "  Senha        : Senha@123"
Write-Host "  Saldo        : R$ $saldoAlice"
Write-Host "  AccessToken  : $tokenAlice"
Write-Host ""
Write-Host "BOB" -ForegroundColor Cyan
Write-Host "  UUID         : $bobId"
Write-Host "  Email        : bob@wallet.com"
Write-Host "  CPF          : $bobCpf"
Write-Host "  Senha        : Senha@123"
Write-Host "  Chave PIX    : bob@wallet.com (tipo email, ativa)"
Write-Host "  Saldo        : R$ $saldoBob"
Write-Host "  AccessToken  : $tokenBob"
Write-Host ""
Write-Host ("-" * 70) -ForegroundColor DarkGray
Write-Host "POSTMAN: cole os valores acima nas Collection Variables:" -ForegroundColor Yellow
Write-Host "  accessToken  -> token da Alice (acima)"
Write-Host "  refreshToken -> $refreshAlice"
Write-Host "  aliceId      -> $aliceId"
Write-Host "  bobId        -> $bobId"
Write-Host "  bobEmail     -> bob@wallet.com"
Write-Host ("-" * 70) -ForegroundColor DarkGray
Write-Host ""
Write-Host "DICA: rode '.\reset-demo.ps1 -CreditAlice 1000' pra ja entrar com saldo." -ForegroundColor DarkGray
Write-Host ""

# Tambem exportar para variaveis de ambiente da sessao (caso queira usar em curl)
$env:DEMO_ALICE_ID = $aliceId
$env:DEMO_ALICE_TOKEN = $tokenAlice
$env:DEMO_BOB_ID = $bobId
$env:DEMO_BOB_TOKEN = $tokenBob
