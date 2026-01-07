// Утилиты для работы с API
// Импортируем новый API клиент
import * as apiClient from './apiClient';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000';

export interface User {
  id: string;
  email: string;
  name: string;
  surname: string;
  role: 'student' | 'parent' | 'teacher' | 'admin';
  group?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
  surname: string;
  role: 'student' | 'parent' | 'teacher' | 'admin';
  group?: string;
}

export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  token?: string;
  refresh_token?: string;
  user?: User;
}

export interface MaterialAttachment {
  type: string;
  url?: string;
  testId?: string;
}

export interface Material {
  id: string;
  title: string;
  description: string;
  attachments: MaterialAttachment[];
  category: string;
  visible: boolean;
}

export interface Task {
  id: string;
  title: string;
  description: string;
  attachments: MaterialAttachment[];
  category: string;
  visible: boolean;
}

export interface Flashcard {
  id: string;
  question: string;
  answer: string;
  category: string;
}

export interface Quiz {
  id: string;
  title: string;
  questions: Array<{ question: string; options: string[]; correct: number }>;
  category: string;
}

// Реэкспортируем функции работы с токенами из нового клиента
export const getAuthToken = apiClient.getAuthToken;
export const setAuthToken = apiClient.setAuthToken;
export const removeAuthToken = apiClient.removeAuthToken;

export const getRefreshToken = apiClient.getRefreshToken;
export const setRefreshToken = apiClient.setRefreshToken;
export const removeRefreshToken = apiClient.removeRefreshToken;

export const getCurrentUser = (): User | null => {
  const userStr = localStorage.getItem('user');
  return userStr ? JSON.parse(userStr) : null;
};
export const setCurrentUser = (user: User): void => localStorage.setItem('user', JSON.stringify(user));
export const removeCurrentUser = (): void => localStorage.removeItem('user');

export const isAuthenticated = (): boolean => !!getAuthToken();

// Универсальная функция для запросов с автоматическим обновлением токена
// async function fetchWithAuth(input: RequestInfo, init?: RequestInit, retry = true): Promise<Response> {
//   let headers: HeadersInit = init?.headers || {};
//   const token = getAuthToken();
//   if (token) {
//     headers = { ...headers, Authorization: `Bearer ${token}` };
//   }
//   try {
//     const response = await fetch(input, { ...init, headers });
//     if (response.status === 401 && retry && getRefreshToken()) {
//       // Пробуем обновить токен
//       const refreshed = await refreshToken();
//       if (refreshed) {
//         // Повторяем исходный запрос с новым токеном
//         return fetchWithAuth(input, init, false);
//       } else {
//         handleAuthError(new Error('Сессия истекла'));
//       }
//     }
//     return response;
//   } catch (err) {
//     throw err;
//   }
// }

// Регистрация пользователя
export const register = async (userData: RegisterRequest): Promise<ApiResponse> => {
  try {
    // Маппинг ролей из старого формата в новый
    const roleMap: Record<string, apiClient.UserRole> = {
      'student': apiClient.UserRole.STUDENT,
      'teacher': apiClient.UserRole.TEACHER,
      'admin': apiClient.UserRole.ADMIN,
      'parent': apiClient.UserRole.STUDENT, // parent маппится в student
    };
    
    const createUserData: apiClient.UserCreateRequest = {
      first_name: userData.name,
      last_name: userData.surname,
      email: userData.email,
      role: roleMap[userData.role] || apiClient.UserRole.STUDENT,
    };
    
    await apiClient.createUser(createUserData);
    
    return { 
      success: true, 
      message: 'Пользователь успешно создан. Логин и пароль будут сгенерированы сервером.' 
    };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : 'Ошибка регистрации';
    return { success: false, message: errorMessage };
  }
};

// Вход пользователя
export const login = async (credentials: LoginRequest): Promise<ApiResponse> => {
  try {
    // Используем новый API клиент
    // Преобразуем email в login (API использует login)
    const loginData = {
      login: credentials.email,
      password: credentials.password,
    };
    
    const response = await apiClient.login(loginData);
    
    // Получаем информацию о текущем пользователе
    const userInfo = await apiClient.getCurrentUser();
    
    // Парсим роль из токена (токен в base64 содержит {id, role})
    let role: 'student' | 'parent' | 'teacher' | 'admin' = 'student';
    try {
      const token = response.access_token;
      // Токен - это base64 закодированный JSON объект {id, role}
      // Пробуем декодировать
      let payload;
      try {
        payload = JSON.parse(atob(token));
      } catch {
        // Если не получилось, возможно токен уже в другом формате
        // Пробуем получить роль из информации о пользователе
        payload = { role: 1 }; // По умолчанию student
      }
      
      const apiRole = payload.role;
      // Маппинг ролей: 1=student, 2=teacher, 3=admin, 4=root
      if (apiRole === 1) role = 'student';
      else if (apiRole === 2) role = 'teacher';
      else if (apiRole === 3) role = 'admin';
      else if (apiRole === 4) role = 'admin'; // root маппится в admin
    } catch (err) {
      console.warn('Не удалось распарсить роль из токена:', err);
      // Если не удалось распарсить, используем значение по умолчанию
    }
    
    // Преобразуем в старый формат для совместимости
    const user: User = {
      id: userInfo.id.toString(),
      email: credentials.email,
      name: userInfo.first_name,
      surname: userInfo.last_name,
      role,
    };
    
    setCurrentUser(user);
    
    return { 
      success: true, 
      message: 'Вход выполнен успешно',
      token: response.access_token,
      refresh_token: response.refresh_token,
      user,
    };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : 'Ошибка авторизации';
    return { success: false, message: errorMessage };
  }
};

