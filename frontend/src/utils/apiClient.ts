// API клиент для работы с бэкендом
// Базовый URL API
// В режиме разработки используем прокси Vite для обхода CORS
const API_BASE_URL = import.meta.env.VITE_API_URL || (import.meta.env.DEV ? '/api' : 'http://85.239.55.179:8000');

// Флаг для работы без бэкенда (локальная разработка) - отключен
const LOCAL_DEV_MODE = false;

// Типы согласно документации API
// Роли теперь строковые: "student", "teacher", "admin"
export type UserRole = 'student' | 'teacher' | 'admin' | 'root';

export type ProgressStatus = 0 | 1 | 2 | 3;
export const ProgressStatus = {
  NOT_STARTED: 0,
  IN_PROGRESS: 1,
  AWAITING_REVIEW: 2,
  COMPLETED: 3,
} as const;

export type AnswerStatus = 0 | 1 | 2 | 3;
export const AnswerStatus = {
  NOT_ANSWERED: 0,
  AWAITING_REVIEW: 1,
  INCORRECT: 2,
  CORRECT: 3,
} as const;

// Типы запросов и ответов
export interface LoginRequest {
  email: string;
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
  id: string | number;
  first_name: string;
  last_name: string;
  middle_name?: string;
  role: UserRole | string;
}

export interface UserCreateRequest {
  first_name: string;
  last_name: string;
  middle_name?: string;
  email: string;
  role: UserRole | string;
  subject_ids?: number[]; // Для учителей - предметы
  group_ids?: number[]; // Для учителей - группы
}

export interface GroupResponse {
  id: string; // UUID
  name: string;
  curator?: UserShort; // Необязательно, так как группа может быть создана без куратора
}

