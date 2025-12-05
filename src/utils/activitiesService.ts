import { quizzes, flashcards, type Quiz, type Flashcard } from '../config/activities';

export type TaskUnion = ((Quiz & { type: 'quiz' }) | (Flashcard & { type: 'flashcard' })) & { hidden?: boolean; requiresConfirmation?: boolean };

const STORAGE_KEY = 'activities';

function getDefaultTasks(): TaskUnion[] {
  return [
    ...quizzes.map(q => ({ ...q, type: 'quiz' as const })),
    ...flashcards.map(f => ({ ...f, type: 'flashcard' as const })),
  ];
}

export function getAllActivities(): TaskUnion[] {
  const data = localStorage.getItem(STORAGE_KEY);
  return data ? JSON.parse(data) : getDefaultTasks();
}

export function addActivity(activity: TaskUnion): void {
  const activities = getAllActivities();
  activities.push({ ...activity, hidden: activity.hidden ?? false });
  localStorage.setItem(STORAGE_KEY, JSON.stringify(activities));
}

export function updateActivity(updatedActivity: TaskUnion): void {
  const activities = getAllActivities().map(a =>
    a.id === updatedActivity.id ? { ...updatedActivity, hidden: updatedActivity.hidden ?? false } : a
  );
  localStorage.setItem(STORAGE_KEY, JSON.stringify(activities));
}

export function removeActivity(id: string): void {
  const activities = getAllActivities().filter(a => a.id !== id);
  localStorage.setItem(STORAGE_KEY, JSON.stringify(activities));
} 