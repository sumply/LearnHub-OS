import { type Component, createSignal, onMount, For, Show } from 'solid-js';
import { useParams, useNavigate, A } from '@solidjs/router';
import { startQuizAttempt, finishQuizAttempt, type StartAttemptResponse, type FinishAttemptResponse, type QuizQuestionContent } from '../utils/apiClient';
import { getCurrentUser } from '../utils/api';
import Header from '../components/Header';

const QuizTake: Component = () => {
  const params = useParams();
  const navigate = useNavigate();
  const quizId = params.quizId;
  
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);
  const [quizData, setQuizData] = createSignal<StartAttemptResponse | null>(null);
  const [currentStep, setCurrentStep] = createSignal(0);
  const [answers, setAnswers] = createSignal<Map<string, string | string[] | number>>(new Map());
  const [showResult, setShowResult] = createSignal(false);
  const [resultData, setResultData] = createSignal<FinishAttemptResponse | null>(null);
  
  const user = getCurrentUser();
  
  // Проверка доступа
  if (!user || user.role !== 'student') {
    return (
      <>
        <Header />
        <div style={{ 'text-align': 'center', 'margin-top': '2rem', color: '#e76f51' }}>
          Только студенты могут проходить квизы
        </div>
      </>
    );
  }

  onMount(async () => {
    if (!quizId) {
      setError('ID квиза не указан');
      return;
    }
    
    await startAttempt();
  });

  const startAttempt = async () => {
    if (!quizId) return;
    
    try {
      setLoading(true);
      setError(null);
      const response = await startQuizAttempt(quizId);
      setQuizData(response);
      setCurrentStep(0);
      setAnswers(new Map());
      setShowResult(false);
    } catch (err) {
      console.error('Ошибка начала попытки:', err);
      setError('Не удалось начать попытку: ' + (err instanceof Error ? err.message : 'Неизвестная ошибка'));
    } finally {
      setLoading(false);
    }
  };

  const getCurrentQuestion = (): QuizQuestionContent | null => {
    const data = quizData();
    if (!data || !data.quiz.content || data.quiz.content.length === 0) return null;
    return data.quiz.content[currentStep()] || null;
  };

  const handleAnswerChange = (questionId: string, value: string | string[] | number) => {
    setAnswers(prev => {
      const newMap = new Map(prev);
      newMap.set(questionId, value);
      return newMap;
    });
  };

  const handleSingleChoice = (questionId: string, option: string) => {
    handleAnswerChange(questionId, option);
  };

  const handleMultipleChoice = (questionId: string, option: string) => {
    const currentAnswers = answers().get(questionId);
    const currentArray = Array.isArray(currentAnswers) ? currentAnswers : [];
    
    if (currentArray.includes(option)) {
      // Убираем вариант
      handleAnswerChange(questionId, currentArray.filter(a => a !== option));
    } else {
      // Добавляем вариант
      handleAnswerChange(questionId, [...currentArray, option]);
    }
  };

  const handleNumericAnswer = (questionId: string, value: string) => {
    const numValue = parseFloat(value);
    if (!isNaN(numValue)) {
      handleAnswerChange(questionId, numValue);
    } else if (value === '') {
      answers().delete(questionId);
      setAnswers(new Map(answers()));
    }
  };

  const goToNext = () => {
    const data = quizData();
    if (!data) return;
    
    if (currentStep() < data.quiz.content.length - 1) {
      setCurrentStep(currentStep() + 1);
    }
  };

  const goToPrevious = () => {
    if (currentStep() > 0) {
      setCurrentStep(currentStep() - 1);
    }
  };

  const finishAttempt = async () => {
    const data = quizData();
    if (!data) return;

    // Проверяем, что все вопросы отвечены
    const unansweredQuestions = data.quiz.content.filter(q => !answers().has(q.id));
    if (unansweredQuestions.length > 0) {
      const confirmFinish = confirm(`Вы не ответили на ${unansweredQuestions.length} вопрос(ов). Завершить попытку?`);
      if (!confirmFinish) return;
    }

    try {
      setLoading(true);
      setError(null);
      
      // Подготавливаем ответы для отправки на сервер
      const answersArray = Array.from(answers().entries()).map(([questionId, answer]) => ({
        question_id: questionId,
        answer: answer
      }));
      
      const response = await finishQuizAttempt(data.attempt.id, answersArray, data.quiz);
      setResultData(response);
      setShowResult(true);
    } catch (err) {
      console.error('Ошибка завершения попытки:', err);
      setError('Не удалось завершить попытку: ' + (err instanceof Error ? err.message : 'Неизвестная ошибка'));
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
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
        {/* Хлебные крошки */}
        <nav style={{ 'margin-bottom': '2rem', color: 'var(--text-secondary)', fontSize: '0.9em', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
          <A href="/tasks" style={{ color: 'var(--text-secondary)', textDecoration: 'none', transition: 'color 0.2s' }} 
             onMouseOver={(e) => e.currentTarget.style.color = '#2563eb'}
             onMouseOut={(e) => e.currentTarget.style.color = 'var(--text-secondary)'}>
            Интерактивные задания
          </A>
          <span style={{ color: 'var(--text-muted)' }}>/</span>
          <A href={`/quiz/${quizId}`} style={{ color: 'var(--text-secondary)', textDecoration: 'none', transition: 'color 0.2s' }}
             onMouseOver={(e) => e.currentTarget.style.color = '#2563eb'}
             onMouseOut={(e) => e.currentTarget.style.color = 'var(--text-secondary)'}>
            {quizData()?.quiz.title || 'Квиз'}
          </A>
          <span style={{ color: 'var(--text-muted)' }}>/</span>
          <span style={{ color: 'var(--text-primary)', fontWeight: 500 }}>
            Прохождение
          </span>
        </nav>
        <Show when={loading()}>
          <div style={{ 'text-align': 'center', padding: '2rem' }}>Загрузка...</div>
        </Show>

        <Show when={error()}>
          <div style={{ 
            color: '#e76f51', 
            padding: '1.5rem', 
            background: 'rgba(231, 111, 81, 0.1)', 
            'border-radius': '12px', 
            'margin-bottom': '2rem',
            border: '1px solid rgba(231, 111, 81, 0.2)',
            fontSize: '1.05em',
            lineHeight: '1.6'
          }}>
            {error()}
          </div>
          <div style={{ 'text-align': 'center', 'margin-top': '1.5rem' }}>
            <button
              onClick={() => navigate('/tasks')}
              style={{
                background: '#2563eb',
                color: '#fff',
                border: 'none',
                borderRadius: '8px',
                padding: '0.75em 2em',
                fontWeight: 600,
                cursor: 'pointer',
                fontSize: '1em'
              }}
            >
              Вернуться к заданиям
            </button>
          </div>
        </Show>

        <Show when={!loading() && !error() && quizData() && !showResult()}>
          {(() => {
            const data = quizData()!;
            const question = getCurrentQuestion();
            
            return (
              <>
                <div style={{ 'margin-bottom': '2rem' }}>
                  <h2 style={{ color: 'var(--text-primary)', 'margin-bottom': '0.5rem' }}>{data.quiz.title}</h2>
                  <div style={{ color: 'var(--text-secondary)', 'margin-bottom': '1rem' }}>
                    {data.quiz.summary}
                  </div>
                  <div style={{ color: 'var(--text-secondary)', fontSize: '0.9em' }}>
                    Предмет: {data.quiz.subject.name} | Максимум попыток: {data.quiz.max_attempts} | 
                    Всего баллов: {data.quiz.total_score}
                  </div>
                </div>

                <div style={{ 'margin-bottom': '1.5rem', color: 'var(--text-secondary)' }}>
                  Вопрос {currentStep() + 1} из {data.quiz.content.length}
                </div>

                <Show when={question}>
                  {(() => {
                    const q = question!;
                    const currentAnswer = answers().get(q.id);
                    
                    return (
                      <div style={{ background: 'var(--bg-secondary)', padding: '2rem', 'border-radius': '12px', 'margin-bottom': '2rem' }}>
                        <h3 style={{ color: 'var(--text-primary)', 'margin-bottom': '1.5rem', fontSize: '1.3em' }}>
                          {q.text}
                        </h3>
                        <div style={{ color: 'var(--text-secondary)', 'margin-bottom': '1.5rem', fontSize: '0.9em' }}>
                          Баллов за вопрос: {q.score}
                        </div>

                        {q.type === 'single' && (
                          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.8em' }}>
                            {q.details.options.map((optionValue, idx) => {
                              const isSelected = currentAnswer === optionValue;
                              
                              return (
                                <label 
                                  key={`${q.id}-option-${idx}`}
                                  style={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    padding: '1em',
                                    border: '2px solid',
                                    borderColor: isSelected ? '#2563eb' : '#e3eafc',
                                    borderRadius: '8px',
                                    cursor: 'pointer',
                                    background: isSelected ? 'var(--bg-tertiary)' : 'var(--bg-primary)',
                                    color: 'var(--text-primary)',
                                    transition: 'all 0.2s'
                                  }}
                                  onClick={(e) => {
                                    e.preventDefault();
                                    e.stopPropagation();
                                    // Используем значение напрямую из массива по индексу
                                    const selectedValue = q.details.options[idx];
                                    handleSingleChoice(q.id, selectedValue);
                                  }}
                                >
                                  <input
                                    type="radio"
                                    name={`question-${q.id}`}
                                    value={optionValue}
                                    checked={isSelected}
                                    readOnly
                                    style={{ marginRight: '0.8em', width: '18px', height: '18px', cursor: 'pointer' }}
                                  />
                                  <span style={{ flex: 1 }}>{optionValue}</span>
                                </label>
                              );
                            })}
                          </div>
                        )}

                        {q.type === 'multiple' && (
                          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.8em' }}>
                            <For each={q.details.options}>
                              {(option) => {
                                const isSelected = Array.isArray(currentAnswer) && currentAnswer.includes(option);
                                return (
                                  <label style={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    padding: '1em',
                                    border: '2px solid',
                                    borderColor: isSelected ? '#2563eb' : '#e3eafc',
                                    borderRadius: '8px',
                                    cursor: 'pointer',
                                    background: isSelected ? 'var(--bg-tertiary)' : 'var(--bg-primary)',
                                    color: 'var(--text-primary)',
                                    transition: 'all 0.2s'
                                  }}>
                                    <input
                                      type="checkbox"
                                      checked={isSelected}
                                      onChange={() => handleMultipleChoice(q.id, option)}
                                      style={{ marginRight: '0.8em', width: '18px', height: '18px', cursor: 'pointer' }}
                                    />
                                    <span style={{ flex: 1 }}>{option}</span>
                                  </label>
                                );
                              }}
                            </For>
                          </div>
                        )}

                        {q.type === 'numeric' && (
                          <div>
                            <input
                              type="number"
                              step="any"
                              value={typeof currentAnswer === 'number' ? currentAnswer : ''}
                              onInput={(e) => handleNumericAnswer(q.id, e.currentTarget.value)}
                              placeholder="Введите числовой ответ"
                              style={{
                                width: '100%',
                                padding: '0.8em',
                                border: '2px solid #e3eafc',
                                borderRadius: '8px',
                                fontSize: '1em',
                                fontFamily: 'inherit',
                                background: 'var(--bg-primary)',
                                color: 'var(--text-primary)'
                              }}
                            />
                          </div>
                        )}
                      </div>
                    );
                  })()}
                </Show>

                <div style={{ display: 'flex', gap: '1em', justifyContent: 'space-between' }}>
                  <button
                    onClick={goToPrevious}
                    disabled={currentStep() === 0}
                    style={{
                      background: currentStep() === 0 ? '#e3eafc' : '#f3f4f6',
                      color: currentStep() === 0 ? '#9ca3af' : '#213547',
                      border: 'none',
                      borderRadius: '8px',
                      padding: '0.75em 2em',
                      fontWeight: 600,
                      cursor: currentStep() === 0 ? 'not-allowed' : 'pointer',
                      opacity: currentStep() === 0 ? 0.6 : 1
                    }}
                  >
                    Назад
                  </button>
                  
                  {currentStep() < data.quiz.content.length - 1 ? (
                    <button
                      onClick={goToNext}
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
                      Следующий вопрос
                    </button>
                  ) : (
                    <button
                      onClick={finishAttempt}
                      disabled={loading()}
                      style={{
                        background: '#10b981',
                        color: '#fff',
                        border: 'none',
                        borderRadius: '8px',
                        padding: '0.75em 2em',
                        fontWeight: 600,
                        cursor: loading() ? 'not-allowed' : 'pointer',
                        opacity: loading() ? 0.6 : 1
                      }}
                    >
                      {loading() ? 'Завершение...' : 'Завершить попытку'}
                    </button>
                  )}
                </div>
              </>
            );
          })()}
        </Show>

        <Show when={showResult() && resultData()}>
          {(() => {
            const result = resultData()!;
            const percentage = result.quiz.total_score > 0 
              ? Math.round((result.attempt.score / result.quiz.total_score) * 100)
              : 0;
            
            return (
              <div style={{ 'text-align': 'center', 'max-width': '800px', margin: '2rem auto' }}>
                <h2 style={{ color: 'var(--text-primary)', 'margin-bottom': '1rem' }}>Результат попытки</h2>
                
                <div style={{ background: 'var(--bg-secondary)', padding: '2rem', 'border-radius': '12px', 'margin-bottom': '2rem' }}>
                  <div style={{ fontSize: '2em', fontWeight: 700, color: '#2563eb', 'margin-bottom': '1rem' }}>
                    {result.attempt.score} / {result.quiz.total_score}
                  </div>
                  <div style={{ fontSize: '1.2em', color: 'var(--text-secondary)', 'margin-bottom': '1rem' }}>
                    {percentage}%
                  </div>
                  <div style={{ color: 'var(--text-secondary)', fontSize: '0.9em' }}>
                    Начало: {formatDate(result.attempt.started_at)}<br />
                    Завершение: {formatDate(result.attempt.ended_at)}
                  </div>
                </div>

                <div style={{ 'margin-bottom': '2rem', 'text-align': 'left' }}>
                  <h3 style={{ color: 'var(--text-primary)', 'margin-bottom': '1.5rem', fontSize: '1.3em', fontWeight: 700 }}>Детали ответов:</h3>
                  <For each={result.attempt.answers}>
                    {(answer) => {
                      const question = result.quiz.content.find(q => q.id === answer.question_id);
                      if (!question) return null;
                      
                      // Получаем правильный ответ в зависимости от типа вопроса
                      let correctAnswer: string | string[] | number = '';
                      if (question.type === 'single') {
                        correctAnswer = question.details.correct;
                      } else if (question.type === 'multiple') {
                        correctAnswer = question.details.correct;
                      } else if (question.type === 'numeric') {
                        correctAnswer = question.details.correct;
                      }
                      
                      // Форматируем ответы для отображения
                      const formatAnswer = (ans: string | string[] | number): string => {
                        if (Array.isArray(ans)) {
                          return ans.length > 0 ? ans.join(', ') : 'Нет ответа';
                        }
                        return String(ans);
                      };
                      
                      const userAnswerFormatted = formatAnswer(answer.answer);
                      const correctAnswerFormatted = formatAnswer(correctAnswer);
                      
                      return (
                        <div style={{ 
                          background: 'var(--bg-secondary)', 
                          padding: '1.5rem', 
                          'border-radius': '12px', 
                          'margin-bottom': '1rem',
                          border: `2px solid ${answer.is_correct ? '#10b981' : '#e76f51'}`,
                          boxShadow: `0 2px 8px ${answer.is_correct ? 'rgba(16,185,129,0.15)' : 'rgba(231,111,81,0.15)'}`
                        }}>
                          <div style={{ color: 'var(--text-primary)', fontWeight: 600, 'margin-bottom': '1rem', fontSize: '1.1em' }}>
                            {question.text}
                          </div>
                          
                          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem', 'margin-bottom': '1rem' }}>
                            <div>
                              <div style={{ 
                                color: 'var(--text-secondary)', 
                                fontSize: '0.85em', 
                                fontWeight: 600,
                                'margin-bottom': '0.3rem',
                                textTransform: 'uppercase',
                                letterSpacing: '0.05em'
                              }}>
                                Ваш ответ:
                              </div>
                              <div style={{ 
                                color: answer.is_correct ? '#10b981' : '#e76f51',
                                fontWeight: 600,
                                fontSize: '1em',
                                padding: '0.5rem 0.75rem',
                                background: answer.is_correct ? 'rgba(16,185,129,0.1)' : 'rgba(231,111,81,0.1)',
                                borderRadius: '6px',
                                display: 'inline-block'
                              }}>
                                {userAnswerFormatted || 'Нет ответа'}
                              </div>
                            </div>
                            
                            <div>
                              <div style={{ 
                                color: 'var(--text-secondary)', 
                                fontSize: '0.85em', 
                                fontWeight: 600,
                                'margin-bottom': '0.3rem',
                                textTransform: 'uppercase',
                                letterSpacing: '0.05em'
                              }}>
                                Правильный ответ:
                              </div>
                              <div style={{ 
                                color: '#10b981',
                                fontWeight: 600,
                                fontSize: '1em',
                                padding: '0.5rem 0.75rem',
                                background: 'rgba(16,185,129,0.1)',
                                borderRadius: '6px',
                                display: 'inline-block'
                              }}>
                                {correctAnswerFormatted}
                              </div>
                            </div>
                          </div>
                          
                          <div style={{ 
                            display: 'flex',
                            alignItems: 'center',
                            gap: '1rem',
                            paddingTop: '0.75rem',
                            borderTop: '1px solid var(--border-color)'
                          }}>
                            <div style={{ 
                              color: answer.is_correct ? '#10b981' : '#e76f51', 
                              fontWeight: 700,
                              fontSize: '1em',
                              display: 'flex',
                              alignItems: 'center',
                              gap: '0.5rem'
                            }}>
                              {answer.is_correct ? (
                                <>
                                  <span style={{ fontSize: '1.2em' }}>✓</span>
                                  <span>Правильно</span>
                                </>
                              ) : (
                                <>
                                  <span style={{ fontSize: '1.2em' }}>✗</span>
                                  <span>Неправильно</span>
                                </>
                              )}
                            </div>
                            <div style={{ 
                              color: 'var(--text-secondary)', 
                              fontSize: '0.9em',
                              marginLeft: 'auto'
                            }}>
                              Баллы: <strong style={{ color: 'var(--text-primary)' }}>{answer.score}</strong> / {question.score}
                            </div>
                          </div>
                        </div>
                      );
                    }}
                  </For>
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
                    cursor: 'pointer',
                    fontSize: '1em'
                  }}
                >
                  Вернуться к заданиям
                </button>
              </div>
            );
          })()}
        </Show>
      </div>
    </>
  );
};

export default QuizTake;
