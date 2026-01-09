#!/bin/bash
# Скрипт для сборки Docker образа

echo "Сборка Docker образа для LearnHub-OS..."

# Собираем образ с тегом learnhub-os:latest
docker build -t learnhub-os:latest .

echo "Сборка завершена!"
echo "Для запуска контейнера используйте:"
echo "  docker run -d -p 8080:80 --name learnhub-os learnhub-os:latest"