export interface GroupCreateRequest {
  name: string;
  curator_id?: string; // UUID, необязательно
  student_ids?: string[]; // UUID[], необязательно
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

// Новый формат вопроса для бэкенда
export interface QuizQuestionDetailsSingle {
  options: string[];
  correct: string; // Одна строка для single choice
}

export interface QuizQuestionDetailsMultiple {
  options: string[];
  correct: string[]; // Массив строк для multiple choice
}

export interface QuizQuestionDetailsNumeric {
  correct: number; // Число для numeric
}

export type QuizQuestionDetails = QuizQuestionDetailsSingle | QuizQuestionDetailsMultiple | QuizQuestionDetailsNumeric;

export interface QuizQuestion {
  text: string;
  score: number; // Баллы за вопрос
  type: 'single' | 'multiple' | 'numeric';
  details: QuizQuestionDetails;
}

export interface QuizCreateRequest {
  title: string;
  owner_id: string; // UUID владельца
  summary: string;
  subject_id: string; // UUID предмета или "все темы"
  deadline?: string; // ISO дата/время (необязательно)
  group_ids?: string[]; // UUID[] групп или ["общий"]
  max_attempts?: number; // Максимальное количество попыток (необязательно)
  questions: QuizQuestion[];
}

export interface QuizShortResponse {
  id: number;
  title: string;
  summary: string;
  total_score: number;
  owner: {
    id: number;
    short_name: string;
    role: UserRole;
  };
  subject: {
    id: number;
    name: string;
  };
  group: Array<{
    id: number;
    name: string;
    curator: {
      id: number;
      short_name: string;
      role: UserRole;
    };
  }>;
}

export interface QuizFullResponse {
  id: number;
  title: string;
  summary: string;
  total_score: number;
  owner: {
    id: number;
    short_name: string;
    role: UserRole;
  };
  subject: {
    id: number;
    name: string;
  };
  group: Array<{
    id: number;
    name: string;
    curator: {
      id: number;
      short_name: string;
      role: UserRole;
    };
  }>;
  questions: QuizQuestion[];
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

// Типы для попыток прохождения квизов (старый формат, оставляем для совместимости)
export interface QuizAttemptAnswer {
  id: string;
  question_id: string;
  answer: string | null;
  score: number;
  is_correct: boolean;
}

export interface QuizAttempt {
  id: string;
  user: {
    id: string;
    first_name: string;
    last_name: string;
    role: string;
  };
  answers: QuizAttemptAnswer[];
  score: number;
  started_at: string;
  ended_at: string;
}

export interface QuizAttemptResponse {
  attempt: QuizAttempt;
  quiz: {
    id: string;
    title: string;
    summary: string;
    owner: {
      id: string;
      first_name: string;
      last_name: string;
      role: string;
    };
    subject: {
      id: string;
      name: string;
    };
    total_score: number;
    deadline: string;
    max_attempts: number;
    created_at: string;
    content: Array<{
      id: string;
      text: string;
      score: number;
      type: string;
      details: any;
    }>;
  };
}

// Новые типы для статистики квиза через /quizzes/{quiz_id}/users
export interface AttemptItem {
  id: string; // UUID
  score: number;
  started_at: string; // ISO date string
  ended_at: string | null; // ISO date string или null
}

export interface UserLastAttempt {
  id: string; // UUID
  first_name: string;
  last_name: string;
  role: string;
  last_attempt: AttemptItem | null; // Может быть null если попыток нет
}

// Типы для прохождения квиза
export interface QuizQuestionContent {
  id: string; // UUID
  text: string;
  score: number;
  type: 'single' | 'multiple' | 'numeric';
  details: QuizQuestionDetails;
}

export interface StartAttemptResponse {
  attempt: {
    id: string; // UUID
    started_at: string; // ISO date string
  };
  quiz: {
    id: string; // UUID
    title: string;
    summary: string;
    owner: {
      id: string;
      first_name: string;
      last_name: string;
      role: string;
    };
    subject: {
      id: string;
      name: string;
    };
    total_score: number;
    deadline: string; // ISO date string
    max_attempts: number;
    created_at: string; // ISO date string
    content: QuizQuestionContent[];
  };
}

export interface FinishAttemptAnswer {
  id: string; // UUID
  question_id: string; // UUID
  answer: string | string[] | number; // Зависит от типа вопроса
  score: number;
  is_correct: boolean;
}

export interface FinishAttemptResponse {
  attempt: {
    id: string; // UUID
    user: {
      id: string;
      first_name: string;
      last_name: string;
      role: string;
    };
    answers: FinishAttemptAnswer[];
    score: number;
    started_at: string; // ISO date string
    ended_at: string; // ISO date string
  };
  quiz: {
    id: string; // UUID
    title: string;
    summary: string;
    owner: {
      id: string;
      first_name: string;
      last_name: string;
      role: string;
    };
    subject: {
      id: string;
      name: string;
    };
    total_score: number;
    deadline: string;
    max_attempts: number;
    created_at: string;
    content: QuizQuestionContent[];
  };
}

// Утилиты для работы с авторизацией (user_id и role)
export interface AuthData {
  user_id: string;
  role: string;
}

export const getAuthData = (): AuthData | null => {
  const userId = localStorage.getItem('user_id');
  const role = localStorage.getItem('user_role');
  if (userId && role) {
    return { user_id: userId, role };
  }
  return null;
};

export const setAuthData = (userId: string, role: string): void => {
  localStorage.setItem('user_id', userId);
  localStorage.setItem('user_role', role);
};

export const removeAuthData = (): void => {
  localStorage.removeItem('user_id');
  localStorage.removeItem('user_role');
};

export const isAuthenticated = (): boolean => {
  return !!getAuthData();
};

// Обратная совместимость (для старых методов)
export const getAuthToken = (): string | null => {
  const authData = getAuthData();
  return authData ? JSON.stringify(authData) : null;
};

export const setAuthToken = (token: string): void => {
  // Игнорируем, так как теперь используем user_id и role
};

export const removeAuthToken = (): void => {
  removeAuthData();
};

export const getRefreshToken = (): string | null => {
  return null; // Больше не используется
};

export const setRefreshToken = (token: string): void => {
  // Игнорируем, так как больше не используется
};

export const removeRefreshToken = (): void => {
  // Игнорируем, так как больше не используется
};

// Базовая функция для запросов
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  // Локальный режим разработки без бэкенда
  if (LOCAL_DEV_MODE) {
    console.warn(`[API] Локальный режим разработки: запрос ${options.method || 'GET'} ${endpoint} пропущен`);
    // Возвращаем пустые данные в зависимости от типа запроса
    if (endpoint.includes('/auth/login')) {
      // Мок-данные для логина
      return {
        id: 'mock_user_' + Date.now(),
        role: 'student',
      } as T;
    }
    if (endpoint.includes('/users/me')) {
      // Мок-данные для текущего пользователя
      return {
        id: 1,
        first_name: 'Тестовый',
        last_name: 'Пользователь',
        role: 'student',
      } as T;
    }
    if (endpoint.includes('/users') && options.method !== 'POST' && !endpoint.includes('/users/')) {
      return { users: [] } as T;
    }
    if (endpoint.includes('/groups') && !endpoint.includes('/groups/')) {
      return { groups: [] } as T;
    }
    if (endpoint.includes('/subjects')) {
      return { subjects: [] } as T;
    }
    if (endpoint.includes('/quizzes') && !endpoint.includes('/quizzes/')) {
      return { quizzes: [] } as T;
    }
    if (endpoint.includes('/progress')) {
      return { progress: [] } as T;
    }
    if (endpoint.includes('/groups/') && endpoint.includes('/students')) {
      return { students: [] } as T;
    }
    if (endpoint.includes('/quizzes/') && endpoint.includes('/attempts')) {
      return { attempts: [] } as T;
    }
    if (endpoint.includes('/users/') && endpoint.includes('/quizzes')) {
      return { quizzes: [] } as T;
    }
    // Для POST/PATCH/DELETE просто возвращаем пустой объект
    return {} as T;
  }

