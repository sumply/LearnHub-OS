import { type Component, createSignal, onMount, For, Show } from 'solid-js';
import { getCurrentUser } from '../utils/api';
import { getUserQuizzes } from '../services/quizService';
import { getAuthData, type QuizShortResponse } from '../utils/apiClient';
import Header from '../components/Header';

interface StudentAttempt {
  quizId: string;
  quizTitle: string;
  lastAttempt: {
    id: string;
    score: number;
    started_at: string;
    ended_at: string | null;
  } | null;
}

const Progress: Component = () => {
  const user = getCurrentUser();
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal<string | null>(null);
  const [studentAttempts, setStudentAttempts] = createSignal<StudentAttempt[]>([]);

  // Проверка доступа
  if (!user || user.role !== 'student') {
    return (
      <>
        <Header />
        <div class="main-shell" style={{ 'max-width': '1200px', marginTop: 'clamp(100px, 12vh, 140px)', marginBottom: '2rem', marginLeft: 'auto', marginRight: 'auto', padding: '0 1.5rem' }}>
          <div style={{ 'text-align': 'center', 'margin-top': '2rem', color: '#e76f51' }}>
            Только студенты могут просматривать свою успеваемость
          </div>
        </div>
      </>
    );
  }

  onMount(async () => {
    await loadStudentProgress();
  });

  const loadStudentProgress = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const authData = getAuthData();
      if (!authData) {
        setError('Не удалось получить данные пользователя');
        setLoading(false);
        return;
      }
      
      // Получаем квизы пользователя через /users/{user_id}/quizzes
      const quizzes = await getUserQuizzes();
      
      // Преобразуем квизы в формат StudentAttempt
      // Предполагаем, что QuizShortResponse теперь включает информацию о попытках
      // Если структура другая, нужно будет обновить интерфейс
      const attempts: StudentAttempt[] = quizzes.map((quiz: any) => {
        // Если в ответе есть информация о попытках, используем её
        // Иначе возвращаем null для lastAttempt
        const lastAttempt = quiz.last_attempt || null;
        
        return {
          quizId: String(quiz.id),
          quizTitle: quiz.title,
          lastAttempt: lastAttempt ? {
            id: lastAttempt.id,
            score: lastAttempt.score,
            started_at: lastAttempt.started_at,
            ended_at: lastAttempt.ended_at
          } : null
        };
      });
      
      setStudentAttempts(attempts);
    } catch (err) {
      console.error('Ошибка загрузки прогресса:', err);
      setError('Не удалось загрузить прогресс: ' + (err instanceof Error ? err.message : 'Неизвестная ошибка'));
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string | undefined) => {
    if (!dateString || dateString === '0001-01-01T00:00:00Z') {
      return '—';
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
    <>
      <Header />
      <div class="main-shell" style={{ 'max-width': '1200px', marginTop: 'clamp(100px, 12vh, 140px)', marginBottom: '2rem', marginLeft: 'auto', marginRight: 'auto', padding: '0 1.5rem' }}>
        <h1 style={{ color: 'var(--text-primary)', 'margin-bottom': '2rem', fontSize: 'clamp(1.75rem, 4vw, 2.5rem)', fontWeight: 700 }}>
          Моя успеваемость
        </h1>

        <Show when={loading()}>
          <div style={{ 'text-align': 'center', padding: '3rem' }}>
            <div style={{ fontSize: '1.1em', color: 'var(--text-secondary)' }}>Загрузка данных...</div>
          </div>
        </Show>

        <Show when={error()}>
          <div style={{ 
            color: '#e76f51', 
            padding: '1.5rem', 
            background: 'rgba(231, 111, 81, 0.1)', 
            'border-radius': '12px', 
            'margin-bottom': '2rem',
            border: '1px solid rgba(231, 111, 81, 0.2)',
            fontSize: '1.05em'
          }}>
            {error()}
          </div>
        </Show>

        <Show when={!loading() && !error()}>
          <Show when={studentAttempts().length === 0} fallback={
            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
              <For each={studentAttempts()}>
                {(attempt) => (
                  <div style={{ 
                    background: 'var(--bg-secondary)', 
                    padding: '1.5rem', 
                    'border-radius': '16px',
                    border: '1px solid var(--border-color)',
                    boxShadow: '0 4px 12px rgba(0,0,0,0.08)',
                    transition: 'all 0.3s ease'
                  }}
                  onMouseOver={(e) => {
                    e.currentTarget.style.transform = 'translateY(-2px)';
                    e.currentTarget.style.boxShadow = '0 6px 16px rgba(0,0,0,0.12)';
                  }}
                  onMouseOut={(e) => {
                    e.currentTarget.style.transform = 'translateY(0)';
                    e.currentTarget.style.boxShadow = '0 4px 12px rgba(0,0,0,0.08)';
                  }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', 'margin-bottom': '1rem' }}>
                      <h3 style={{ 
                        color: 'var(--text-primary)', 
                        fontWeight: 700, 
                        fontSize: '1.2em',
                        margin: 0,
                        flex: 1
                      }}>
                        {attempt.quizTitle}
                      </h3>
                    </div>
                    
                    <Show when={attempt.lastAttempt} fallback={
                      <div style={{ 
                        color: 'var(--text-secondary)', 
                        fontSize: '1em',
                        padding: '1rem',
                        background: 'rgba(108, 117, 125, 0.1)',
                        borderRadius: '8px',
                        textAlign: 'center'
                      }}>
                        Попыток не было
                      </div>
                    }>
                      {(() => {
                        const lastAttempt = attempt.lastAttempt!;
                        return (
                          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1rem' }}>
                            <div style={{ 
                              background: 'linear-gradient(135deg, rgba(37,99,235,0.08) 0%, rgba(37,99,235,0.03) 100%)',
                              padding: '1rem', 
                              'border-radius': '12px',
                              border: '1.5px solid rgba(37,99,235,0.15)'
                            }}>
                              <div style={{ 
                                color: '#2563eb', 
                                fontSize: '0.85em', 
                                fontWeight: 600,
                                'margin-bottom': '0.5rem',
                                textTransform: 'uppercase',
                                letterSpacing: '0.05em'
                              }}>
                                Баллы
                              </div>
                              <div style={{ 
                                color: 'var(--text-primary)', 
                                fontWeight: 700, 
                                fontSize: '1.5em'
                              }}>
                                {lastAttempt.score}
                              </div>
                            </div>
                            
                            <div style={{ 
                              background: 'linear-gradient(135deg, rgba(16,185,129,0.08) 0%, rgba(16,185,129,0.03) 100%)',
                              padding: '1rem', 
                              'border-radius': '12px',
                              border: '1.5px solid rgba(16,185,129,0.15)'
                            }}>
                              <div style={{ 
                                color: '#10b981', 
                                fontSize: '0.85em', 
                                fontWeight: 600,
                                'margin-bottom': '0.5rem',
                                textTransform: 'uppercase',
                                letterSpacing: '0.05em'
                              }}>
                                Начало попытки
                              </div>
                              <div style={{ 
                                color: 'var(--text-primary)', 
                                fontWeight: 600, 
                                fontSize: '1em',
                                lineHeight: '1.4'
                              }}>
                                {formatDate(lastAttempt.started_at)}
                              </div>
                            </div>
                            
                            <Show when={lastAttempt.ended_at}>
                              <div style={{ 
                                background: 'linear-gradient(135deg, rgba(139,92,246,0.08) 0%, rgba(139,92,246,0.03) 100%)',
                                padding: '1rem', 
                                'border-radius': '12px',
                                border: '1.5px solid rgba(139,92,246,0.15)'
                              }}>
                                <div style={{ 
                                  color: '#8b5cf6', 
                                  fontSize: '0.85em', 
                                  fontWeight: 600,
                                  'margin-bottom': '0.5rem',
                                  textTransform: 'uppercase',
                                  letterSpacing: '0.05em'
                                }}>
                                  Завершение
                                </div>
                                <div style={{ 
                                  color: 'var(--text-primary)', 
                                  fontWeight: 600, 
                                  fontSize: '1em',
                                  lineHeight: '1.4'
                                }}>
                                  {formatDate(lastAttempt.ended_at)}
                                </div>
                              </div>
                            </Show>
                          </div>
                        );
                      })()}
                    </Show>
                  </div>
                )}
              </For>
            </div>
          }>
            <div style={{ 
              'text-align': 'center', 
              padding: '3rem',
              background: 'var(--bg-secondary)',
              borderRadius: '16px',
              border: '1px solid var(--border-color)'
            }}>
              <div style={{ fontSize: '1.2em', color: 'var(--text-secondary)', 'margin-bottom': '0.5rem' }}>
                У вас пока нет попыток прохождения квизов
              </div>
              <div style={{ fontSize: '1em', color: 'var(--text-muted)' }}>
                Начните прохождение квизов, чтобы увидеть здесь свой прогресс
              </div>
            </div>
          </Show>
        </Show>
      </div>
    </>
  );
};

export default Progress;
