const express = require('express');
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const cors = require('cors');
const { body, validationResult } = require('express-validator');

const app = express();
const PORT = process.env.PORT || 3000;
const JWT_SECRET = process.env.JWT_SECRET || 'your-secret-key';

// Middleware
app.use(cors());
app.use(express.json());

// Имитация базы данных (в реальном проекте используйте MongoDB, PostgreSQL и т.д.)
let users = [];

// Middleware для проверки JWT токена
const authenticateToken = (req, res, next) => {
  const authHeader = req.headers['authorization'];
  const token = authHeader && authHeader.split(' ')[1];
  
  if (!token) {
    return res.status(401).json({ 
      success: false, 
      message: 'Токен не предоставлен' 
    });
  }
  
  jwt.verify(token, JWT_SECRET, (err, user) => {
    if (err) {
      return res.status(401).json({ 
        success: false, 
        message: 'Недействительный токен' 
      });
    }
    req.user = user;
    next();
  });
};

// Валидация данных регистрации
const validateRegister = [
  body('email').isEmail().withMessage('Некорректный email'),
  body('password').isLength({ min: 6 }).withMessage('Пароль должен содержать минимум 6 символов'),
  body('name').notEmpty().withMessage('Имя обязательно'),
  body('surname').notEmpty().withMessage('Фамилия обязательна'),
  body('role').isIn(['student', 'parent', 'teacher']).withMessage('Некорректная роль'),
];

// Валидация данных входа
const validateLogin = [
  body('email').isEmail().withMessage('Некорректный email'),
  body('password').notEmpty().withMessage('Пароль обязателен'),
];

// Генерация уникального ID
const generateId = () => {
  return Date.now().toString(36) + Math.random().toString(36).substr(2);
};

// 1. Регистрация пользователя
app.post('/api/auth/register', validateRegister, async (req, res) => {
  try {
    // Проверка ошибок валидации
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(422).json({
        success: false,
        message: errors.array()[0].msg
      });
    }

    const { email, password, name, surname, role, group } = req.body;

    // Проверка, что email не занят
    const existingUser = users.find(user => user.email === email);
    if (existingUser) {
      return res.status(400).json({
        success: false,
        message: 'Email уже используется'
      });
    }

    // Хеширование пароля
    const hashedPassword = await bcrypt.hash(password, 10);

    // Создание пользователя
    const newUser = {
      id: generateId(),
      email,
      password: hashedPassword,
      name,
      surname,
      role,
      group: role !== 'teacher' ? group : undefined,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };

    users.push(newUser);

    // Удаляем пароль из ответа
    const { password: _, ...userWithoutPassword } = newUser;

    res.status(200).json({
      success: true,
      message: 'Пользователь успешно зарегистрирован',
      user: userWithoutPassword
    });

  } catch (error) {
    console.error('Ошибка регистрации:', error);
    res.status(500).json({
      success: false,
      message: 'Внутренняя ошибка сервера'
    });
  }
});

// 2. Вход пользователя
app.post('/api/auth/login', validateLogin, async (req, res) => {
  try {
    // Проверка ошибок валидации
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(422).json({
        success: false,
        message: errors.array()[0].msg
      });
    }

    const { email, password } = req.body;

    // Поиск пользователя
    const user = users.find(u => u.email === email);
    if (!user) {
      return res.status(401).json({
        success: false,
        message: 'Неверный email или пароль'
      });
    }

    // Проверка пароля
    const isValidPassword = await bcrypt.compare(password, user.password);
    if (!isValidPassword) {
      return res.status(401).json({
        success: false,
        message: 'Неверный email или пароль'
      });
    }

    // Создание JWT токена
    const token = jwt.sign(
      { 
        userId: user.id, 
        email: user.email, 
        role: user.role 
      },
      JWT_SECRET,
      { expiresIn: '24h' }
    );

    // Удаляем пароль из ответа
    const { password: _, ...userWithoutPassword } = user;

    res.status(200).json({
      success: true,
      message: 'Вход выполнен успешно',
      token,
      user: userWithoutPassword
    });

  } catch (error) {
    console.error('Ошибка входа:', error);
    res.status(500).json({
      success: false,
      message: 'Внутренняя ошибка сервера'
    });
  }
});

// 3. Выход пользователя
app.post('/api/auth/logout', authenticateToken, (req, res) => {
  // В JWT аутентификации выход обычно происходит на клиенте
  // путем удаления токена. Здесь можно добавить логику
  // для добавления токена в черный список, если нужно.
  
  res.status(200).json({
    success: true,
    message: 'Выход выполнен успешно'
  });
});

// 4. Проверка токена
app.get('/api/auth/verify', authenticateToken, (req, res) => {
  const user = users.find(u => u.id === req.user.userId);
  
  if (!user) {
    return res.status(401).json({
      success: false,
      message: 'Пользователь не найден'
    });
  }

  const { password: _, ...userWithoutPassword } = user;

  res.status(200).json({
    success: true,
    message: 'Токен действителен',
    user: userWithoutPassword
  });
});

// 5. Получение профиля пользователя
app.get('/api/auth/profile', authenticateToken, (req, res) => {
  const user = users.find(u => u.id === req.user.userId);
  
  if (!user) {
    return res.status(404).json({
      success: false,
      message: 'Пользователь не найден'
    });
  }

  const { password: _, ...userWithoutPassword } = user;

  res.status(200).json({
    success: true,
    user: userWithoutPassword
  });
});

// 6. Обновление профиля пользователя
app.put('/api/auth/profile', authenticateToken, [
  body('name').optional().notEmpty().withMessage('Имя не может быть пустым'),
  body('surname').optional().notEmpty().withMessage('Фамилия не может быть пустой'),
  body('group').optional().notEmpty().withMessage('Группа не может быть пустой'),
], (req, res) => {
  try {
    // Проверка ошибок валидации
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(422).json({
        success: false,
        message: errors.array()[0].msg
      });
    }

    const userIndex = users.findIndex(u => u.id === req.user.userId);
    
    if (userIndex === -1) {
      return res.status(404).json({
        success: false,
        message: 'Пользователь не найден'
      });
    }

    // Обновление данных пользователя
    const { name, surname, group } = req.body;
    if (name) users[userIndex].name = name;
    if (surname) users[userIndex].surname = surname;
    if (group && users[userIndex].role !== 'teacher') {
      users[userIndex].group = group;
    }
    users[userIndex].updatedAt = new Date().toISOString();

    const { password: _, ...userWithoutPassword } = users[userIndex];

    res.status(200).json({
      success: true,
      message: 'Профиль обновлен успешно',
      user: userWithoutPassword
    });

  } catch (error) {
    console.error('Ошибка обновления профиля:', error);
    res.status(500).json({
      success: false,
      message: 'Внутренняя ошибка сервера'
    });
  }
});

// Тестовый эндпоинт для проверки работы сервера
app.get('/api/health', (req, res) => {
  res.status(200).json({
    success: true,
    message: 'Сервер работает',
    timestamp: new Date().toISOString()
  });
});

// Обработка ошибок
app.use((err, req, res, next) => {
  console.error('Ошибка сервера:', err);
  res.status(500).json({
    success: false,
    message: 'Внутренняя ошибка сервера'
  });
});

// Обработка несуществующих маршрутов
app.use('*', (req, res) => {
  res.status(404).json({
    success: false,
    message: 'Маршрут не найден'
  });
});

app.listen(PORT, () => {
  console.log(`Сервер запущен на порту ${PORT}`);
  console.log(`API доступен по адресу: http://localhost:${PORT}/api`);
});

module.exports = app; 