# Скрипт для установки зависимостей
# Использование: .\install-deps.ps1

Write-Host "📦 Установка зависимостей Learn Hub..." -ForegroundColor Green

# Переходим в папку frontend
Set-Location -Path "frontend"

# Устанавливаем зависимости
Write-Host "Установка зависимостей фронтенда..." -ForegroundColor Yellow
npm install

Write-Host "✅ Зависимости установлены!" -ForegroundColor Green

# Возвращаемся в корневую директорию
Set-Location -Path ".."
