# Learn Hub OS

Образовательная платформа нового поколения для учеников, учителей, родителей и администраторов.

## Структура проекта

```
LearnHub-OS/
├── frontend/          # Frontend приложение (Solid.js + Vite)
└── package.json      # Корневой package.json для удобного запуска
```

## Быстрый старт

### Установка зависимостей

```bash
# Установить зависимости фронтенда
npm run install:frontend

# Или перейти в папку frontend и установить там
cd frontend
npm install
```

### Запуск проекта

Из корневой директории проекта:

```bash
# Запустить dev-сервер
npm run dev

# Собрать production версию
npm run build

# Предпросмотр production сборки
npm run preview
```

Или из папки `frontend`:

```bash
cd frontend
npm run dev
```

## Дополнительная информация

- Подробная документация по фронтенду: [frontend/README.md](frontend/README.md)
- Документация по адаптивному дизайну: [frontend/RESPONSIVE_DESIGN.md](frontend/RESPONSIVE_DESIGN.md)
- Инструкции по установке: [frontend/INSTALL.md](frontend/INSTALL.md)
- Docker инструкции: [frontend/README_DOCKER.md](frontend/README_DOCKER.md)
