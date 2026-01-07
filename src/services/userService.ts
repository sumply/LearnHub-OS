// Сервис для работы с пользователями через API
import * as apiClient from '../utils/apiClient';
import { getCurrentUser as getStoredUser, setCurrentUser } from '../utils/api';
import type { User } from '../utils/api';

// Кэш для пользователей
let usersCache: apiClient.UserShort[] | null = null;
let cacheTimestamp: number = 0;
const CACHE_TTL = 60000; // 1 минута

// Получить всех пользователей
export async function getAllUsers(): Promise<apiClient.UserShort[]> {
  const now = Date.now();
  if (usersCache && (now - cacheTimestamp) < CACHE_TTL) {
    return usersCache;
  }

  try {
    usersCache = await apiClient.getUsers();
    cacheTimestamp = now;
    return usersCache;
  } catch (error) {
    console.error('Ошибка загрузки пользователей:', error);
    // Возвращаем кэш, если есть, иначе пустой массив
    return usersCache || [];
  }
}

// Получить пользователя по ID
export async function getUserById(userId: number): Promise<apiClient.UserFull | null> {
  try {
    return await apiClient.getUserById(userId);
  } catch (error) {
    console.error('Ошибка загрузки пользователя:', error);
    return null;
  }
}

// Получить текущего пользователя
export async function getCurrentUser(): Promise<User | null> {
  try {
    const userInfo = await apiClient.getCurrentUser();
    const storedUser = getStoredUser();
    
    // Обновляем информацию о пользователе
    const user: User = {
      id: userInfo.id.toString(),
      email: storedUser?.email || '',
      name: userInfo.first_name,
      surname: userInfo.last_name,
      role: storedUser?.role || 'student', // Роль берем из токена или хранилища
    };
    
    setCurrentUser(user);
    return user;
  } catch (error) {
    console.error('Ошибка загрузки текущего пользователя:', error);
    return getStoredUser();
  }
}

// Найти пользователя по email (через поиск по всем пользователям)
export async function findUserByEmail(email: string): Promise<User | null> {
  try {
    const users = await getAllUsers();
    // В API нет поиска по email, поэтому нужно получить полную информацию
    // Это не оптимально, но API не предоставляет такой функционал
    const storedUser = getStoredUser();
    if (storedUser && storedUser.email === email) {
      return storedUser;
    }
    return null;
  } catch (error) {
    console.error('Ошибка поиска пользователя:', error);
    return null;
  }
}

// Проверить, занят ли email
export async function isEmailTaken(email: string): Promise<boolean> {
  // API не предоставляет прямую проверку email
  // Можно попробовать создать пользователя и обработать ошибку
  // Или просто вернуть false и позволить серверу проверить
  return false;
}

// Очистить кэш
export function clearUsersCache(): void {
  usersCache = null;
  cacheTimestamp = 0;
}


