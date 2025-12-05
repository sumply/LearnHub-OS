# API Спецификация для Бэкенда

## Базовый URL
```
http://188.225.24.208:8000
```

## Аутентификация

### 1. Регистрация пользователя
**Эндпоинт:** `POST /registration`

**Заголовки:**
```
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "email": "string",
  "password": "string",
  "name": "string",
  "surname": "string",
  "role": "student" | "parent" | "teacher",
  "group": "string" // опционально
}
```

**Ответ:**
```json
{
  "success": true,
  "message": "string",
  "data": {
    "user": {
      "id": "string",
      "email": "string",
      "name": "string",
      "surname": "string",
      "role": "student" | "parent" | "teacher",
      "group": "string"
    },
    "token": "string",
    "refresh_token": "string"
  }
}
```

### 2. Авторизация пользователя
**Эндпоинт:** `POST /authorization`

**Заголовки:**
```
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Ответ:**
```json
{
  "success": true,
  "message": "string",
  "token": "string",
  "refresh_token": "string",
  "user": {
    "id": "string",
    "email": "string",
    "name": "string",
    "surname": "string",
    "role": "student" | "parent" | "teacher",
    "group": "string"
  }
}
```

### 3. Обновление access-токена
**Эндпоинт:** `POST /refresh_access_token`

**Заголовки:**
```
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "refresh_token": "string"
}
```

**Ответ:**
```json
{
  "success": true,
  "message": "string",
  "token": "string"
}
```

### 4. Проверка токена
**Эндпоинт:** `GET /profile`

**Заголовки:**
```
Authorization: Bearer {token}
```

**Ответ:**
```json
{
  "success": true,
  "message": "string",
  "user": {
    "id": "string",
    "email": "string",
    "name": "string",
    "surname": "string",
    "role": "student" | "parent" | "teacher",
    "group": "string"
  }
}
```

## Материалы

### 5. Получение материалов
**Эндпоинт:** `GET /materials`

**Заголовки:**
```
Authorization: Bearer {token}
```

**Ответ:**
```json
{
  "success": true,
  "message": "string",
  "materials": [
    {
      "id": "string",
      "title": "string",
      "description": "string",
      "attachments": [
        {
          "type": "string",
          "url": "string",
          "testId": "string"
        }
      ],
      "category": "string",
      "visible": true
    }
  ]
}
```

## Задания

### 6. Получение заданий
**Эндпоинт:** `GET /tasks`

**Заголовки:**
```
Authorization: Bearer {token}
```

**Ответ:**
```json
{
  "success": true,
  "message": "string",
  "tasks": [
    {
      "id": "string",
      "title": "string",
      "description": "string",
      "attachments": [
        {
          "type": "string",
          "url": "string",
          "testId": "string"
        }
      ],
      "category": "string",
      "visible": true
    }
  ]
}
```

## Карточки (Flashcards)

### 7. Получение карточек
**Эндпоинт:** `GET /flashcards`

**Заголовки:**
```
Authorization: Bearer {token}
```

**Ответ:**
```json
{
  "success": true,
  "message": "string",
  "flashcards": [
    {
      "id": "string",
      "question": "string",
      "answer": "string",
      "category": "string"
    }
  ]
}
```

## Тесты (Quizzes)

### 8. Получение тестов
**Эндпоинт:** `GET /quizzes`

**Заголовки:**
```
Authorization: Bearer {token}
```

**Ответ:**
```json
{
  "success": true,
  "message": "string",
  "quizzes": [
    {
      "id": "string",
      "title": "string",
      "questions": [
        {
          "question": "string",
          "options": ["string"],
          "correct": 0
        }
      ],
      "category": "string"
    }
  ]
}
```

## Общие требования

### Коды ответов HTTP
- `200` - Успешный запрос
- `401` - Неавторизованный доступ (неверный токен)
- `400` - Ошибка в запросе
- `500` - Внутренняя ошибка сервера

### Структура ответа
Все ответы должны содержать поля:
- `success`: boolean - статус выполнения операции
- `message`: string - сообщение о результате
- `data` или конкретные поля с данными (опционально)

### Аутентификация
Для защищенных эндпоинтов требуется заголовок:
```
Authorization: Bearer {access_token}
```

### Обработка ошибок
При ошибке аутентификации (401) сервер должен возвращать:
```json
{
  "success": false,
  "message": "Unauthorized access"
}
```

## Типы данных

### User
```typescript
interface User {
  id: string;
  email: string;
  name: string;
  surname: string;
  role: 'student' | 'parent' | 'teacher';
  group?: string;
}
```

### MaterialAttachment
```typescript
interface MaterialAttachment {
  type: string;
  url?: string;
  testId?: string;
}
```

### Material/Task
```typescript
interface Material {
  id: string;
  title: string;
  description: string;
  attachments: MaterialAttachment[];
  category: string;
  visible: boolean;
}
```

### Flashcard
```typescript
interface Flashcard {
  id: string;
  question: string;
  answer: string;
  category: string;
}
```

### Quiz
```typescript
interface Quiz {
  id: string;
  title: string;
  questions: Array<{
    question: string;
    options: string[];
    correct: number;
  }>;
  category: string;
}
``` 