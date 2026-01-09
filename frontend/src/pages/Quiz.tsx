import { type Component, createSignal, onMount, For, Show, createResource } from 'solid-js';
import { checkToken, getCurrentUser } from '../utils/api';
import { getQuizzes, getQuizById } from '../services/quizService';
import { getProgress, startProgress, finishProgress, updateAnswer } from '../services/progressService';
import * as apiClient from '../utils/apiClient';

const Quiz: Component = () => {
  const [hasAccess, setHasAccess] = createSignal<boolean | null>(null);
  const [error, setError] = createSignal<string | null>(null);
  const [loading, setLoading] = createSignal(true);
  const [currentQuizShort, setCurrentQuizShort] = createSignal<apiClient.QuizShortResponse | null>(null);
  const [currentQuizFull, setCurrentQuizFull] = createSignal<apiClient.QuizFullResponse | null>(null);
  const [currentProgress, setCurrentProgress] = createSignal<apiClient.QuizProgressResponse | null>(null);
  const [step, setStep] = createSignal(0);
  // Для каждого вопроса: либо индекс выбранного варианта, либо текст ответа
  const [answers, setAnswers] = createSignal<(number | string)[]>([]);
  const [currentChoice, setCurrentChoice] = createSignal<number | null>(null);
  const [currentTextAnswer, setCurrentTextAnswer] = createSignal<string>('');
  const [showResult, setShowResult] = createSignal(false);
  const [score, setScore] = createSignal(0);
  const [loadingQuiz, setLoadingQuiz] = createSignal(false);

  // Загружаем квизы и прогресс через API
  // Используем функцию-обертку, которая будет вызвана только после проверки токена
  const [quizzesData, { refetch: refetchQuizzes }] = createResource(async () => {
    const ok = await checkToken();
    if (!ok) return [];
    return await getQuizzes();
  });
  
  const [progressData, { refetch: refetchProgress }] = createResource(async () => {
    const ok = await checkToken();
    if (!ok) return [];
    return await getProgress();
  });

  onMount(async () => {
    setLoading(true);
    const ok = await checkToken();
    setHasAccess(ok);
    if (!ok) {
      setTimeout(() => {
        window.location.href = '/login';
      }, 2000);
      setLoading(false);
      return;
    }
    
    // Явно загружаем квизы и прогресс после проверки токена
    try {
      await refetchQuizzes();
      await refetchProgress();
    } catch (error) {
      console.error('Ошибка загрузки данных:', error);
      setError('Не удалось загрузить данные');
    }
    
    setLoading(false);
  });

  // Фильтруем квизы по роли пользователя
  const filteredQuizzes = () => {
    const quizzes = quizzesData();
    const user = getCurrentUser();
    if (!quizzes || !user) return [];
    
    // TODO: Реализовать фильтрацию по группам и предметам когда будет доступна полная информация о пользователе
    return quizzes;
  };

  async function startQuiz(quiz: apiClient.QuizShortResponse) {
    try {
      setLoadingQuiz(true);
      // Находим прогресс для этого квиза
      const progress = progressData()?.find(p => p.quiz.id === quiz.id);
      if (!progress) {
        setError('Прогресс для этого квиза не найден');
        setLoadingQuiz(false);
        return;
      }
      
      // Начинаем прохождение
      await startProgress(progress.id);
      
      // Загружаем полный квиз с вопросами
      const fullQuiz = await getQuizById(quiz.id);
      
      setCurrentQuizShort(quiz);
      setCurrentQuizFull(fullQuiz);
      setCurrentProgress(progress);
      setStep(0);
      setAnswers([]);
      setCurrentChoice(null);
      setCurrentTextAnswer('');
      setShowResult(false);
      setScore(0);
      setError(null);
    } catch (error) {
      console.error('Ошибка начала квиза:', error);
      setError('Не удалось начать квиз: ' + (error instanceof Error ? error.message : 'Неизвестная ошибка'));
    } finally {
      setLoadingQuiz(false);
    }
  }

  async function saveAnswer() {
    const quiz = currentQuizFull();
    const progress = currentProgress();
    if (!quiz || !progress) return;
    
    const currentQuestion = quiz.questions[step()];
    if (!currentQuestion) return;
    
    try {
      let answerText = '';
      
      // Если есть варианты ответов, используем выбранный вариант
      if (currentQuestion.options && currentQuestion.options.length > 0) {
        if (currentChoice() !== null) {
          answerText = currentQuestion.options[currentChoice()!].text;
        } else {
          alert('Пожалуйста, выберите вариант ответа');
          return;
        }
      } else {
        // Если вариантов нет - открытый вопрос
        answerText = currentTextAnswer().trim();
        if (!answerText) {
          alert('Пожалуйста, введите ответ');
          return;
        }
      }
      
      // Отправляем ответ (answer_id соответствует индексу вопроса)
      await updateAnswer(progress.id, step(), answerText);
      
      // Сохраняем ответ локально
      setAnswers([...answers(), answerText]);
      setCurrentChoice(null);
      setCurrentTextAnswer('');
      
      // Переходим к следующему вопросу или завершаем
      if (step() < quiz.questions.length - 1) {
        setStep(step() + 1);
      } else {
        // Завершение теста
        await finishProgress(progress.id);
        // Обновляем прогресс для получения финального score
        const updatedProgress = progressData()?.find(p => p.id === progress.id);
        if (updatedProgress) {
          setScore(updatedProgress.score);
        }
        setShowResult(true);
      }
    } catch (error) {
      console.error('Ошибка сохранения ответа:', error);
      setError('Не удалось сохранить ответ: ' + (error instanceof Error ? error.message : 'Неизвестная ошибка'));
    }
  }
  
  function goToPreviousQuestion() {
    if (step() > 0) {
      setStep(step() - 1);
      // Восстанавливаем предыдущий ответ
      const prevAnswer = answers()[step() - 1];
      if (typeof prevAnswer === 'string') {
        const quiz = currentQuizFull();
        if (quiz && quiz.questions[step() - 1]) {
          const question = quiz.questions[step() - 1];
          if (question.options && question.options.length > 0) {
            const optionIndex = question.options.findIndex(opt => opt.text === prevAnswer);
            if (optionIndex >= 0) {
              setCurrentChoice(optionIndex);
            }
          } else {
            setCurrentTextAnswer(prevAnswer);
          }
        }
      }
    }
  }

  if (hasAccess() === null || loading()) {
    return <div style={{'text-align': 'center', 'margin-top': '2rem'}}>Загрузка...</div>;
  }
  if (!hasAccess()) {
    return <div style={{'text-align': 'center', 'margin-top': '2rem', color: '#e76f51'}}>Нет доступа. Пожалуйста, войдите в систему.</div>;
  }
  if (error()) {
    return <div style={{'text-align': 'center', 'margin-top': '2rem', color: '#e76f51'}}>{error()}</div>;
  }

  return (
    <div class="quiz-page" style={{'max-width': '1400px', margin: '2rem auto', padding: '0 1rem'}}>
      <h2>Тесты</h2>
      {!currentQuizShort() && (
        <Show when={filteredQuizzes().length > 0} fallback={<div>Нет доступных тестов.</div>}>
          <ul style={{'list-style': 'none', padding: 0, display: 'grid', gap: '3rem', 'grid-template-columns': 'repeat(auto-fit, minmax(300px, 1fr))'}}>
            <For each={filteredQuizzes()}>{quiz => (
              <li style={{'margin-bottom': '0', 'background': 'var(--bg-secondary)', 'border-radius': '8px', padding: '1.2rem', 'box-shadow': '0 2px 8px rgba(0,0,0,0.04)', color: 'var(--text-primary)'}}>
                <h3>{quiz.title}</h3>
                <div style={{color: '#666', 'margin-bottom': '0.5rem'}}><b>Предмет:</b> {quiz.subject?.name || 'Не указан'}</div>
                <div style={{color: '#666', 'margin-bottom': '0.5rem'}}><b>Описание:</b> {quiz.summary}</div>
                <button style={{marginTop: '1em', background: '#2563eb', color: '#fff', border: 'none', borderRadius: '7px', padding: '0.5em 1.2em', fontWeight: 600, cursor: 'pointer'}} onClick={() => startQuiz(quiz)}>Начать тест</button>
              </li>
            )}</For>
          </ul>
        </Show>
      )}
      {loadingQuiz() && (
        <div style={{'text-align': 'center', 'margin-top': '2rem'}}>Загрузка квиза...</div>
      )}
      {currentQuizFull() && !showResult() && !loadingQuiz() && (
        <div style={{'margin-top': '2em', 'max-width': '900px', margin: '2rem auto'}}>
          <h3>{currentQuizFull()?.title}</h3>
          <div style={{'margin-bottom': '1em', color: '#888'}}>
            Вопрос {step() + 1} из {currentQuizFull()?.questions.length || 0}
          </div>
          <div style={{'font-weight': 500, 'margin-bottom': '1.5em', padding: '1em', background: 'var(--bg-secondary)', borderRadius: '8px', color: 'var(--text-primary)'}}>
            {currentQuizFull()?.summary || ''}
          </div>
          
          {currentQuizFull()?.questions[step()] && (
            <div style={{'margin-bottom': '2em'}}>
              <h4 style={{'margin-bottom': '1em', fontSize: '1.2em'}}>
                {currentQuizFull()!.questions[step()].text}
              </h4>
              
              {currentQuizFull()!.questions[step()].options && currentQuizFull()!.questions[step()].options!.length > 0 ? (
                // Вопрос с вариантами ответов
                <div style={{display: 'flex', flexDirection: 'column', gap: '0.8em'}}>
                  <For each={currentQuizFull()!.questions[step()].options}>
                    {(option, idx) => (
                      <label style={{
                        display: 'flex',
                        alignItems: 'center',
                        padding: '1em',
                        border: '2px solid',
                        borderColor: currentChoice() === idx() ? '#2563eb' : '#e3eafc',
                        borderRadius: '8px',
                        cursor: 'pointer',
                        background: currentChoice() === idx() ? 'var(--bg-tertiary)' : 'var(--bg-secondary)',
                        color: 'var(--text-primary)',
                        transition: 'all 0.2s'
                      }}>
                        <input
                          type="radio"
                          name={`question-${step()}`}
                          checked={currentChoice() === idx()}
                          onChange={() => setCurrentChoice(idx())}
                          style={{marginRight: '0.8em', width: '18px', height: '18px', cursor: 'pointer'}}
                        />
                        <span style={{flex: 1}}>{option.text}</span>
                      </label>
                    )}
                  </For>
                </div>
              ) : (
                // Открытый вопрос
                <div>
                  <textarea
                    value={currentTextAnswer()}
                    onInput={(e) => setCurrentTextAnswer(e.currentTarget.value)}
                    placeholder="Введите ваш ответ..."
                    style={{
                      width: '100%',
                      minHeight: '120px',
                      padding: '0.8em',
                      border: '2px solid #e3eafc',
                      borderRadius: '8px',
                      fontSize: '1em',
                      fontFamily: 'inherit',
                      resize: 'vertical'
                    }}
                  />
                </div>
              )}
              
              <div style={{display: 'flex', gap: '1em', marginTop: '2em', justifyContent: 'space-between'}}>
                <button
                  onClick={goToPreviousQuestion}
                  disabled={step() === 0}
                  style={{
                    background: step() === 0 ? '#e3eafc' : '#f3f4f6',
                    color: step() === 0 ? '#9ca3af' : '#213547',
                    border: 'none',
                    borderRadius: '7px',
                    padding: '0.6em 1.5em',
                    fontWeight: 600,
                    cursor: step() === 0 ? 'not-allowed' : 'pointer',
                    opacity: step() === 0 ? 0.6 : 1
                  }}
                >
                  Назад
                </button>
                <button
                  onClick={saveAnswer}
                  style={{
                    background: '#2563eb',
                    color: '#fff',
                    border: 'none',
                    borderRadius: '7px',
                    padding: '0.6em 1.5em',
                    fontWeight: 600,
                    cursor: 'pointer'
                  }}
                >
                  {step() < (currentQuizFull()?.questions.length || 0) - 1 ? 'Следующий вопрос' : 'Завершить тест'}
                </button>
              </div>
            </div>
          )}
        </div>
      )}
      {showResult() && (
        <div style={{marginTop: '2em', textAlign: 'center', maxWidth: '600px', margin: '2rem auto'}}>
          <h3>Результат теста</h3>
          <div style={{fontSize: '1.2em', margin: '1em 0'}}>
            Тест завершен
          </div>
          {currentQuizFull() && (
            <div style={{margin: '1.5em 0', padding: '1em', background: 'var(--bg-secondary)', borderRadius: '8px', color: 'var(--text-primary)'}}>
              <div style={{fontSize: '1.1em', marginBottom: '0.5em'}}>
                <strong>Набрано баллов:</strong> {score()} из {currentQuizFull()!.total_score}
              </div>
              <div style={{color: '#2563eb', fontSize: '0.95em'}}>
                {currentQuizFull()!.total_score > 0 
                  ? `Процент выполнения: ${Math.round((score() / currentQuizFull()!.total_score) * 100)}%`
                  : 'Ожидается проверка преподавателем'}
              </div>
            </div>
          )}
          <div style={{color: '#2a9d8f', fontWeight: 600, marginBottom: '1em'}}>
            Результаты сохранены
          </div>
          <button style={{'margin-top': '1em', background: '#2563eb', color: '#fff', border: 'none', 'border-radius': '7px', padding: '0.5em 1.2em', 'font-weight': 600, cursor: 'pointer'}} onClick={() => { 
            setCurrentQuizShort(null); 
            setCurrentQuizFull(null);
            setShowResult(false); 
            setCurrentProgress(null);
            setStep(0);
            setAnswers([]);
            setCurrentChoice(null);
            setCurrentTextAnswer('');
            setScore(0);
          }}>К списку тестов</button>
        </div>
      )}
    </div>
  );
};

export default Quiz; 