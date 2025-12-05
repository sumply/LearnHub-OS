import { type Component, createSignal, onMount, For, Show } from 'solid-js';
import { checkToken } from '../utils/api';
import { quizzes } from '../config/activities';
import type { Quiz as QuizType } from '../config/activities';
import { getCurrentUser } from '../utils/api';
import { getAllUsers } from '../config/users';

const Quiz: Component = () => {
  const [hasAccess, setHasAccess] = createSignal<boolean | null>(null);
  const [quizzes, setQuizzes] = createSignal<QuizType[]>([]);
  const [error, setError] = createSignal<string | null>(null);
  const [loading, setLoading] = createSignal(true);
  const [currentQuiz, setCurrentQuiz] = createSignal<QuizType | null>(null);
  const [step, setStep] = createSignal(0);
  // Для каждого вопроса: либо индекс выбранного варианта, либо текст ответа
  const [answers, setAnswers] = createSignal<(number | string)[]>([]);
  const [currentChoice, setCurrentChoice] = createSignal<number | null>(null);
  const [showResult, setShowResult] = createSignal(false);
  const [score, setScore] = createSignal(0);

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
    try {
      const user = getCurrentUser() as any;
      let filtered = quizzes();
      if (user) {
        if (user.role === 'student') {
          filtered = quizzes().filter((q: any) => q.category === user.group);
        } else if (user.role === 'teacher') {
          filtered = quizzes().filter((q: any) => Array.isArray(user.subjects) && user.subjects.includes(q.category));
        } else if (user.role === 'parent') {
          filtered = quizzes(); // доработать по детям
        }
      }
      setQuizzes(filtered);
    } catch (e: any) {
      setError(e.message || 'Ошибка загрузки тестов');
    } finally {
      setLoading(false);
    }
  });

  function startQuiz(quiz: QuizType) {
    setCurrentQuiz(quiz);
    setStep(0);
    setAnswers([]);
    setShowResult(false);
    setScore(0);
  }

  function nextQuestion() {
    const quiz = currentQuiz() as QuizType;
    if (!quiz) return;
    const q = quiz.questions[step()];
    let answerToSave: number | string = '';
    if (q.type === 'choice') {
      if (currentChoice() === null) return;
      answerToSave = currentChoice()!;
    }
    // для open вопроса ответ уже добавляется через textarea
    setAnswers([...answers(), answerToSave]);
    setCurrentChoice(null);
    if (step() < (quiz.questions.length || 0) - 1) {
      setStep(step() + 1);
    } else {
      // Завершение теста
      let correct = 0;
      let hasOpen = false;
      const allAnswers = [...answers(), answerToSave];
      quiz.questions.forEach((q, i) => {
        if (q.type === 'choice' && allAnswers[i] === q.correct) correct++;
        if (q.type === 'open') hasOpen = true;
      });
      setScore(correct);
      setShowResult(true);
      // Сохраняем результат в activityHistory
      const user = getCurrentUser();
      if (user && user.role === 'student') {
        const all = getAllUsers();
        const idx = all.findIndex((u: any) => u.id === user.id);
        if (idx !== -1) {
          const history = (all[idx] as any).activityHistory || [];
          if (quiz.requiresTeacherCheck || hasOpen) {
            history.push({
              activityId: quiz.id,
              type: 'quiz',
              title: quiz.title,
              date: new Date().toISOString(),
              answers: allAnswers,
              status: 'pending',
            });
          } else {
            history.push({
              activityId: quiz.id,
              type: 'quiz',
              title: quiz.title,
              date: new Date().toISOString(),
              answers: allAnswers,
              score: correct,
              maxScore: quiz.questions.filter(q => q.type === 'choice').length,
              status: correct >= Math.ceil(quiz.questions.filter(q => q.type === 'choice').length * 0.6) ? 'passed' : 'failed',
            });
          }
          all[idx].activityHistory = history;
          localStorage.setItem('users', JSON.stringify(all));
          localStorage.setItem('user', JSON.stringify(all[idx]));
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
    <div class="quiz-page" style={{'max-width': '800px', margin: '2rem auto'}}>
      <h2>Тесты</h2>
      {!currentQuiz() && (
        <Show when={quizzes().length > 0} fallback={<div>Нет доступных тестов.</div>}>
          <ul style={{'list-style': 'none', padding: 0}}>
            <For each={quizzes()}>{quiz => (
              <li style={{'margin-bottom': '2rem', 'background': '#f8f9fa', 'border-radius': '8px', padding: '1.2rem', 'box-shadow': '0 2px 8px rgba(0,0,0,0.04)'}}>
                <h3>{quiz.title}</h3>
                <div style={{color: '#666', 'margin-bottom': '0.5rem'}}><b>Категория:</b> {quiz.category}</div>
                <div><b>Вопросов:</b> {quiz.questions.length}</div>
                <button style={{marginTop: '1em', background: '#2563eb', color: '#fff', border: 'none', borderRadius: '7px', padding: '0.5em 1.2em', fontWeight: 600, cursor: 'pointer'}} onClick={() => startQuiz(quiz)}>Начать тест</button>
              </li>
            )}</For>
          </ul>
        </Show>
      )}
      {currentQuiz() && !showResult() && (
        <div style={{'margin-top': '2em'}}>
          <h3>{currentQuiz()?.title}</h3>
          <div style={{'margin-bottom': '1em', color: '#888'}}>Вопрос {step() + 1} из {currentQuiz()?.questions.length}</div>
          <div style={{'font-weight': 500, 'margin-bottom': '1em'}}>{currentQuiz()?.questions[step()].question}</div>
          <Show
            when={currentQuiz()?.questions[step()].type === 'choice'}
            fallback={
              <div style={{'margin-bottom': '1em'}}>
                <textarea
                  style={{width: '100%', 'min-height': '80px', 'border-radius': '7px', border: '1px solid #ccc', padding: '0.5em', 'font-size': '1em'}}
                  placeholder="Введите ваш ответ..."
                  onInput={e => (e.target as HTMLTextAreaElement).value}
                  id="open-answer"
                />
                <button
                  style={{'margin-top': '1em', background: '#2563eb', color: '#fff', border: 'none', 'border-radius': '7px', padding: '0.5em 1.2em', 'font-weight': 600, cursor: 'pointer'}}
                  onClick={() => {
                    const val = (document.getElementById('open-answer') as HTMLTextAreaElement)?.value || '';
                    setAnswers([...answers(), val]);
                    setStep(step() + 1);
                  }}
                >Ответить</button>
              </div>
            }
          >
            <ul style={{'list-style': 'none', padding: 0}}>
              {currentQuiz()?.questions[step()].options?.map((opt, idx) => (
                <li style={{'margin-bottom': '0.7em'}}>
                  <button
                    style={{
                      background: currentChoice() === idx ? '#2563eb' : '#e3eafc',
                      color: currentChoice() === idx ? '#fff' : '#2563eb',
                      border: 'none',
                      'border-radius': '7px',
                      padding: '0.5em 1.2em',
                      'font-weight': 600,
                      cursor: 'pointer',
                      width: '100%',
                      'text-align': 'left',
                      outline: currentChoice() === idx ? '2px solid #2563eb' : 'none',
                      transition: 'background 0.2s, color 0.2s',
                    }}
                    onClick={() => setCurrentChoice(idx)}
                  >{opt}</button>
                </li>
              ))}
            </ul>
            <button
              style={{
                marginTop: '1em',
                background: currentChoice() !== null ? '#2563eb' : '#e3eafc',
                color: currentChoice() !== null ? '#fff' : '#888',
                border: 'none',
                'border-radius': '7px',
                padding: '0.5em 1.2em',
                'font-weight': 600,
                cursor: currentChoice() !== null ? 'pointer' : 'not-allowed',
                width: '100%',
                fontSize: '1.1em',
                transition: 'background 0.2s, color 0.2s',
              }}
              disabled={currentChoice() === null}
              onClick={nextQuestion}
            >{step() === (currentQuiz()?.questions.length || 0) - 1 ? 'Завершить тест' : 'Далее'}</button>
          </Show>
        </div>
      )}
      {showResult() && (
        <div style={{marginTop: '2em', textAlign: 'center'}}>
          <h3>Результат теста</h3>
          <div style={{fontSize: '1.2em', margin: '1em 0'}}>Правильных ответов: {score()} из {currentQuiz()?.questions.length}</div>
          <div style={{color: score() >= Math.ceil((currentQuiz()?.questions.length || 0) * 0.6) ? '#2a9d8f' : '#e76f51', fontWeight: 600, marginBottom: '1em'}}>
            {score() >= Math.ceil((currentQuiz()?.questions.length || 0) * 0.6) ? 'Тест пройден!' : 'Тест не пройден'}
          </div>
          <button style={{'margin-top': '1em', background: '#2563eb', color: '#fff', border: 'none', 'border-radius': '7px', padding: '0.5em 1.2em', 'font-weight': 600, cursor: 'pointer'}} onClick={() => { setCurrentQuiz(null); setShowResult(false); }}>К списку тестов</button>
        </div>
      )}
    </div>
  );
};

export default Quiz; 