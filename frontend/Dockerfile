# Многостадийная сборка для фронтенда
# Стадия 1: Сборка приложения
FROM node:20-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY package*.json ./

# Устанавливаем зависимости
RUN npm ci

# Копируем исходный код
COPY . .

# Принимаем build args для переменных окружения
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

# Собираем приложение
RUN npm run build

# Стадия 2: Production сервер
FROM nginx:alpine

# Копируем собранные файлы из стадии сборки
COPY --from=builder /app/dist /usr/share/nginx/html

# Копируем конфигурацию nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Открываем порт 80
EXPOSE 80

# Запускаем nginx
CMD ["nginx", "-g", "daemon off;"]