  if (!API_BASE_URL) {
    throw new Error('API_BASE_URL не настроен. Установите VITE_API_URL или включите proxy в vite.config.ts');
  }

  const authData = getAuthData();
  const method = options.method || 'GET';
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> || {}),
  };

  // Для GET запросов отправляем user_id и role в заголовках
  // Для POST/PATCH/PUT - в теле запроса
  if (authData && method === 'GET') {
    headers['X-User-Id'] = authData.user_id;
    headers['X-User-Role'] = authData.role;
  }

  try {
    // Подготавливаем тело запроса с user_id и role для POST/PATCH/PUT запросов
    let requestBody = options.body;
    
    // Если это POST/PATCH/PUT и есть данные авторизации, добавляем user_id и role в тело
    // НО: не перезаписываем role, если он уже есть в теле запроса (например, при создании пользователя)
    if (authData && (method === 'POST' || method === 'PATCH' || method === 'PUT')) {
      try {
        const bodyData = requestBody ? JSON.parse(requestBody as string) : {};
        bodyData.user_id = authData.user_id;
        // Добавляем role только если его нет в теле запроса (чтобы не перезаписывать роль при создании пользователя)
        if (!bodyData.role) {
          bodyData.role = authData.role;
        }
        requestBody = JSON.stringify(bodyData);
      } catch (e) {
        // Если не удалось распарсить тело, создаем новое
        requestBody = JSON.stringify({
          user_id: authData.user_id,
          role: authData.role,
        });
      }
    }
    
    // Отладочная информация
    if (import.meta.env.DEV) {
      const bodyForLog = requestBody ? (() => {
        try {
          return JSON.parse(requestBody as string);
        } catch {
          return requestBody;
        }
      })() : undefined;
      console.log(`[API] ${method} ${endpoint}`, {
        hasAuth: !!authData,
        userId: authData?.user_id,
        role: authData?.role,
        url: `${API_BASE_URL}${endpoint}`,
        body: bodyForLog
      });
    }
    
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
      body: requestBody,
    });

    // Обработка сетевых ошибок
    if (!response.ok) {
      // Если 401 - неавторизован
      if (response.status === 401) {
        // Очищаем данные авторизации и перенаправляем на логин
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
      const backendUrl = import.meta.env.DEV ? 'http://85.239.55.179:8000 (через proxy /api)' : API_BASE_URL;
      throw new Error(`Не удалось подключиться к серверу. Проверьте, что бэкенд доступен на ${backendUrl}`);
    }
    throw error;
  }
}

