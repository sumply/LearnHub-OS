// API клиент для работы с бэкендом
// Базовый URL API
// В режиме разработки используем прокси Vite для обхода CORS
const API_BASE_URL = import.meta.env.VITE_API_URL || (import.meta.env.DEV ? '/api' : 'http://localhost:3000');

// Типы согласно документации API
export enum UserRole {
  STUDENT = 1,
  TEACHER = 2,
  ADMIN = 3,
  ROOT = 4,
}

export enum ProgressStatus {
  NOT_STARTED = 0,
  IN_PROGRESS = 1,
  AWAITING_REVIEW = 2,
  COMPLETED = 3,
}

export enum AnswerStatus {
  NOT_ANSWERED = 0,
  AWAITING_REVIEW = 1,
  INCORRECT = 2,
  CORRECT = 3,
}

// Типы запросов и ответов
export interface LoginRequest {
  login: string;
  password: string;
}

export interface LoginResponse {
  refresh_token: string;
  access_token: string;
}

export interface UserShort {
  id: number;
  short_name: string;
}

export interface UserFull {
  id: number;
  first_name: string;
  last_name: string;
  middle_name?: string;
}

export interface UserCreateRequest {
  first_name: string;
  last_name: string;
  middle_name?: string;
  email: string;
  role: UserRole;
}

export interface GroupResponse {
  id: number;
  name: string;
  curator: UserShort;
}

export interface GroupCreateRequest {
  name: string;
  curator_id: number;
}

export interface GroupAddStudentsRequest {
  student_ids: number[];
}

export interface SubjectResponse {
  id: number;
  name: string;
}

export interface SubjectCreateRequest {
  name: string;
}

export interface QuizOption {
  text: string;
  is_correct: boolean;
}

export interface QuizQuestion {
  text: string;
  options?: QuizOption[];
}

export interface QuizCreateRequest {
  title: string;
  summary: string;
  subject_id: number;
  group_ids?: number[];
  questions: QuizQuestion[];
}

export interface QuizShortResponse {
  id: number;
  title: string;
  summary: string;
  owner: UserShort;
  subject: SubjectResponse;
  group?: GroupResponse[];
}

export interface QuizProgressResponse {
  id: number;
  quiz: {
    id: number;
    title: string;
    summary: string;
    owner: UserShort;
    subject: SubjectResponse;
  };
  user: UserShort;
  status: ProgressStatus;
  score: number;
  start_date?: string;
  completed_date?: string;
}

export interface AnswerPatchRequest {
  text: string;
}

export interface ApiError {
  error: string;
}

// Утилиты для работы с токенами
export const getAuthToken = (): string | null => {
  return localStorage.getItem('accessToken');
};

export const setAuthToken = (token: string): void => {
  localStorage.setItem('accessToken', token);
};

export const removeAuthToken = (): void => {
  localStorage.removeItem('accessToken');
};

export const getRefreshToken = (): string | null => {
  return localStorage.getItem('refreshToken');
};

export const setRefreshToken = (token: string): void => {
  localStorage.setItem('refreshToken', token);
};

export const removeRefreshToken = (): void => {
  localStorage.removeItem('refreshToken');
};

export const isAuthenticated = (): boolean => {
  return !!getAuthToken();
};

// Базовая функция для запросов
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getAuthToken();
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  try {
    // Отладочная информация
    if (import.meta.env.DEV) {
      console.log(`[API] ${options.method || 'GET'} ${endpoint}`, {
        hasToken: !!token,
        url: `${API_BASE_URL}${endpoint}`
      });
    }
    
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });

    // Обработка сетевых ошибок
    if (!response.ok) {
      // Если 401 - неавторизован, возможно токен истек
      if (response.status === 401) {
        // Очищаем токен и перенаправляем на логин
        logout();
        throw new Error('Сессия истекла. Пожалуйста, войдите снова.');
      }

      let errorMessage = `HTTP error! status: ${response.status}`;
      try {
        const errorData: ApiError = await response.json();
        errorMessage = errorData.error || errorMessage;
      } catch {
        // Если не удалось распарсить ошибку, используем стандартное сообщение
        if (response.status === 404) {
          errorMessage = 'Ресурс не найден';
        } else if (response.status === 500) {
          errorMessage = 'Ошибка сервера';
        }
      }
      throw new Error(errorMessage);
    }

    // Если ответ пустой (например, для DELETE запросов)
    if (response.status === 204 || response.headers.get('content-length') === '0') {
      return {} as T;
    }

    return await response.json();
  } catch (error) {
    // Обработка сетевых ошибок (CORS, таймаут, и т.д.)
    if (error instanceof TypeError && error.message === 'Failed to fetch') {
      throw new Error('Не удалось подключиться к серверу. Проверьте, что бэкенд запущен на http://localhost:3000');
    }
    throw error;
  }
}

