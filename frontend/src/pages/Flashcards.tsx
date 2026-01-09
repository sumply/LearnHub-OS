import { type Component, createSignal, onMount, For, Show } from 'solid-js';
import { checkToken } from '../utils/api';
import { flashcards } from '../config/activities';
import type { Flashcard } from '../config/activities';
import { getCurrentUser } from '../utils/api';
import { getAllUsers } from '../config/users';
import Header from '../components/Header';

const Flashcards: Component = () => {
  const [hasAccess, setHasAccess] = createSignal<boolean | null>(null);
  const [flashcards, setFlashcards] = createSignal<Flashcard[]>([]);
  const [error, setError] = createSignal<string | null>(null);
  const [loading, setLoading] = createSignal(true);
  const [currentIdx, setCurrentIdx] = createSignal(0);
  const [showAnswer, setShowAnswer] = createSignal(false);
  const [session, setSession] = createSignal<{ correct: number; total: number } | null>(null);

  function startSession() {
    setCurrentIdx(0);
    setShowAnswer(false);
    setSession(null);
  }

  function handleShowAnswer() {
    setShowAnswer(true);
  }

  function handleNextCard(isCorrect: boolean) {
    const next = currentIdx() + 1;
    if (!session()) setSession({ correct: isCorrect ? 1 : 0, total: 1 });
    else setSession({ correct: session()!.correct + (isCorrect ? 1 : 0), total: session()!.total + 1 });
    setShowAnswer(false);
    if (next < flashcards().length) {
      setCurrentIdx(next);
    } else {
      // Сессия завершена — сохранить результат
      const user = getCurrentUser();
      if (user && user.role === 'student') {
        const all = getAllUsers();
        const idx = all.findIndex(u => u.id === user.id);
        if (idx !== -1) {
          const history = all[idx].activityHistory || [];
          history.push({
            activityId: 'flashcards_' + Date.now(),
            type: 'flashcard',
            title: 'Сессия карточек (' + (new Date().toLocaleDateString('ru-RU')) + ')',
            date: new Date().toISOString(),
            score: session() ? session()!.correct + (isCorrect ? 1 : 0) : (isCorrect ? 1 : 0),
            maxScore: session() ? session()!.total + 1 : 1,
            status: ((session() ? session()!.correct : 0) + (isCorrect ? 1 : 0)) >= Math.ceil((session() ? session()!.total : 0) + 1 * 0.6) ? 'passed' : 'failed',
          });
          all[idx].activityHistory = history;
          localStorage.setItem('users', JSON.stringify(all));
          localStorage.setItem('user', JSON.stringify(all[idx]));
        }
      }
      setSession(null);
      setCurrentIdx(0);
    }
  }

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
      let filtered = flashcards;
      if (user) {
        if (user.role === 'student') {
          filtered = flashcards.filter((f: any) => f.category === user.group);
        } else if (user.role === 'teacher') {
          filtered = flashcards.filter((f: any) => Array.isArray(user.subjects) && user.subjects.includes(f.category));
        } else if (user.role === 'parent') {
          filtered = flashcards; // доработать по детям
        }
      }
      setFlashcards(filtered);
    } catch (e: any) {
      setError(e.message || 'Ошибка загрузки карточек');
    } finally {
      setLoading(false);
    }
  });

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
    <>
      <Header />
      <div class="main-shell">
        <h2>Карточки</h2>
        {flashcards().length > 0 && session() === null && (
          <button style={{marginBottom: '1.5em', background: '#2563eb', color: '#fff', border: 'none', borderRadius: '7px', padding: '0.5em 1.2em', fontWeight: 600, cursor: 'pointer'}} onClick={startSession}>Начать сессию</button>
        )}
        {session() !== null && (
          <div style={{marginTop: '2em', textAlign: 'center'}}>
            <h3>Сессия карточек</h3>
            <div style={{fontSize: '1.1em', margin: '1em 0'}}>Карточка {currentIdx() + 1} из {flashcards().length}</div>
            <div style={{fontWeight: 500, marginBottom: '1em'}}>{flashcards()[currentIdx()].question}</div>
            {!showAnswer() ? (
              <button style={{background: '#e3eafc', color: '#2563eb', border: 'none', borderRadius: '7px', padding: '0.5em 1.2em', fontWeight: 600, cursor: 'pointer'}} onClick={handleShowAnswer}>Показать ответ</button>
            ) : (
              <>
                <div style={{margin: '1em 0', color: '#213547', fontWeight: 600}}>Ответ: {flashcards()[currentIdx()].answer}</div>
                <div style={{display: 'flex', justifyContent: 'center', gap: '1.5em'}}>
                  <button style={{background: '#2a9d8f', color: '#fff', border: 'none', borderRadius: '7px', padding: '0.5em 1.2em', fontWeight: 600, cursor: 'pointer'}} onClick={() => handleNextCard(true)}>Я знал</button>
                  <button style={{background: '#e76f51', color: '#fff', border: 'none', borderRadius: '7px', padding: '0.5em 1.2em', fontWeight: 600, cursor: 'pointer'}} onClick={() => handleNextCard(false)}>Не знал</button>
                </div>
              </>
            )}
          </div>
        )}
        <Show when={flashcards().length > 0} fallback={<div>Нет доступных карточек.</div>}>
          <ul style={{'list-style': 'none', padding: 0}}>
            <For each={flashcards()}>{fc => (
              <li style={{'margin-bottom': '1.5rem', 'background': '#f8f9fa', 'border-radius': '8px', padding: '1.2rem', 'box-shadow': '0 2px 8px rgba(0,0,0,0.04)'}}>
                <div><b>Вопрос:</b> {fc.question}</div>
                <div><b>Ответ:</b> {fc.answer}</div>
                <div style={{color: '#888'}}><b>Категория:</b> {fc.category}</div>
              </li>
            )}</For>
          </ul>
        </Show>
      </div>
    </>
  );
};

export default Flashcards; 