# Script para rodar testes e coverage
Write-Host "=== Executando Testes da Digital Wallet API ===" -ForegroundColor Green
Write-Host ""

# 1. Rodar todos os testes
Write-Host "1. Rodando todos os testes..." -ForegroundColor Yellow
$env:CGO_ENABLED="1"
go test ./... -v

Write-Host ""
Write-Host "=== Gerando Coverage Report ===" -ForegroundColor Green
Write-Host ""

# 2. Rodar testes com coverage
Write-Host "2. Gerando coverage..." -ForegroundColor Yellow
go test ./... -coverprofile=coverage.out -covermode=atomic

Write-Host ""

# 3. Mostrar coverage por pacote
Write-Host "3. Coverage por pacote:" -ForegroundColor Yellow
go test ./... -cover

Write-Host ""

# 4. Gerar relatório HTML
Write-Host "4. Gerando relatório HTML..." -ForegroundColor Yellow
if (Test-Path coverage.out) {
    go tool cover -html=coverage.out -o coverage.html
    Write-Host ""
    Write-Host "=== Relatórios Gerados ===" -ForegroundColor Green
    Write-Host "- coverage.out (dados brutos)" -ForegroundColor Cyan
    Write-Host "- coverage.html (relatório visual)" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Abra coverage.html no navegador para ver o relatório detalhado!" -ForegroundColor Green

    # 5. Abrir relatório automaticamente
    Write-Host ""
    $response = Read-Host "Deseja abrir o relatório HTML agora? (S/N)"
    if ($response -eq "S" -or $response -eq "s") {
        Start-Process (Resolve-Path coverage.html)
    }
} else {
    Write-Host "Arquivo coverage.out não foi gerado" -ForegroundColor Red
}
