// Типы для тестов и карточек
import type { MaterialAttachment } from '../utils/api';

export interface Quiz {
  id: string;
  title: string;
  requiresTeacherCheck: boolean; // нужен ли разбор учителем
  questions: Array<{
    id: string;
    type: 'choice' | 'open';
    question: string;
    options?: string[]; // только для choice
    correct?: number;   // только для choice
  }>;
  category: string;
  teacher: string;
  attachments?: MaterialAttachment[];
}

export interface Flashcard {
  id: string;
  question: string;
  answer: string;
  category: string;
  teacher: string;
  attachments?: MaterialAttachment[];
}

// Тестовые тесты
export const quizzes: Quiz[] = [
  {
    id: 'quiz1',
    title: 'Тест по математике',
    requiresTeacherCheck: false,
    category: 'math',
    teacher: 'Иван Петров',
    questions: [
      {
        id: 'q1',
        type: 'choice',
        question: '2 + 2 = ?',
        options: ['3', '4', '5', '6'],
        correct: 1,
      },
      {
        id: 'q2',
        type: 'choice',
        question: '5 * 3 = ?',
        options: ['8', '15', '10', '20'],
        correct: 1,
      },
    ],
  },
  {
    id: 'quiz2',
    title: 'Тест по физике',
    requiresTeacherCheck: false,
    category: 'phys',
    teacher: 'Марина Громова',
    questions: [
      {
        id: 'q1',
        type: 'choice',
        question: 'Сколько секунд в одной минуте?',
        options: ['60', '100', '30', '120'],
        correct: 0,
      },
      {
        id: 'q2',
        type: 'choice',
        question: 'Что измеряется в Ньютонах?',
        options: ['Масса', 'Сила', 'Время', 'Температура'],
        correct: 1,
      },
    ],
  },
  {
    id: 'quiz3',
    title: 'Тест по английскому языку',
    requiresTeacherCheck: false,
    category: 'eng',
    teacher: 'Ольга Сидорова',
    questions: [
      {
        id: 'q1',
        type: 'choice',
        question: 'What is the capital of Great Britain?',
        options: ['London', 'Paris', 'Berlin', 'Madrid'],
        correct: 0,
      },
      {
        id: 'q2',
        type: 'choice',
        question: 'Choose the correct form: "He ___ to school every day."',
        options: ['go', 'goes', 'going', 'gone'],
        correct: 1,
      },
    ],
  },
  {
    id: 'quiz4',
    title: 'Продвинутый тест по математике',
    requiresTeacherCheck: false,
    category: 'math',
    teacher: 'Владимир Ковалёв',
    questions: [
      {
        id: 'q1',
        type: 'choice',
        question: 'x² + 2x + 1 = 0, найдите x',
        options: ['x = -1', 'x = 1', 'x = 0', 'x = 2'],
        correct: 0,
      },
    ],
  },
  {
    id: 'quiz5',
    title: 'Тест по английскому языку (продвинутый)',
    requiresTeacherCheck: false,
    category: 'eng',
    teacher: 'Владимир Ковалёв',
    questions: [
      {
        id: 'q1',
        type: 'choice',
        question: 'Choose the correct form: "If I ___ rich, I would buy a house."',
        options: ['am', 'was', 'were', 'will be'],
        correct: 2,
      },
    ],
  },
  // Пример теста с открытым вопросом и ручной проверкой
  {
    id: 'quiz6',
    title: 'Смешанный тест',
    requiresTeacherCheck: true,
    category: 'math',
    teacher: 'Иван Петров',
    questions: [
      {
        id: 'q1',
        type: 'choice',
        question: '2 + 2 = ?',
        options: ['3', '4', '5', '6'],
        correct: 1,
      },
      {
        id: 'q2',
        type: 'open',
        question: 'Объясните, почему 2 + 2 = 4.',
      },
    ],
  },
].map(q => ({
  ...q,
  questions: q.questions.map((question: any) => ({
    type: question.type || 'choice',
    ...question
  }))
}));

// Тестовые карточки
export const flashcards: Flashcard[] = [
  {
    id: 'fc1',
    question: 'Столица Франции?',
    answer: 'Париж',
    category: 'eng',
    teacher: 'Ольга Сидорова',
  },
  {
    id: 'fc2',
    question: '3 * 3 = ?',
    answer: '9',
    category: 'math',
    teacher: 'Иван Петров',
  },
  {
    id: 'fc3',
    question: 'Символ химического элемента Водород?',
    answer: 'H',
    category: 'phys',
    teacher: 'Марина Громова',
  },
  {
    id: 'fc4',
    question: 'Как переводится слово "apple"?',
    answer: 'Яблоко',
    category: 'eng',
    teacher: 'Ольга Сидорова',
  },
  {
    id: 'fc5',
    question: '5 + 7 = ?',
    answer: '12',
    category: 'math',
    teacher: 'Иван Петров',
  },
  {
    id: 'fc6',
    question: 'Сколько минут в одном часе?',
    answer: '60',
    category: 'phys',
    teacher: 'Марина Громова',
  },
]; 