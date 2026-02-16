# Установка зависимостей для фронтенда

## Требования

- **Node.js** версии 18.x или выше (рекомендуется 20.x)
- **npm** (обычно устанавливается вместе с Node.js) версии 9.x или выше

## Проверка версий

Перед установкой проверьте версии:

```bash
node --version
npm --version
```

## Установка зависимостей

### Вариант 1: Используя npm (рекомендуется)

```bash
# Перейдите в папку frontend
cd frontend

# Установите все зависимости
npm install
```

### Вариант 2: Используя npm ci (для production/CI)

```bash
# Перейдите в папку frontend
cd frontend

# Установите зависимости точно по package-lock.json (быстрее и надежнее)
npm ci
```

### Вариант 3: Если есть проблемы с npm

```bash
# Перейдите в папку frontend
cd frontend

# Очистите кэш npm
npm cache clean --force

# Удалите node_modules и package-lock.json (если нужно)
rm -rf node_modules package-lock.json

# Установите заново
npm install
```

## После установки

После успешной установки зависимостей вы можете:

1. **Запустить dev-сервер:**
   ```bash
   npm run dev
   ```

2. **Собрать production версию:**
   ```bash
   npm run build
   ```

3. **Предпросмотр production сборки:**
   ```bash
   npm run preview
   ```

## Установленные зависимости

### Основные зависимости (dependencies):
- `@solidjs/router` ^0.15.3 - Роутинг для Solid.js
- `solid-js` ^1.9.7 - Фреймворк Solid.js

### Зависимости для разработки (devDependencies):
- `typescript` ~5.8.3 - TypeScript компилятор
- `vite` ^6.0.0 - Сборщик и dev-сервер
- `vite-plugin-solid` ^2.11.7 - Плагин Vite для Solid.js

## Устранение проблем

### Ошибка "ERESOLVE unable to resolve dependency"
```bash
npm install --legacy-peer-deps
```

### Ошибка "npm ERR! code EACCES"
Используйте sudo (Linux/Mac) или запустите терминал от имени администратора (Windows):
```bash
sudo npm install
```

### Медленная установка
Используйте альтернативный реестр:
```bash
npm install --registry https://registry.npmmirror.com
```

### Проблемы с версией Node.js
Установите правильную версию Node.js через nvm (Node Version Manager):
```bash
# Установка nvm (Linux/Mac)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# Установка Node.js 20
nvm install 20
nvm use 20
```

## Быстрая команда для копирования

```bash
cd frontend && npm install
```
