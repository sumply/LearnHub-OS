// Типы предметов и групп
export interface Subject {
  id: string;
  name: string;
}

export interface Group {
  id: string;
  name: string;
  studentIds: string[];
}

// Тестовые предметы
export const subjects: Subject[] = [
  { id: 'math', name: 'Математика' },
  { id: 'phys', name: 'Физика' },
  { id: 'rus', name: 'Русский язык' },
  { id: 'eng', name: 'Английский язык' },
  { id: 'db', name: 'Базы данных' },
  { id: 'os', name: 'Операционные системы' },
  { id: 'networks', name: 'Компьютерные сети' },
  { id: 'games', name: 'Игровые технологии' },
];

// Тестовые группы
export const groups: Group[] = [
  { id: 'A', name: '11А', studentIds: ['4'] }, // Петя Иванов
  { id: 'B', name: '11Б', studentIds: [] },
]; 