// API методы

// Авторизация
export const login = async (credentials: LoginRequest): Promise<{ user_id: string; role: string }> => {
  // Локальный режим разработки - мок-авторизация
  if (LOCAL_DEV_MODE) {
    console.warn('[API] Локальный режим: используется мок-авторизация');
    const mockUserId = 'mock_user_' + Date.now();
    const mockRole = 'student';
    setAuthData(mockUserId, mockRole);
    return { user_id: mockUserId, role: mockRole };
  }

  const response = await apiRequest<{ id: string; role: string }>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(credentials),
  });
  
  // Сохраняем user_id и role из ответа
  const userId = response.id;
  const role = response.role;
  setAuthData(userId, role);
  
  return { user_id: userId, role };
};

// Пользователи
export const createUser = async (userData: UserCreateRequest): Promise<void> => {
  await apiRequest<void>('/users', {
    method: 'POST',
    body: JSON.stringify(userData),
  });
};

export const getUsers = async (): Promise<UserShort[]> => {
  const response = await apiRequest<{ users: UserShort[] }>('/users');
  return response.users || [];
};

export const getCurrentUser = async (): Promise<UserFull> => {
  const authData = getAuthData();
  if (!authData) {
    throw new Error('Пользователь не авторизован');
  }
  const response = await apiRequest<{ user: UserFull }>(`/users/${authData.user_id}`, {
    method: 'GET',
  });
  // Проверяем структуру ответа
  if (response && 'user' in response) {
    return response.user;
  }
  // Если ответ пришел напрямую как UserFull
  return response as unknown as UserFull;
};

export const getUserById = async (userId: number): Promise<UserFull> => {
  const response = await apiRequest<{ user: UserFull }>(`/users/${userId}`);
  // Проверяем структуру ответа
  if (response && 'user' in response) {
    return response.user;
  }
  // Если ответ пришел напрямую как UserFull
  return response as unknown as UserFull;
};

// Группы
export const createGroup = async (groupData: GroupCreateRequest): Promise<void> => {
  await apiRequest<void>('/groups', {
    method: 'POST',
    body: JSON.stringify(groupData),
  });
};

export const getGroups = async (): Promise<GroupResponse[]> => {
  const response = await apiRequest<{ groups: GroupResponse[] }>('/groups');
  return response.groups || [];
};

export const addStudentsToGroup = async (
  groupId: string | number,
  studentIds: GroupAddStudentsRequest
): Promise<void> => {
  const groupIdStr = String(groupId);
  await apiRequest<void>(`/groups/${groupIdStr}/students`, {
    method: 'POST',
    body: JSON.stringify(studentIds),
  });
};

// Получить студентов группы
export const getGroupStudents = async (groupId: string | number): Promise<UserShort[]> => {
  const groupIdStr = String(groupId);
  const response = await apiRequest<{ students: UserShort[] }>(`/groups/${groupIdStr}/students`);
  return response.students || [];
};

// Предметы
export const createSubject = async (subjectData: SubjectCreateRequest): Promise<void> => {
  await apiRequest<void>('/subjects', {
    method: 'POST',
    body: JSON.stringify(subjectData),
  });
};

