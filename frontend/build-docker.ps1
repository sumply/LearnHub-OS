# PowerShell скрипт для сборки Docker образа

Write-Host "Сборка Docker образа для LearnHub-OS..." -ForegroundColor Green

# Собираем образ с тегом learnhub-os:latest
docker build -t learnhub-os:latest .

if ($LASTEXITCODE -eq 0) {
    Write-Host "Сборка завершена успешно!" -ForegroundColor Green
    Write-Host "Для запуска контейнера используйте:" -ForegroundColor Yellow
    Write-Host "  docker run -d -p 8080:80 --name learnhub-os learnhub-os:latest" -ForegroundColor Cyan
} else {
    Write-Host "Ошибка при сборке образа!" -ForegroundColor Red
    exit 1
}
