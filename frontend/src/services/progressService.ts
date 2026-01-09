// Сервис для работы с прогрессом прохождения квизов через API
import * as apiClient from '../utils/apiClient';

// Кэш для прогресса
let progressCache: apiClient.QuizProgressResponse[] | null = null;
let cacheTimestamp: number = 0;
const CACHE_TTL = 30000; // 30 секунд (прогресс обновляется чаще)

// Получить весь прогресс
export async function getProgress(): Promise<apiClient.QuizProgressResponse[]> {
  const now = Date.now();
  if (progressCache && (now - cacheTimestamp) < CACHE_TTL) {
    return progressCache;
  }

  try {
    progressCache = await apiClient.getProgress();
    cacheTimestamp = now;
    return progressCache;
  } catch (error) {
    console.error('Ошибка загрузки прогресса:', error);
    return progressCache || [];
  }
}

// Начать прохождение квиза
export async function startProgress(progressId: number): Promise<void> {
  try {
    await apiClient.startProgress(progressId);
    // Очищаем кэш
    progressCache = null;
    cacheTimestamp = 0;
  } catch (error) {
    console.error('Ошибка начала прохождения:', error);
    throw error;
  }
}

// Завершить прохождение квиза
export async function finishProgress(progressId: number): Promise<void> {
  try {
    await apiClient.finishProgress(progressId);
    // Очищаем кэш
    progressCache = null;
    cacheTimestamp = 0;
  } catch (error) {
    console.error('Ошибка завершения прохождения:', error);
    throw error;
  }
}

// Обновить ответ на вопрос
export async function updateAnswer(
  progressId: number,
  answerId: number,
  text: string
): Promise<void> {
  try {
    await apiClient.patchAnswer(progressId, answerId, { text });
    // Очищаем кэш
    progressCache = null;
    cacheTimestamp = 0;
  } catch (error) {
    console.error('Ошибка обновления ответа:', error);
    throw error;
  }
}

// Отметить ответ как правильный
export async function markAnswerCorrect(progressId: number, answerId: number): Promise<void> {
  try {
    await apiClient.markAnswerCorrect(progressId, answerId);
    // Очищаем кэш
    progressCache = null;
    cacheTimestamp = 0;
  } catch (error) {
    console.error('Ошибка отметки ответа как правильного:', error);
    throw error;
  }
}

// Отметить ответ как неправильный
export async function markAnswerIncorrect(progressId: number, answerId: number): Promise<void> {
  try {
    await apiClient.markAnswerIncorrect(progressId, answerId);
    // Очищаем кэш
    progressCache = null;
    cacheTimestamp = 0;
  } catch (error) {
    console.error('Ошибка отметки ответа как неправильного:', error);
    throw error;
  }
}

// Очистить кэш
export function clearCache(): void {
  progressCache = null;
  cacheTimestamp = 0;
}

