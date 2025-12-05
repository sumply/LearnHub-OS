import { type Component, createSignal, Show, For } from 'solid-js';
import { A } from '@solidjs/router';
import { getCurrentUser } from '../utils/api';
import { getAllUsers } from '../config/users';
import { getAllActivities } from '../utils/activitiesService';
import { quizzes } from '../config/activities';

const TeacherDashboard: Component = () => {
  const user = getCurrentUser();
  const isTeacher = user?.role === 'teacher';
  const allUsers = getAllUsers();
  const allActivities = getAllActivities();

  // Классы, где учитель ведёт уроки (по своим предметам)
  const mySubjects = user?.subjects || [];
  const myClasses = Array.from(new Set(
    allUsers.filter(u => u.role === 'student' && mySubjects.some((subj: string) => {
      // ищем задания по этому предмету
      return allActivities.some(a => a.teacher === user.name && a.category === subj && a.type === 'quiz');
    })).map(u => u.group)
  ));
  const [selectedClass, setSelectedClass] = createSignal<string | null>(null);

  // Ученики выбранного класса
  const studentsInClass = () =>
    allUsers.filter(u => u.role === 'student' && u.group === selectedClass());

  // История по ученику только по заданиям этого учителя
  function getStudentHistory(student: any) {
    return (student.activityHistory || []).filter((a: any) => {
      const act = allActivities.find(act => act.id === a.activityId);
      return act && act.teacher === user.name;
    });
  }

  // Проверка заданий на проверке
  function handleCheck(activity: any, student: any, status: string, score: number, maxScore: number, comment: string) {
    const all = getAllUsers();
    const idx = all.findIndex(u => u.id === student.id);
    if (idx !== -1) {
      const hist = all[idx].activityHistory || [];
      const actIdx = hist.findIndex(a => a.activityId === activity.activityId && a.date === activity.date);
      if (actIdx !== -1) {
        hist[actIdx].status = status;
        hist[actIdx].score = score;
        hist[actIdx].maxScore = maxScore;
        hist[actIdx].comment = comment;
        all[idx].activityHistory = hist;
        localStorage.setItem('users', JSON.stringify(all));
        window.location.reload();
      }
    }
  }

  return (
    <div class="dashboard">
      <h2>Добро пожаловать, учитель!</h2>
      <ul>
        <li><A href="/teacher/groups">Мои классы</A></li>
        <li><A href="/teacher/materials">Мои материалы</A></li>
        <li><A href="/tasks">Интерактивные задания</A></li>
        <li><A href="/quiz">Тесты</A></li>
      </ul>
      {isTeacher && (
        <div style={{'margin-top': '2rem'}}>
          <button style={{padding: '10px 20px', background: '#2a9d8f', color: '#fff', border: 'none', 'border-radius': '5px', cursor: 'pointer', 'margin-right': '1rem'}}>Добавить материал</button>
          <button style={{padding: '10px 20px', background: '#e76f51', color: '#fff', border: 'none', 'border-radius': '5px', cursor: 'pointer', 'margin-right': '1rem'}}>Создать тест</button>
          <button style={{padding: '10px 20px', background: '#f4a261', color: '#fff', border: 'none', 'border-radius': '5px', cursor: 'pointer', 'margin-right': '1rem'}}>Поставить оценку</button>
          <button style={{padding: '10px 20px', background: '#264653', color: '#fff', border: 'none', 'border-radius': '5px', cursor: 'pointer'}}>Написать комментарий родителю/ученику</button>
        </div>
      )}
      {isTeacher && (
        <div style={{ marginTop: '2.5rem' }}>
          <h3 style={{ color: '#2563eb', marginBottom: '1.2em', fontSize: '1.18em' }}>Успеваемость моих классов</h3>
          <div style={{ marginBottom: '1.5em' }}>
            <b>Мои классы:</b>
            {myClasses.length === 0 && <span style={{ color: '#888', marginLeft: '1em' }}>Нет классов</span>}
            <For each={myClasses}>{group => (
              <button style={{ marginLeft: '1em', marginBottom: '0.5em', padding: '0.5em 1.2em', borderRadius: '8px', border: '1.5px solid #2563eb', background: selectedClass() === group ? '#2563eb' : '#fff', color: selectedClass() === group ? '#fff' : '#2563eb', cursor: 'pointer', fontWeight: 600 }} onClick={() => setSelectedClass(group)}>{group}</button>
            )}</For>
          </div>
          <Show when={selectedClass()}>
            <div style={{ marginTop: '1.5em' }}>
              <h4 style={{ color: '#2563eb', marginBottom: '0.7em' }}>Ученики класса {selectedClass()}</h4>
              <table style={{ width: '100%', borderCollapse: 'collapse', background: 'var(--bg-primary)', borderRadius: '8px', overflow: 'hidden', boxShadow: '0 2px 8px rgba(37,99,235,0.10)' }}>
                <thead>
                  <tr style={{ background: '#e3eafc', color: '#213547' }}>
                    <th style={{ padding: '0.6em', textAlign: 'left' }}>ФИО</th>
                    <th style={{ padding: '0.6em', textAlign: 'left' }}>История заданий</th>
                  </tr>
                </thead>
                <tbody>
                  <For each={studentsInClass()}>{student => (
                    <tr>
                      <td style={{ padding: '0.5em', fontWeight: 600 }}>{student.surname} {student.name}</td>
                      <td style={{ padding: '0.5em' }}>
                        <For each={getStudentHistory(student)}>{activity => (
                          <div style={{ marginBottom: '0.7em', background: '#f7fafd', borderRadius: '8px', padding: '0.5em 1em', border: '1px solid #e3eafc' }}>
                            <b>{activity.title}</b> — <span style={{ color: activity.status === 'passed' ? '#2a9d8f' : activity.status === 'failed' ? '#e76f51' : '#888' }}>{activity.status === 'passed' ? 'Пройдено' : activity.status === 'failed' ? 'Не пройдено' : 'На проверке'}</span>
                            {activity.status === 'pending' && (
                              <div style={{ marginTop: '0.7em', background: '#f8f9fa', borderRadius: '8px', padding: '0.7em 1em', border: '1px solid #e3eafc' }}>
                                {/* Детализация по вопросам */}
                                {activity.type === 'quiz' && Array.isArray(activity.answers) && (() => {
                                  const quiz = quizzes.find(q => q.id === activity.activityId);
                                  if (!quiz) return <div style={{ color: '#e76f51', marginTop: '1em' }}>Тест не найден</div>;
                                  return (
                                    <div style={{ marginBottom: '1em' }}>
                                      <b>Ответы ученика:</b>
                                      <ul style={{ paddingLeft: '1.2em', marginTop: '0.5em' }}>
                                        {quiz.questions.map((q, qidx) => (
                                          <li style={{ marginBottom: '0.7em' }}>
                                            <div><b>Вопрос {qidx + 1}:</b> {q.question}</div>
                                            <div style={{ marginLeft: '0.7em' }}>
                                              <b>Ответ ученика:</b> {typeof activity.answers[qidx] === 'string' ? activity.answers[qidx] : (q.options && typeof activity.answers[qidx] === 'number' ? q.options[activity.answers[qidx] as number] : '')}
                                            </div>
                                            {q.type === 'choice' && q.correct !== undefined && (
                                              <div style={{ marginLeft: '0.7em' }}>
                                                <b>Правильный ответ:</b> {q.options ? q.options[q.correct] : ''}
                                              </div>
                                            )}
                                          </li>
                                        ))}
                                      </ul>
                                    </div>
                                  );
                                })()}
                                {/* Форма проверки */}
                                <form onSubmit={e => {
                                  e.preventDefault();
                                  const form = e.target as HTMLFormElement;
                                  const status = (form.elements.namedItem('status') as HTMLSelectElement).value;
                                  const score = +(form.elements.namedItem('score') as HTMLInputElement).value;
                                  const maxScore = +(form.elements.namedItem('maxScore') as HTMLInputElement).value;
                                  const comment = (form.elements.namedItem('comment') as HTMLInputElement).value;
                                  handleCheck(activity, student, status, score, maxScore, comment);
                                }} style={{ display: 'flex', gap: '0.5em', alignItems: 'center', marginTop: '1em', flexWrap: 'wrap' }}>
                                  <select name="status" defaultValue="passed" style={{ fontSize: '1em' }}>
                                    <option value="passed">Пройдено</option>
                                    <option value="failed">Не пройдено</option>
                                  </select>
                                  <input name="score" type="number" min="0" placeholder="Баллы" style={{ width: '60px', fontSize: '1em' }} />
                                  <input name="maxScore" type="number" min="0" placeholder="Макс." style={{ width: '60px', fontSize: '1em' }} />
                                  <input name="comment" type="text" placeholder="Комментарий" style={{ minWidth: '180px', fontSize: '1em' }} />
                                  <button type="submit" style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: '8px', padding: '0.3em 1em', fontWeight: 600, fontSize: '1em', cursor: 'pointer' }}>Сохранить</button>
                                </form>
                              </div>
                            )}
                            {activity.files && activity.files.length > 0 && (
                              <div style={{ marginTop: '0.3em', color: '#2563eb', fontSize: '0.95em' }}>
                                <b>Файлы:</b> {activity.files.map(f => <span style={{ marginRight: '0.7em' }}>{f}</span>)}
                              </div>
                            )}
                            {activity.status !== 'in_progress' && (activity.score !== undefined) && (
                              <span style={{ marginLeft: '1em', color: '#213547' }}>Оценка: {activity.score}{activity.maxScore ? ` / ${activity.maxScore}` : ''}</span>
                            )}
                          </div>
                        )}</For>
                      </td>
                    </tr>
                  )}</For>
                </tbody>
              </table>
            </div>
          </Show>
        </div>
      )}
    </div>
  );
};

export default TeacherDashboard; 