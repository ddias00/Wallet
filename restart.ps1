# Reiniciar API
Write-Host "Reiniciando API..." -ForegroundColor Yellow
docker-compose restart api
Start-Sleep -Seconds 5
Write-Host "API reiniciada!" -ForegroundColor Green
