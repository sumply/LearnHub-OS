// Сервис для работы с квизами через API
import * as apiClient from '../utils/apiClient';

// Кэш для квизов
let quizzesCache: apiClient.QuizShortResponse[] | null = null;
let cacheTimestamp: number = 0;
const CACHE_TTL = 60000; // 1 минута

// Получить все квизы
export async function getQuizzes(): Promise<apiClient.QuizShortResponse[]> {
  const now = Date.now();
  if (quizzesCache && (now - cacheTimestamp) < CACHE_TTL) {
    return quizzesCache;
  }

  try {
    quizzesCache = await apiClient.getQuizzes();
    cacheTimestamp = now;
    return quizzesCache;
  } catch (error) {
    console.error('Ошибка загрузки квизов:', error);
    return quizzesCache || [];
  }
}

// Создать квиз
export async function createQuiz(quizData: apiClient.QuizCreateRequest): Promise<void> {
  try {
    await apiClient.createQuiz(quizData);
    // Очищаем кэш
    quizzesCache = null;
    cacheTimestamp = 0;
  } catch (error) {
    console.error('Ошибка создания квиза:', error);
    throw error;
  }
}

// Удалить квиз
export async function deleteQuiz(quizId: number): Promise<void> {
  try {
    await apiClient.deleteQuiz(quizId);
    // Очищаем кэш
    quizzesCache = null;
    cacheTimestamp = 0;
  } catch (error) {
    console.error('Ошибка удаления квиза:', error);
    throw error;
  }
}

// Очистить кэш
export function clearCache(): void {
  quizzesCache = null;
  cacheTimestamp = 0;
}

