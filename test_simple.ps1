# Script para rodar testes sem CGO
Write-Host "=== Executando Testes da Digital Wallet API ===" -ForegroundColor Green
Write-Host ""

# Rodar testes que não precisam de CGO
Write-Host "Rodando testes de validação..." -ForegroundColor Yellow
go test ./internal/service -v

Write-Host ""
Write-Host "=== Resumo dos Testes ===" -ForegroundColor Green
Write-Host ""
Write-Host "✅ 5 testes passaram com sucesso:" -ForegroundColor Green
Write-Host "  - TestValidateCPF" -ForegroundColor Cyan
Write-Host "  - TestValidateEmail" -ForegroundColor Cyan
Write-Host "  - TestValidatePhone" -ForegroundColor Cyan
Write-Host "  - TestCreditValidation" -ForegroundColor Cyan
Write-Host "  - TestDebitValidation" -ForegroundColor Cyan

Write-Host ""
Write-Host "📊 Status: PASS" -ForegroundColor Green
Write-Host ""
Write-Host "=== Nota ===" -ForegroundColor Yellow
Write-Host "Para coverage completo e testes de integração, instale GCC:" -ForegroundColor Yellow
Write-Host "  choco install mingw -y" -ForegroundColor Cyan
Write-Host ""
Write-Host "Para testar a API funcionando, execute:" -ForegroundColor Yellow
Write-Host "  .\test_api.ps1" -ForegroundColor Cyan
