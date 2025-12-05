// Утилиты для работы с API

const API_BASE_URL = 'http://188.225.24.208:8000';

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

export const getAuthToken = (): string | null => localStorage.getItem('authToken');
export const setAuthToken = (token: string): void => localStorage.setItem('authToken', token);
export const removeAuthToken = (): void => localStorage.removeItem('authToken');

export const getRefreshToken = (): string | null => localStorage.getItem('refreshToken');
export const setRefreshToken = (token: string): void => localStorage.setItem('refreshToken', token);
export const removeRefreshToken = (): void => localStorage.removeItem('refreshToken');

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
  // const response = await fetch(`${API_BASE_URL}/registration`, {
  //   method: 'POST',
  //   headers: { 'Content-Type': 'application/json' },
  //   body: JSON.stringify(userData),
  // });
  // return await response.json();
  return { success: true, message: 'Заглушка регистрации (сервер отключён)' };
};

// Вход пользователя
export const login = async (credentials: LoginRequest): Promise<ApiResponse> => {
  // const response = await fetch(`${API_BASE_URL}/authorization`, {
  //   method: 'POST',
  //   headers: { 'Content-Type': 'application/json' },
  //   body: JSON.stringify(credentials),
  // });
  // const data = await response.json();
  // if (data.access && data.refresh) {
  //   setAuthToken(data.access);
  //   setRefreshToken(data.refresh);
  //   return { success: true, message: 'Вход выполнен успешно' };
  // }
  // return { success: false, message: 'Ошибка авторизации' };
  return { success: true, message: 'Заглушка входа (сервер отключён)' };
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
  removeAuthToken();
  removeRefreshToken();
  removeCurrentUser();
  return {
    success: true,
    message: 'Выход выполнен успешно',
  };
};

export const handleAuthError = (error: Error): void => {
  console.error('Ошибка авторизации:', error);
  removeAuthToken();
  removeRefreshToken();
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