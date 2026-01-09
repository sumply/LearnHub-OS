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
    // Проверяем авторизацию перед запросом
    if (!apiClient.isAuthenticated()) {
      console.warn('Пользователь не авторизован. Пользователи не загружены.');
      return [];
    }
    
    usersCache = await apiClient.getUsers();
    cacheTimestamp = now;
    return usersCache;
  } catch (error) {
    console.error('Ошибка загрузки пользователей:', error);
    // Если ошибка авторизации, очищаем кэш
    if (error instanceof Error && error.message.includes('Сессия истекла')) {
      usersCache = null;
      cacheTimestamp = 0;
    }
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
    
    // Пытаемся получить роль из токена
    let role: 'student' | 'parent' | 'teacher' | 'admin' = storedUser?.role || 'student';
    try {
      const token = apiClient.getAuthToken();
      if (token) {
        // JWT токен состоит из трех частей, разделенных точками: header.payload.signature
        const parts = token.split('.');
        if (parts.length >= 2) {
          // Декодируем payload
          const payload = JSON.parse(atob(parts[1]));
          const apiRole = payload.role;
          // Маппинг ролей: 0=student, 1=teacher, 2=admin, 3=root (суперпользователь)
          if (apiRole === 0) role = 'student';
          else if (apiRole === 1) role = 'teacher';
          else if (apiRole === 2) role = 'admin';
          else if (apiRole === 3) role = 'admin'; // root (суперпользователь) маппится в admin
        }
      }
    } catch (err) {
      // Если не удалось распарсить, используем сохраненную роль
      console.warn('Не удалось распарсить роль из токена при загрузке пользователя:', err);
    }
    
    // Обновляем информацию о пользователе
    const user: User = {
      id: userInfo.id.toString(),
      email: storedUser?.email || '',
      name: userInfo.first_name,
      surname: userInfo.last_name,
      role, // Используем роль из токена
    };
    
    setCurrentUser(user);
    return user;
  } catch (error) {
    console.error('Ошибка загрузки текущего пользователя:', error);
    return getStoredUser();
  }
}

// Получить всех пользователей с полной информацией и ролью
// Поскольку API не возвращает роль напрямую, мы получаем полную информацию
// и используем специальный подход для определения роли
export interface UserWithRole extends apiClient.UserFull {
  role: apiClient.UserRole;
}

let usersWithRoleCache: UserWithRole[] | null = null;
let usersWithRoleCacheTimestamp: number = 0;

export async function getAllUsersWithRoles(): Promise<UserWithRole[]> {
  const now = Date.now();
  if (usersWithRoleCache && (now - usersWithRoleCacheTimestamp) < CACHE_TTL) {
    return usersWithRoleCache;
  }

  try {
    if (!apiClient.isAuthenticated()) {
      console.warn('Пользователь не авторизован. Пользователи не загружены.');
      return [];
    }
    
    // Получаем список всех пользователей
    const usersShort = await apiClient.getUsers();
    
    // Для каждого пользователя получаем полную информацию (включая роль из API)
    const usersWithRoles: (UserWithRole | null)[] = await Promise.all(
      usersShort.map(async (userShort) => {
        try {
          const userFull = await apiClient.getUserById(userShort.id);
          
          // Роль теперь приходит из API в userFull.role
          // Проверяем, что роль есть в ответе
          if (userFull.role === undefined || userFull.role === null) {
            console.warn(`Пользователь ${userFull.id} не имеет роли в ответе API`);
            // Используем роль по умолчанию только если её нет в ответе
            return {
              ...userFull,
              role: apiClient.UserRole.STUDENT,
            };
          }
          
          if (import.meta.env.DEV) {
            console.log(`[getAllUsersWithRoles] Пользователь ${userFull.id} (${userFull.last_name} ${userFull.first_name}): роль из API = ${userFull.role}`);
          }
          
          return {
            ...userFull,
            role: userFull.role, // Роль уже есть в ответе API
          };
        } catch (error) {
          console.error(`Ошибка загрузки пользователя ${userShort.id}:`, error);
          return null;
        }
      })
    );
    
    // Фильтруем null значения
    const validUsers = usersWithRoles.filter((u): u is UserWithRole => u !== null);
    
    usersWithRoleCache = validUsers;
    usersWithRoleCacheTimestamp = now;
    return validUsers;
  } catch (error) {
    console.error('Ошибка загрузки пользователей с ролями:', error);
    return usersWithRoleCache || [];
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
  usersWithRoleCache = null;
  cacheTimestamp = 0;
  usersWithRoleCacheTimestamp = 0;
}
