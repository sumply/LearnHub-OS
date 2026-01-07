// Сервис для работы с предметами и группами через API
import * as apiClient from '../utils/apiClient';

// Кэш для предметов и групп
let subjectsCache: apiClient.SubjectResponse[] | null = null;
let groupsCache: apiClient.GroupResponse[] | null = null;
let subjectsCacheTimestamp: number = 0;
let groupsCacheTimestamp: number = 0;
const CACHE_TTL = 60000; // 1 минута

// Получить все предметы
export async function getSubjects(): Promise<apiClient.SubjectResponse[]> {
  const now = Date.now();
  if (subjectsCache && (now - subjectsCacheTimestamp) < CACHE_TTL) {
    return subjectsCache;
  }

  try {
    // Проверяем авторизацию перед запросом
    if (!apiClient.isAuthenticated()) {
      console.warn('Пользователь не авторизован. Предметы не загружены.');
      return [];
    }
    
    subjectsCache = await apiClient.getSubjects();
    subjectsCacheTimestamp = now;
    return subjectsCache;
  } catch (error) {
    console.error('Ошибка загрузки предметов:', error);
    // Если ошибка авторизации, очищаем кэш
    if (error instanceof Error && error.message.includes('Сессия истекла')) {
      subjectsCache = null;
      subjectsCacheTimestamp = 0;
    }
    return subjectsCache || [];
  }
}

// Получить все группы
export async function getGroups(): Promise<apiClient.GroupResponse[]> {
  const now = Date.now();
  if (groupsCache && (now - groupsCacheTimestamp) < CACHE_TTL) {
    return groupsCache;
  }

  try {
    groupsCache = await apiClient.getGroups();
    groupsCacheTimestamp = now;
    return groupsCache;
  } catch (error) {
    console.error('Ошибка загрузки групп:', error);
    return groupsCache || [];
  }
}

// Создать предмет
export async function createSubject(name: string): Promise<void> {
  try {
    await apiClient.createSubject({ name });
    // Очищаем кэш
    subjectsCache = null;
    subjectsCacheTimestamp = 0;
  } catch (error) {
    console.error('Ошибка создания предмета:', error);
    throw error;
  }
}

// Создать группу
export async function createGroup(name: string, curatorId: number): Promise<void> {
  try {
    await apiClient.createGroup({ name, curator_id: curatorId });
    // Очищаем кэш
    groupsCache = null;
    groupsCacheTimestamp = 0;
  } catch (error) {
    console.error('Ошибка создания группы:', error);
    throw error;
  }
}

// Добавить студентов в группу
export async function addStudentsToGroup(groupId: number, studentIds: number[]): Promise<void> {
  try {
    await apiClient.addStudentsToGroup(groupId, { student_ids: studentIds });
  } catch (error) {
    console.error('Ошибка добавления студентов в группу:', error);
    throw error;
  }
}

// Очистить кэш
export function clearCache(): void {
  subjectsCache = null;
  groupsCache = null;
  subjectsCacheTimestamp = 0;
  groupsCacheTimestamp = 0;
}

