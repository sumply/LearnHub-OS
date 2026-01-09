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

// Временные тесты удалены - теперь используем API
// Данные теперь загружаются через quizService.getQuizzes()
export const quizzes: Quiz[] = [];

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