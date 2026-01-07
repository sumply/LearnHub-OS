import { type Component, createSignal, onMount, For, Show, createResource } from 'solid-js';
import { checkToken, getCurrentUser } from '../utils/api';
import { getQuizzes } from '../services/quizService';
import { getProgress, startProgress, finishProgress, updateAnswer } from '../services/progressService';
import * as apiClient from '../utils/apiClient';

const Quiz: Component = () => {
  const [hasAccess, setHasAccess] = createSignal<boolean | null>(null);
  const [error, setError] = createSignal<string | null>(null);
  const [loading, setLoading] = createSignal(true);
  const [currentQuiz, setCurrentQuiz] = createSignal<apiClient.QuizShortResponse | null>(null);
  const [currentProgress, setCurrentProgress] = createSignal<apiClient.QuizProgressResponse | null>(null);
  const [step, setStep] = createSignal(0);
  // Для каждого вопроса: либо индекс выбранного варианта, либо текст ответа
  const [answers, setAnswers] = createSignal<(number | string)[]>([]);
  const [currentChoice, setCurrentChoice] = createSignal<number | null>(null);
  const [showResult, setShowResult] = createSignal(false);
  const [score, setScore] = createSignal(0);

  // Загружаем квизы и прогресс через API
  const [quizzesData] = createResource(getQuizzes);
  const [progressData] = createResource(getProgress);

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
      // Находим прогресс для этого квиза
      const progress = progressData()?.find(p => p.quiz.id === quiz.id);
      if (progress) {
        // Начинаем прохождение
        await startProgress(progress.id);
        setCurrentProgress(progress);
      }
      setCurrentQuiz(quiz);
      setStep(0);
      setAnswers([]);
      setShowResult(false);
      setScore(0);
    } catch (error) {
      console.error('Ошибка начала квиза:', error);
      setError('Не удалось начать квиз');
    }
  }

  async function nextQuestion() {
    const quiz = currentQuiz();
    const progress = currentProgress();
    if (!quiz || !progress) return;
    
    // TODO: Получить полную информацию о квизе с вопросами
    // Пока используем упрощенную логику
    const answerToSave: number | string = currentChoice() !== null ? currentChoice()! : '';
    
    setAnswers([...answers(), answerToSave]);
    setCurrentChoice(null);
    
    // TODO: Получить количество вопросов из полного квиза
    // Пока предполагаем, что есть вопросы
    if (step() < 10) { // Временное значение
      setStep(step() + 1);
    } else {
      // Завершение теста
      try {
        await finishProgress(progress.id);
        setShowResult(true);
        // TODO: Рассчитать score на основе ответов
        setScore(0);
      } catch (error) {
        console.error('Ошибка завершения квиза:', error);
        setError('Не удалось завершить квиз');
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
    <div class="quiz-page" style={{'max-width': '800px', margin: '2rem auto'}}>
      <h2>Тесты</h2>
      {!currentQuiz() && (
        <Show when={filteredQuizzes().length > 0} fallback={<div>Нет доступных тестов.</div>}>
          <ul style={{'list-style': 'none', padding: 0}}>
            <For each={filteredQuizzes()}>{quiz => (
              <li style={{'margin-bottom': '2rem', 'background': '#f8f9fa', 'border-radius': '8px', padding: '1.2rem', 'box-shadow': '0 2px 8px rgba(0,0,0,0.04)'}}>
                <h3>{quiz.title}</h3>
                <div style={{color: '#666', 'margin-bottom': '0.5rem'}}><b>Предмет:</b> {quiz.subject.name}</div>
                <div style={{color: '#666', 'margin-bottom': '0.5rem'}}><b>Описание:</b> {quiz.summary}</div>
                <button style={{marginTop: '1em', background: '#2563eb', color: '#fff', border: 'none', borderRadius: '7px', padding: '0.5em 1.2em', fontWeight: 600, cursor: 'pointer'}} onClick={() => startQuiz(quiz)}>Начать тест</button>
              </li>
            )}</For>
          </ul>
        </Show>
      )}
      {currentQuiz() && !showResult() && (
        <div style={{'margin-top': '2em'}}>
          <h3>{currentQuiz()?.title}</h3>
          <div style={{'margin-bottom': '1em', color: '#888'}}>Прохождение квиза</div>
          <div style={{'font-weight': 500, 'margin-bottom': '1em'}}>
            {currentQuiz()?.summary || 'Начните прохождение квиза'}
          </div>
          <div style={{color: '#666', marginBottom: '1em'}}>
            <p>Для прохождения квиза необходимо получить полную информацию о вопросах.</p>
            <p>Эта функциональность будет реализована после добавления endpoint для получения полного квиза.</p>
          </div>
          <button
            style={{
              marginTop: '1em',
              background: '#2563eb',
              color: '#fff',
              border: 'none',
              'border-radius': '7px',
              padding: '0.5em 1.2em',
              'font-weight': 600,
              cursor: 'pointer',
            }}
            onClick={async () => {
              const progress = currentProgress();
              if (progress) {
                try {
                  await finishProgress(progress.id);
                  setShowResult(true);
                } catch (error) {
                  console.error('Ошибка завершения:', error);
                }
              }
            }}
          >Завершить тест</button>
        </div>
      )}
      {showResult() && (
        <div style={{marginTop: '2em', textAlign: 'center'}}>
          <h3>Результат теста</h3>
          <div style={{fontSize: '1.2em', margin: '1em 0'}}>Тест завершен</div>
          <div style={{color: '#2a9d8f', fontWeight: 600, marginBottom: '1em'}}>
            Результаты сохранены
          </div>
          <button style={{'margin-top': '1em', background: '#2563eb', color: '#fff', border: 'none', 'border-radius': '7px', padding: '0.5em 1.2em', 'font-weight': 600, cursor: 'pointer'}} onClick={() => { setCurrentQuiz(null); setShowResult(false); setCurrentProgress(null); }}>К списку тестов</button>
        </div>
      )}
    </div>
  );
};

export default Quiz; 