export const getSubjects = async (): Promise<SubjectResponse[]> => {
  const response = await apiRequest<{ subjects: SubjectResponse[] }>('/subjects');
  return response.subjects || [];
};

// Квизы
export const createQuiz = async (quizData: QuizCreateRequest): Promise<void> => {
  await apiRequest<void>('/quizzes', {
    method: 'POST',
    body: JSON.stringify(quizData),
  });
};

export const getQuizzes = async (): Promise<QuizShortResponse[]> => {
  const response = await apiRequest<{ quizzes: QuizShortResponse[] }>('/quizzes');
  return response.quizzes || [];
};

export const getQuizById = async (quizId: number): Promise<QuizFullResponse> => {
  const response = await apiRequest<{ quiz: QuizFullResponse }>(`/quizzes/${quizId}`);
  // Проверяем структуру ответа
  if (response && 'quiz' in response) {
    return response.quiz;
  }
  // Если ответ пришел напрямую как QuizFullResponse
  return response as unknown as QuizFullResponse;
};

export const deleteQuiz = async (quizId: number): Promise<void> => {
  await apiRequest<void>(`/quizzes/${quizId}`, {
    method: 'DELETE',
  });
};

// Получение попыток прохождения квиза (старый эндпоинт)
export const getQuizAttempts = async (quizId: string | number): Promise<QuizAttemptResponse[]> => {
  const response = await apiRequest<{ attempts: QuizAttemptResponse[] }>(`/quizzes/${quizId}/attempts`);
  return response.attempts || [];
};

// Получение студентов с их последними попытками для квиза
export const getQuizUsers = async (quizId: string | number): Promise<UserLastAttempt[]> => {
  const response = await apiRequest<UserLastAttempt[]>(`/quizzes/${quizId}/users`);
  return response || [];
};

// Начать попытку прохождения квиза
export const startQuizAttempt = async (quizId: string | number): Promise<StartAttemptResponse> => {
  const response = await apiRequest<StartAttemptResponse>(`/quizzes/${quizId}/attempt`, {
    method: 'POST',
  });
  return response;
};

// Завершить попытку прохождения квиза
// Ответы должны быть отправлены отдельно через PATCH /attempts/{attempt_id}/answers/{answer_id}
// или включены в тело запроса finish (зависит от реализации бэкенда)
export const finishQuizAttempt = async (attemptId: string, answers?: Array<{ question_id: string; answer: string | string[] | number }>): Promise<FinishAttemptResponse> => {
  const body = answers ? { answers } : undefined;
  const response = await apiRequest<FinishAttemptResponse>(`/attempts/${attemptId}/finish`, {
    method: 'POST',
    body: body ? JSON.stringify(body) : undefined,
  });
  return response;
};

// Прогресс
export const getProgress = async (): Promise<QuizProgressResponse[]> => {
  const response = await apiRequest<{ progress: QuizProgressResponse[] }>('/progress');
  return response.progress || [];
};

// Получить квизы пользователя
export const getUserQuizzes = async (userId?: string): Promise<QuizShortResponse[]> => {
  const authData = getAuthData();
  let targetUserId = userId;
  
  // Если userId не передан, берем из authData
  if (!targetUserId && authData) {
    targetUserId = authData.user_id;
  }
  
  // Проверяем, что user_id валидный (не 'true', не пустая строка)
  const userIdStr = String(targetUserId || '');
  if (!targetUserId || userIdStr === 'true' || userIdStr === 'false' || userIdStr.trim() === '') {
    throw new Error('Пользователь не авторизован или user_id неверный');
  }
  
  // Убеждаемся, что user_id - это строка
  const userIdString = String(targetUserId);
  const response = await apiRequest<{ quizzes: QuizShortResponse[] }>(`/users/${userIdString}/quizzes`);
  return response.quizzes || [];
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
  removeAuthData();
  localStorage.removeItem('user');
};