// Обновление access-токена
export const refreshToken = async (): Promise<boolean> => {
  // const refresh_token = getRefreshToken();
  // if (!refresh_token) return false;
  // try {
  //   const response = await fetch(`${API_BASE_URL}/refresh_access_token`, {
  //     method: 'POST',
  //     headers: { 'Content-Type': 'application/json' },
  //     body: JSON.stringify({ refresh_token }),
  //   });
  //   const data = await response.json();
  //   if (data.success && data.token) {
  //     setAuthToken(data.token);
  //     return true;
  //   } else {
  //     removeAuthToken();
  //     removeRefreshToken();
  //     removeCurrentUser();
  //     return false;
  //   }
  // } catch {
  //   removeAuthToken();
  //   removeRefreshToken();
  //   removeCurrentUser();
  //   return false;
  // }
    return false;
};

export const logout = async (): Promise<ApiResponse> => {
  apiClient.logout();
  removeCurrentUser();
  return {
    success: true,
    message: 'Выход выполнен успешно',
  };
};

export const handleAuthError = (error: Error): void => {
  console.error('Ошибка авторизации:', error);
  apiClient.logout();
  removeCurrentUser();
  window.location.href = '/login';
};

export const checkToken = async (): Promise<boolean> => {
  // const token = getAuthToken();
  // if (!token) return false;
  // try {
  //   const response = await fetch('http://188.225.24.208:8000/profile', {
  //     method: 'GET',
  //     headers: { 'Authorization': `Bearer ${token}` },
  //   });
  //   if (response.status === 200) return true;
  //   return false;
  // } catch {
  //   return false;
  // }
  return true; // Заглушка: всегда валидный токен
};

export const getMaterials = async (): Promise<Material[]> => {
  // const token = getAuthToken();
  // if (!token) throw new Error('Нет токена');
  // const response = await fetch('http://188.225.24.208:8000/materials', {
  //   method: 'GET',
  //   headers: { 'Authorization': `Bearer ${token}` },
  // });
  // if (!response.ok) throw new Error('Ошибка загрузки материалов');
  // const data = await response.json();
  // if (!data.success) throw new Error('Ошибка: ' + (data.message || 'Не удалось получить материалы'));
  // return data.materials;
  return [];
};

export const getTasks = async (): Promise<Task[]> => {
  // const token = getAuthToken();
  // if (!token) throw new Error('Нет токена');
  // const response = await fetch('http://188.225.24.208:8000/tasks', {
  //   method: 'GET',
  //   headers: { 'Authorization': `Bearer ${token}` },
  // });
  // if (!response.ok) throw new Error('Ошибка загрузки заданий');
  // const data = await response.json();
  // if (!data.success) throw new Error('Ошибка: ' + (data.message || 'Не удалось получить задания'));
  // return data.tasks;
  return [];
};

export const getFlashcards = async (): Promise<Flashcard[]> => {
  // const token = getAuthToken();
  // if (!token) throw new Error('Нет токена');
  // const response = await fetch('http://188.225.24.208:8000/flashcards', {
  //   method: 'GET',
  //   headers: { 'Authorization': `Bearer ${token}` },
  // });
  // if (!response.ok) throw new Error('Ошибка загрузки карточек');
  // const data = await response.json();
  // if (!data.success) throw new Error('Ошибка: ' + (data.message || 'Не удалось получить карточки'));
  // return data.flashcards;
  return [];
};

export const getQuizzes = async (): Promise<Quiz[]> => {
  // const token = getAuthToken();
  // if (!token) throw new Error('Нет токена');
  // const response = await fetch('http://188.225.24.208:8000/quizzes', {
  //   method: 'GET',
  //   headers: { 'Authorization': `Bearer ${token}` },
  // });
  // if (!response.ok) throw new Error('Ошибка загрузки тестов');
  // const data = await response.json();
  // if (!data.success) throw new Error('Ошибка: ' + (data.message || 'Не удалось получить тесты'));
  // return data.quizzes;
  return [];
};

// Пример использования fetchWithAuth для защищённых запросов:
// export const getProfile = async (): Promise<ApiResponse<User>> => {
//   const response = await fetchWithAuth(`${API_BASE_URL}/profile`, { method: 'GET' });
//   return await response.json();
// }; 