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

// Тестовые пользователи
export const users: User[] = [
  {
    id: '1',
    email: 'admin@demo.ru',
    password: 'admin123',
    name: 'Админ',
    surname: 'Системный',
    role: 'admin',
  },
  {
    id: '2',
    email: 'teacher@demo.ru',
    password: 'teacher123',
    name: 'Иван',
    surname: 'Петров',
    role: 'teacher',
    subjects: ['math', 'phys'],
  },
  {
    id: '3',
    email: 'parent@demo.ru',
    password: 'parent123',
    name: 'Мария',
    surname: 'Иванова',
    role: 'parent',
    childrenIds: ['4'],
  },
  {
    id: '4',
    email: 'student@demo.ru',
    password: 'student123',
    name: 'Петя',
    surname: 'Иванов',
    role: 'student',
    group: '11А',
    activityHistory: [
      {
        activityId: 'quiz1',
        type: 'quiz',
        title: 'Тест по математике',
        date: '2024-01-15T10:00:00Z',
        score: 8,
        maxScore: 10,
        status: 'passed'
      },
      {
        activityId: 'quiz2',
        type: 'quiz',
        title: 'Тест по физике',
        date: '2024-01-16T14:30:00Z',
        score: 6,
        maxScore: 10,
        status: 'passed'
      }
    ]
  },
  // Новый ученик 1
  {
    id: '5',
    email: 'student2@demo.ru',
    password: 'student234',
    name: 'Саша',
    surname: 'Кузнецов',
    role: 'student',
    group: '10Б',
    activityHistory: [
      {
        activityId: 'quiz1',
        type: 'quiz',
        title: 'Тест по математике',
        date: '2024-01-15T11:00:00Z',
        score: 9,
        maxScore: 10,
        status: 'passed'
      },
      {
        activityId: 'quiz4',
        type: 'quiz',
        title: 'Продвинутый тест по математике',
        date: '2024-01-18T13:00:00Z',
        score: 8,
        maxScore: 10,
        status: 'passed'
      }
    ]
  },
  // Новый ученик 2
  {
    id: '6',
    email: 'student3@demo.ru',
    password: 'student345',
    name: 'Алина',
    surname: 'Смирнова',
    role: 'student',
    group: '9В',
    activityHistory: [
      {
        activityId: 'quiz3',
        type: 'quiz',
        title: 'Тест по английскому языку',
        date: '2024-01-17T09:00:00Z',
        score: 7,
        maxScore: 10,
        status: 'passed'
      },
      {
        activityId: 'quiz5',
        type: 'quiz',
        title: 'Тест по английскому языку (продвинутый)',
        date: '2024-01-19T10:00:00Z',
        status: 'in_progress'
      }
    ]
  },
  // Новый ученик 3
  {
    id: '7',
    email: 'student4@demo.ru',
    password: 'student456',
    name: 'Дмитрий',
    surname: 'Орлов',
    role: 'student',
    group: '11А',
    activityHistory: [
      {
        activityId: 'quiz1',
        type: 'quiz',
        title: 'Тест по математике',
        date: '2024-01-15T12:00:00Z',
        score: 5,
        maxScore: 10,
        status: 'failed'
      },
      {
        activityId: 'quiz2',
        type: 'quiz',
        title: 'Тест по физике',
        date: '2024-01-16T15:00:00Z',
        status: 'in_progress'
      }
    ]
  },
  // Новый ученик 4
  {
    id: '8',
    email: 'student5@demo.ru',
    password: 'student567',
    name: 'Екатерина',
    surname: 'Васильева',
    role: 'student',
    group: '10Б',
    activityHistory: [
      {
        activityId: 'quiz3',
        type: 'quiz',
        title: 'Тест по английскому языку',
        date: '2024-01-17T10:00:00Z',
        status: 'in_progress'
      }
    ]
  },
  // Новый учитель 1
  {
    id: '9',
    email: 'teacher2@demo.ru',
    password: 'teacher234',
    name: 'Ольга',
    surname: 'Сидорова',
    role: 'teacher',
    subjects: ['eng'],
  },
  // Новый учитель 2
  {
    id: '10',
    email: 'teacher3@demo.ru',
    password: 'teacher345',
    name: 'Владимир',
    surname: 'Ковалёв',
    role: 'teacher',
    subjects: ['math', 'eng'],
  },
  // Новый учитель 3
  {
    id: '11',
    email: 'teacher4@demo.ru',
    password: 'teacher456',
    name: 'Марина',
    surname: 'Громова',
    role: 'teacher',
    subjects: ['phys'],
  },
];

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