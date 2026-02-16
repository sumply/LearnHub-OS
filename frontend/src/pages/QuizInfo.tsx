import { type Component, createSignal, onMount, Show } from 'solid-js';
import { useParams, useNavigate, A } from '@solidjs/router';
import { getQuizzes, type QuizShortResponse } from '../services/quizService';
import { getCurrentUser } from '../utils/api';
import Header from '../components/Header';

const QuizInfo: Component = () => {
  const params = useParams();
  const navigate = useNavigate();
  const quizId = params.quizId;
  
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal<string | null>(null);
  const [quiz, setQuiz] = createSignal<QuizShortResponse | null>(null);
  
  const user = getCurrentUser();
  
  onMount(async () => {
    if (!quizId) {
      setError('ID квиза не указан');
      setLoading(false);
      return;
    }
    
    try {
      setLoading(true);
      setError(null);
      const quizzes = await getQuizzes();
      const foundQuiz = quizzes.find(q => String(q.id) === quizId);
      
      if (!foundQuiz) {
        setError('Квиз не найден');
        return;
      }
      
      setQuiz(foundQuiz);
    } catch (err) {
      console.error('Ошибка загрузки квиза:', err);
      setError('Не удалось загрузить информацию о квизе: ' + (err instanceof Error ? err.message : 'Неизвестная ошибка'));
    } finally {
      setLoading(false);
    }
  });

  const handleStartAttempt = () => {
    if (!quizId) return;
    navigate(`/quiz/take/${quizId}`);
  };

  return (
    <>
      <Header />
      <div class="main-shell" style={{ 'max-width': '1200px', marginTop: 'clamp(100px, 12vh, 140px)', marginBottom: '2rem', marginLeft: 'auto', marginRight: 'auto', padding: '0 1.5rem' }}>
        {/* Хлебные крошки */}
        <nav style={{ 'margin-bottom': '2.5rem', color: 'var(--text-secondary)', fontSize: '0.95em', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
          <A href="/tasks" style={{ color: 'var(--text-secondary)', textDecoration: 'none', transition: 'color 0.2s' }} 
             onMouseOver={(e) => e.currentTarget.style.color = '#2563eb'}
             onMouseOut={(e) => e.currentTarget.style.color = 'var(--text-secondary)'}>
            Интерактивные задания
          </A>
          <span style={{ color: 'var(--text-muted)' }}>/</span>
          <span style={{ color: 'var(--text-primary)', fontWeight: 500 }}>
            {quiz()?.title || 'Загрузка...'}
          </span>
        </nav>

        <Show when={loading()}>
          <div style={{ 'text-align': 'center', padding: '2rem' }}>Загрузка информации о квизе...</div>
        </Show>

        <Show when={error()}>
          <div style={{ color: '#e76f51', padding: '1rem', background: 'rgba(231, 111, 81, 0.1)', 'border-radius': '8px', 'margin-bottom': '1rem' }}>
            {error()}
          </div>
          <button
            onClick={() => navigate('/tasks')}
            style={{
              background: '#2563eb',
              color: '#fff',
              border: 'none',
              borderRadius: '8px',
              padding: '0.75em 2em',
              fontWeight: 600,
              cursor: 'pointer'
            }}
          >
            Вернуться к заданиям
          </button>
        </Show>

        <Show when={!loading() && !error() && quiz()}>
          {(() => {
            const quizData = quiz()!;
            
            return (
              <div>
                {/* Основная карточка квиза */}
                <div style={{ 
                  background: 'var(--bg-primary)', 
                  padding: '3rem 2.5rem', 
                  'border-radius': '24px', 
                  'margin-bottom': '2rem', 
                  'box-shadow': '0 8px 32px rgba(0,0,0,0.08), 0 2px 8px rgba(0,0,0,0.04)',
                  border: '1px solid var(--border-color)',
                  transition: 'box-shadow 0.3s ease'
                }}>
                  <h1 style={{ 
                    color: 'var(--text-primary)', 
                    'margin-bottom': '1rem', 
                    fontSize: 'clamp(1.75rem, 4vw, 2.5rem)', 
                    fontWeight: 700,
                    lineHeight: '1.2',
                    letterSpacing: '-0.02em'
                  }}>
                    {quizData.title}
                  </h1>
                  
                  <div style={{ 
                    color: 'var(--text-secondary)', 
                    'margin-bottom': '2.5rem', 
                    fontSize: '1.15em', 
                    lineHeight: '1.7',
                    maxWidth: '800px'
                  }}>
                    {quizData.summary}
                  </div>

                  {/* Современные информационные карточки */}
                  <div style={{ 
                    display: 'grid', 
                    gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))',
                    gap: '1.25rem', 
                    'margin-bottom': '2.5rem'
                  }}>
                    {/* Карточка предмета */}
                    <div style={{ 
                      background: 'linear-gradient(135deg, rgba(37,99,235,0.08) 0%, rgba(37,99,235,0.03) 100%)',
                      padding: '1.5rem', 
                      'border-radius': '16px',
                      border: '1.5px solid rgba(37,99,235,0.15)',
                      transition: 'all 0.3s ease',
                      position: 'relative',
                      overflow: 'hidden'
                    }}
                    onMouseOver={(e) => {
                      e.currentTarget.style.transform = 'translateY(-2px)';
                      e.currentTarget.style.boxShadow = '0 8px 20px rgba(37,99,235,0.15)';
                    }}
                    onMouseOut={(e) => {
                      e.currentTarget.style.transform = 'translateY(0)';
                      e.currentTarget.style.boxShadow = 'none';
                    }}>
                      <div style={{ 
                        color: '#2563eb', 
                        fontSize: '0.85em', 
                        fontWeight: 600,
                        'margin-bottom': '0.5rem',
                        textTransform: 'uppercase',
                        letterSpacing: '0.05em'
                      }}>
                        Предмет
                      </div>
                      <div style={{ 
                        color: 'var(--text-primary)', 
                        fontWeight: 700, 
                        fontSize: '1.3em',
                        lineHeight: '1.3'
                      }}>
                        {quizData.subject?.name || 'Не указан'}
                      </div>
                    </div>
                    
                    {/* Карточка преподавателя */}
                    <div style={{ 
                      background: 'linear-gradient(135deg, rgba(139,92,246,0.08) 0%, rgba(139,92,246,0.03) 100%)',
                      padding: '1.5rem', 
                      'border-radius': '16px',
                      border: '1.5px solid rgba(139,92,246,0.15)',
                      transition: 'all 0.3s ease',
                      position: 'relative',
                      overflow: 'hidden'
                    }}
                    onMouseOver={(e) => {
                      e.currentTarget.style.transform = 'translateY(-2px)';
                      e.currentTarget.style.boxShadow = '0 8px 20px rgba(139,92,246,0.15)';
                    }}
                    onMouseOut={(e) => {
                      e.currentTarget.style.transform = 'translateY(0)';
                      e.currentTarget.style.boxShadow = 'none';
                    }}>
                      <div style={{ 
                        color: '#8b5cf6', 
                        fontSize: '0.85em', 
                        fontWeight: 600,
                        'margin-bottom': '0.5rem',
                        textTransform: 'uppercase',
                        letterSpacing: '0.05em'
                      }}>
                        Преподаватель
                      </div>
                      <div style={{ 
                        color: 'var(--text-primary)', 
                        fontWeight: 700, 
                        fontSize: '1.3em',
                        lineHeight: '1.3'
                      }}>
                        {quizData.owner?.short_name || 'Неизвестно'}
                      </div>
                    </div>
                    
                    {/* Карточка баллов */}
                    <div style={{ 
                      background: 'linear-gradient(135deg, rgba(16,185,129,0.08) 0%, rgba(16,185,129,0.03) 100%)',
                      padding: '1.5rem', 
                      'border-radius': '16px',
                      border: '1.5px solid rgba(16,185,129,0.15)',
                      transition: 'all 0.3s ease',
                      position: 'relative',
                      overflow: 'hidden'
                    }}
                    onMouseOver={(e) => {
                      e.currentTarget.style.transform = 'translateY(-2px)';
                      e.currentTarget.style.boxShadow = '0 8px 20px rgba(16,185,129,0.15)';
                    }}
                    onMouseOut={(e) => {
                      e.currentTarget.style.transform = 'translateY(0)';
                      e.currentTarget.style.boxShadow = 'none';
                    }}>
                      <div style={{ 
                        color: '#10b981', 
                        fontSize: '0.85em', 
                        fontWeight: 600,
                        'margin-bottom': '0.5rem',
                        textTransform: 'uppercase',
                        letterSpacing: '0.05em'
                      }}>
                        Максимум баллов
                      </div>
                      <div style={{ 
                        color: 'var(--text-primary)', 
                        fontWeight: 700, 
                        fontSize: '1.3em',
                        lineHeight: '1.3'
                      }}>
                        {quizData.total_score}
                      </div>
                    </div>
                  </div>

                  <Show when={user && user.role === 'student'}>
                    <div style={{ 
                      'text-align': 'center', 
                      'margin-top': '2.5rem',
                      paddingTop: '2rem',
                      borderTop: '1px solid var(--border-color)'
                    }}>
                      <button
                        onClick={handleStartAttempt}
                        style={{
                          background: 'linear-gradient(135deg, #2563eb 0%, #1e4ed8 100%)',
                          color: '#fff',
                          border: 'none',
                          borderRadius: '14px',
                          padding: '1.1em 3.5em',
                          fontWeight: 700,
                          fontSize: '1.15em',
                          cursor: 'pointer',
                          transition: 'all 0.3s ease',
                          'box-shadow': '0 6px 20px rgba(37,99,235,0.35), 0 2px 8px rgba(37,99,235,0.2)',
                          position: 'relative',
                          overflow: 'hidden'
                        }}
                        onMouseOver={(e) => {
                          e.currentTarget.style.transform = 'translateY(-3px) scale(1.02)';
                          e.currentTarget.style.boxShadow = '0 8px 24px rgba(37,99,235,0.45), 0 4px 12px rgba(37,99,235,0.25)';
                        }}
                        onMouseOut={(e) => {
                          e.currentTarget.style.transform = 'translateY(0) scale(1)';
                          e.currentTarget.style.boxShadow = '0 6px 20px rgba(37,99,235,0.35), 0 2px 8px rgba(37,99,235,0.2)';
                        }}
                      >
                        Начать попытку
                      </button>
                    </div>
                  </Show>
                </div>
              </div>
            );
          })()}
        </Show>
      </div>
    </>
  );
};

export default QuizInfo;
