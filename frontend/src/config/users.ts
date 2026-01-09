// Типы и классы пользователей
export type UserRole = 'student' | 'parent' | 'teacher' | 'admin';

export interface BaseUser {
  id: string;
  email: string;
  password: string;
  name: string;
  surname: string;
  role: UserRole;
}

export interface Admin extends BaseUser {
  role: 'admin';
}
export interface Teacher extends BaseUser {
  role: 'teacher';
  subjects: string[];
}
export interface Parent extends BaseUser {
  role: 'parent';
  childrenIds: string[];
}
export interface Student extends BaseUser {
  role: 'student';
  group: string;
  activityHistory?: Array<{
    activityId: string;
    type: string;
    title: string;
    date: string;
    score?: number;
    maxScore?: number;
    status: 'passed' | 'failed' | 'in_progress';
  }>;
}

export type User = Admin | Teacher | Parent | Student;

// Временные пользователи удалены - теперь используем API
// Данные теперь загружаются через userService.getAllUsers()
export const users: User[] = [];

// Получить всех пользователей (из localStorage + дефолтные)
export function getAllUsers(): User[] {
  const stored = localStorage.getItem('users');
  let arr: User[] = [];
  if (stored) {
    try {
      arr = JSON.parse(stored);
    } catch {}
  }
  // Склеиваем дефолтных и новых, без дубликатов по email
  const emails = new Set(arr.map(u => u.email));
  return [...arr, ...users.filter(u => !emails.has(u.email))];
}

// Добавить пользователя в localStorage
export function addUser(newUser: User): void {
  const all = getAllUsers();
  all.push(newUser);
  localStorage.setItem('users', JSON.stringify(all));
}

// Проверить, занят ли email
export function isEmailTaken(email: string): boolean {
  return !!getAllUsers().find(u => u.email === email);
}

// Функция поиска пользователя по email
export function findUserByEmail(email: string): User | undefined {
  return getAllUsers().find(u => u.email === email);
}

// Функция проверки пароля
export function checkUserPassword(user: User, password: string): boolean {
  return user.password === password;
}

// Функция проверки роли
export function isUserRole(user: User, role: UserRole): boolean {
  return user.role === role;
} 