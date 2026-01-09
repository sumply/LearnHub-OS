# API Спецификация для Системы Авторизации

## Базовый URL
```
http://localhost:3000/api/auth
```

## 1. Регистрация пользователя

### Endpoint
```
POST /api/auth/register
```

### Заголовки
```
Content-Type: application/json
```

### Тело запроса (JSON)
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "Иван",
  "surname": "Иванов",
  "role": "student",
  "group": "11А"
}
```

### Поля:
- `email` (string, обязательное) - email пользователя
- `password` (string, обязательное, минимум 6 символов) - пароль
- `name` (string, обязательное) - имя пользователя
- `surname` (string, обязательное) - фамилия пользователя
- `role` (string, обязательное) - роль: "student", "parent", "teacher"
- `group` (string, опциональное) - класс/группа (только для student и parent)

### Успешный ответ (200 OK)
```json
{
  "success": true,
  "message": "Пользователь успешно зарегистрирован",
  "user": {
    "id": "user_id_123",
    "email": "user@example.com",
    "name": "Иван",
    "surname": "Иванов",
    "role": "student",
    "group": "11А"
  }
}
```

### Ошибки:
- **400 Bad Request** - Неверные данные
```json
{
  "success": false,
  "message": "Email уже используется"
}
```

- **422 Unprocessable Entity** - Ошибки валидации
```json
{
  "success": false,
  "message": "Пароль должен содержать минимум 6 символов"
}
```

## 2. Вход пользователя

### Endpoint
```
POST /api/auth/login
```

### Заголовки
```
Content-Type: application/json
```

### Тело запроса (JSON)
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Поля:
- `email` (string, обязательное) - email пользователя
- `password` (string, обязательное) - пароль

### Успешный ответ (200 OK)
```json
{
  "success": true,
  "message": "Вход выполнен успешно",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user_id_123",
    "email": "user@example.com",
    "name": "Иван",
    "surname": "Иванов",
    "role": "student",
    "group": "11А"
  }
}
```

### Ошибки:
- **401 Unauthorized** - Неверные учетные данные
```json
{
  "success": false,
  "message": "Неверный email или пароль"
}
```

- **400 Bad Request** - Неверные данные
```json
{
  "success": false,
  "message": "Email и пароль обязательны"
}
```

## 3. Выход пользователя

### Endpoint
```
POST /api/auth/logout
```

### Заголовки
```
Authorization: Bearer <token>
Content-Type: application/json
```

### Успешный ответ (200 OK)
```json
{
  "success": true,
  "message": "Выход выполнен успешно"
}
```

## 4. Проверка токена

### Endpoint
```
GET /api/auth/verify
```

### Заголовки
```
Authorization: Bearer <token>
```

### Успешный ответ (200 OK)
```json
{
  "success": true,
  "message": "Токен действителен",
  "user": {
    "id": "user_id_123",
    "email": "user@example.com",
    "name": "Иван",
    "surname": "Иванов",
    "role": "student",
    "group": "11А"
  }
}
```

### Ошибки:
- **401 Unauthorized** - Недействительный токен
```json
{
  "success": false,
  "message": "Недействительный токен"
}
```

## 5. Получение профиля пользователя

### Endpoint
```
GET /api/auth/profile
```

### Заголовки
```
Authorization: Bearer <token>
```

### Успешный ответ (200 OK)
```json
{
  "success": true,
  "user": {
    "id": "user_id_123",
    "email": "user@example.com",
    "name": "Иван",
    "surname": "Иванов",
    "role": "student",
    "group": "11А",
    "createdAt": "2024-01-01T00:00:00.000Z",
    "updatedAt": "2024-01-01T00:00:00.000Z"
  }
}
```

## Общие требования к бэкенду:

### 1. Валидация данных
- Email должен быть валидным форматом
- Пароль минимум 6 символов
- Все обязательные поля должны быть заполнены
- Email должен быть уникальным в системе

### 2. Безопасность
- Пароли должны хешироваться (bcrypt, argon2)
- JWT токены для аутентификации
- Время жизни токена: 24 часа
- Refresh токены (опционально)

### 3. Обработка ошибок
- Всегда возвращать JSON с полями `success` и `message`
- Использовать соответствующие HTTP статус коды
- Логировать ошибки на сервере

### 4. База данных
Создать таблицу `users` со следующими полями:
```sql
CREATE TABLE users (
  id VARCHAR(255) PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  surname VARCHAR(255) NOT NULL,
  role ENUM('student', 'parent', 'teacher') NOT NULL,
  group_name VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 5. Middleware для аутентификации
Создать middleware для проверки JWT токена:
```javascript
const authenticateToken = (req, res, next) => {
  const authHeader = req.headers['authorization'];
  const token = authHeader && authHeader.split(' ')[1];
  
  if (!token) {
    return res.status(401).json({ success: false, message: 'Токен не предоставлен' });
  }
  
  jwt.verify(token, process.env.JWT_SECRET, (err, user) => {
    if (err) {
      return res.status(401).json({ success: false, message: 'Недействительный токен' });
    }
    req.user = user;
    next();
  });
};
```

## Примеры использования на фронтенде:

### Сохранение токена после входа:
```javascript
localStorage.setItem('authToken', data.token);
localStorage.setItem('user', JSON.stringify(data.user));
```

### Отправка запросов с токеном:
```javascript
const token = localStorage.getItem('authToken');
const response = await fetch('/api/auth/profile', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
```

### Проверка авторизации:
```javascript
const token = localStorage.getItem('authToken');
if (!token) {
  window.location.href = '/login';
}
``` 