// API методы

// Авторизация
export const login = async (credentials: LoginRequest): Promise<LoginResponse> => {
  const response = await apiRequest<LoginResponse>('/login', {
    method: 'POST',
    body: JSON.stringify(credentials),
  });
  
  setAuthToken(response.access_token);
  setRefreshToken(response.refresh_token);
  
  return response;
};

// Пользователи
export const createUser = async (userData: UserCreateRequest): Promise<void> => {
  await apiRequest<void>('/users', {
    method: 'POST',
    body: JSON.stringify(userData),
  });
};

export const getUsers = async (): Promise<UserShort[]> => {
  return await apiRequest<UserShort[]>('/users');
};

export const getCurrentUser = async (): Promise<UserFull> => {
  return await apiRequest<UserFull>('/users/me');
};

export const getUserById = async (userId: number): Promise<UserFull> => {
  return await apiRequest<UserFull>(`/users/${userId}`);
};

// Группы
export const createGroup = async (groupData: GroupCreateRequest): Promise<void> => {
  await apiRequest<void>('/groups', {
    method: 'POST',
    body: JSON.stringify(groupData),
  });
};

export const getGroups = async (): Promise<GroupResponse[]> => {
  return await apiRequest<GroupResponse[]>('/groups');
};

export const addStudentsToGroup = async (
  groupId: number,
  studentIds: GroupAddStudentsRequest
): Promise<void> => {
  await apiRequest<void>(`/groups/${groupId}/students`, {
    method: 'POST',
    body: JSON.stringify(studentIds),
  });
};

// Предметы
export const createSubject = async (subjectData: SubjectCreateRequest): Promise<void> => {
  await apiRequest<void>('/subject', {
    method: 'POST',
    body: JSON.stringify(subjectData),
  });
};

export const getSubjects = async (): Promise<SubjectResponse[]> => {
  return await apiRequest<SubjectResponse[]>('/subject');
};

// Квизы
export const createQuiz = async (quizData: QuizCreateRequest): Promise<void> => {
  await apiRequest<void>('/quizzes', {
    method: 'POST',
    body: JSON.stringify(quizData),
  });
};

export const getQuizzes = async (): Promise<QuizShortResponse[]> => {
  return await apiRequest<QuizShortResponse[]>('/quizzes');
};

export const deleteQuiz = async (quizId: number): Promise<void> => {
  await apiRequest<void>(`/quizzes/${quizId}`, {
    method: 'DELETE',
  });
};

// Прогресс
export const getProgress = async (): Promise<QuizProgressResponse[]> => {
  return await apiRequest<QuizProgressResponse[]>('/progress');
};

export const startProgress = async (progressId: number): Promise<void> => {
  await apiRequest<void>(`/progress/${progressId}/start`, {
    method: 'POST',
  });
};

export const finishProgress = async (progressId: number): Promise<void> => {
  await apiRequest<void>(`/progress/${progressId}/finish`, {
    method: 'POST',
  });
};

export const patchAnswer = async (
  progressId: number,
  answerId: number,
  answerData: AnswerPatchRequest
): Promise<void> => {
  await apiRequest<void>(`/progress/${progressId}/answer/${answerId}`, {
    method: 'PATCH',
    body: JSON.stringify(answerData),
  });
};

export const markAnswerCorrect = async (
  progressId: number,
  answerId: number
): Promise<void> => {
  await apiRequest<void>(`/progress/${progressId}/answer/${answerId}/correct`, {
    method: 'POST',
  });
};

export const markAnswerIncorrect = async (
  progressId: number,
  answerId: number
): Promise<void> => {
  await apiRequest<void>(`/progress/${progressId}/answer/${answerId}/incorrect`, {
    method: 'POST',
  });
};

// Выход
export const logout = (): void => {
  removeAuthToken();
  removeRefreshToken();
  localStorage.removeItem('user');
};

