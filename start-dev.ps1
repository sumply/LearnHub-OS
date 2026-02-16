# Скрипт для запуска dev-сервера фронтенда
# Использование: .\start-dev.ps1

Write-Host "🚀 Запуск dev-сервера Learn Hub..." -ForegroundColor Green

# Переходим в папку frontend
Set-Location -Path "frontend"

# Проверяем наличие node_modules
if (-not (Test-Path "node_modules")) {
    Write-Host "📦 Установка зависимостей..." -ForegroundColor Yellow
    npm install
}

# Запускаем dev-сервер
Write-Host "✨ Запуск Vite dev-сервера..." -ForegroundColor Cyan
npm run dev
