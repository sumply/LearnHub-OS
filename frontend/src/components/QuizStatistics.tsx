import { type Component, createSignal, For, Show, createResource } from 'solid-js';
import { getUserQuizzes, getUsers, type UserShort } from '../utils/apiClient';
import { getCurrentUser } from '../utils/api';

interface QuizStatisticsProps {
  quizId: string;
  quizTitle: string;
  onClose: () => void;
}

interface UserWithAttempt {
  id: string;
  first_name: string;
  last_name: string;
  role: string;
  last_attempt: {
    id: string;
    score: number;
    started_at: string;
    ended_at: string | null;
  } | null;
}

const QuizStatistics: Component<QuizStatisticsProps> = (props) => {
  const [error, setError] = createSignal<string | null>(null);
  
  const user = getCurrentUser();

  // Загружаем студентов с их попытками для квиза через /users/{user_id}/quizzes
  const [usersData, { refetch }] = createResource(
    async () => {
      try {
        setError(null);
        
        // Получаем всех пользователей
        const allUsers = await getUsers();
        
        // Фильтруем только студентов
        const students = allUsers.filter(u => u.role === 'student');
        
        // Для каждого студента получаем его квизы через /users/{user_id}/quizzes
        const usersWithAttempts: UserWithAttempt[] = await Promise.all(
          students.map(async (student) => {
            try {
              const quizzes = await getUserQuizzes(String(student.id));
              
              // Ищем нужный квиз в списке квизов студента
              const quiz = quizzes.find((q: any) => String(q.id) === String(props.quizId));
              
              // Если квиз найден и есть информация о попытке, используем её
              const lastAttempt = quiz?.last_attempt || null;
              
              return {
                id: student.id,
                first_name: student.first_name,
                last_name: student.last_name,
                role: student.role,
                last_attempt: lastAttempt ? {
                  id: lastAttempt.id,
                  score: lastAttempt.score,
                  started_at: lastAttempt.started_at,
                  ended_at: lastAttempt.ended_at
                } : null
              };
            } catch (err) {
              console.error(`Ошибка загрузки квизов для студента ${student.id}:`, err);
              return {
                id: student.id,
                first_name: student.first_name,
                last_name: student.last_name,
                role: student.role,
                last_attempt: null
              };
            }
          })
        );
        
        return usersWithAttempts;
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Ошибка загрузки статистики';
        setError(errorMessage);
        return [];
      }
    }
  );

  const formatDate = (dateString: string | null) => {
    if (!dateString) {
      return 'Не завершено';
    }
    try {
      const date = new Date(dateString);
      return date.toLocaleString('ru-RU', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      });
    } catch {
      return dateString;
    }
  };

  return (
    <div
      style={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        background: 'rgba(0, 0, 0, 0.5)',
        display: 'flex',
        'align-items': 'center',
        'justify-content': 'center',
        'z-index': 1000,
        padding: '2rem',
        overflow: 'auto',
      }}
      onClick={(e) => {
        if (e.target === e.currentTarget) {
          props.onClose();
        }
      }}
    >
      <div
        style={{
          background: 'var(--bg-primary)',
          'border-radius': '16px',
          padding: '2rem',
          'max-width': '1200px',
          width: '100%',
          'max-height': '90vh',
          overflow: 'auto',
          'box-shadow': '0 8px 32px rgba(0, 0, 0, 0.2)',
        }}
        onClick={(e) => e.stopPropagation()}
      >
        <div style={{ display: 'flex', 'justify-content': 'space-between', 'align-items': 'center', 'margin-bottom': '1.5rem' }}>
          <h2 style={{ margin: 0, color: 'var(--text-primary)' }}>Статистика квиза: {props.quizTitle}</h2>
          <button
            onClick={props.onClose}
            style={{
              background: 'var(--accent-danger)',
              color: '#fff',
              border: 'none',
              'border-radius': '8px',
              padding: '0.5rem 1rem',
              cursor: 'pointer',
              'font-size': '1rem',
            }}
          >
            ✕ Закрыть
          </button>
        </div>

        <Show when={error()}>
          <div style={{ color: 'var(--accent-danger)', 'margin-bottom': '1rem', padding: '1rem', background: 'rgba(220, 53, 69, 0.1)', 'border-radius': '8px' }}>
            {error()}
          </div>
        </Show>

        <Show when={usersData.loading}>
          <div style={{ 'text-align': 'center', padding: '2rem' }}>Загрузка статистики...</div>
        </Show>

        <Show when={!usersData.loading && !error()}>
          {/* Таблица статистики */}
          <div style={{ overflow: 'auto' }}>
            <table
              style={{
                width: '100%',
                'border-collapse': 'collapse',
                background: 'var(--bg-primary)',
              }}
            >
              <thead>
                <tr style={{ background: 'var(--bg-tertiary)', 'border-bottom': '2px solid var(--border-color)' }}>
                  <th style={{ padding: '0.75rem', 'text-align': 'left', color: 'var(--text-primary)' }}>Ученик</th>
                  <th style={{ padding: '0.75rem', 'text-align': 'center', color: 'var(--text-primary)' }}>Баллы</th>
                  <th style={{ padding: '0.75rem', 'text-align': 'center', color: 'var(--text-primary)' }}>Начало попытки</th>
                  <th style={{ padding: '0.75rem', 'text-align': 'center', color: 'var(--text-primary)' }}>Завершение попытки</th>
                </tr>
              </thead>
              <tbody>
                <Show when={usersData() && usersData()!.length > 0} fallback={
                  <tr>
                    <td colSpan={4} style={{ padding: '2rem', 'text-align': 'center', color: 'var(--text-muted)' }}>
                      Нет студентов, проходивших этот квиз
                    </td>
                  </tr>
                }>
                  <For each={usersData()}>
                    {(userData) => {
                      const fullName = `${userData.last_name} ${userData.first_name}`;
                      const hasAttempt = userData.last_attempt !== null;

                      return (
                        <tr style={{ 'border-bottom': '1px solid var(--border-color)' }}>
                          <td style={{ padding: '0.75rem', color: 'var(--text-primary)' }}>
                            {fullName}
                          </td>
                          <td style={{ padding: '0.75rem', 'text-align': 'center', color: 'var(--text-primary)' }}>
                            {hasAttempt ? userData.last_attempt!.score : (
                              <span style={{ color: 'var(--text-muted)' }}>Попыток не было</span>
                            )}
                          </td>
                          <td style={{ padding: '0.75rem', 'text-align': 'center', color: 'var(--text-primary)' }}>
                            {hasAttempt ? formatDate(userData.last_attempt!.started_at) : (
                              <span style={{ color: 'var(--text-muted)' }}>—</span>
                            )}
                          </td>
                          <td style={{ padding: '0.75rem', 'text-align': 'center', color: 'var(--text-primary)' }}>
                            {hasAttempt ? formatDate(userData.last_attempt!.ended_at) : (
                              <span style={{ color: 'var(--text-muted)' }}>—</span>
                            )}
                          </td>
                        </tr>
                      );
                    }}
                  </For>
                </Show>
              </tbody>
            </table>
          </div>
        </Show>
      </div>
    </div>
  );
};

export default QuizStatistics;
