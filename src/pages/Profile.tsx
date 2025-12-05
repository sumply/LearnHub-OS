import { getCurrentUser } from '../utils/api';
import { useParams } from '@solidjs/router';
import Header from '../components/Header';
import { Show, For, createSignal, createEffect } from 'solid-js';
import { getAllUsers } from '../config/users';
import { getAllActivities } from '../utils/activitiesService';
import { quizzes } from '../config/activities';
// Импортировать типы
import type { Student, User } from '../config/users';
// Импортируем список всех групп и предметов
import { groups, subjects } from '../config/subjectsGroups';

function TeacherClassPerformance(props: { user: User }) {
  const user = props.user as Student | User;
  const allUsers = getAllUsers();
  const allActivities = getAllActivities();
  const mySubjects = Array.isArray((user as any).subjects) ? (user as any).subjects : [];
  const myClasses = Array.from(new Set(
    allUsers
      .filter((u) => u.role === 'student' && (u as Student).activityHistory && (u as Student).activityHistory!.length > 0)
      .filter((student) => (student as Student).activityHistory!.some((activity: any) => {
        const act = allActivities.find((a) => a.id === activity.activityId);
        return act && mySubjects.includes(act.category);
      }))
      .map((u) => (u as any).group)
  ));
  const [selectedClass, setSelectedClass] = createSignal<string|null>(null);
  const [commentEdit, setCommentEdit] = createSignal<any>(null);
  function saveComment(student: Student, activity: any, value: string) {
    const all = getAllUsers();
    const idx = all.findIndex((u) => u.id === student.id);
    if (idx !== -1) {
      const hist = (all[idx] as Student).activityHistory || [];
      const actIdx = hist.findIndex((a: any) => a.activityId === activity.activityId && a.date === activity.date);
      if (actIdx !== -1) {
        hist[actIdx].comment = value;
        (all[idx] as Student).activityHistory = hist;
    localStorage.setItem('users', JSON.stringify(all));
        setCommentEdit(null);
    window.location.reload();
  }
    }
  }
  function studentsInClass() {
    return allUsers.filter((u) => u.role === 'student' && (u as any).group === selectedClass());
  }
  function getStudentHistory(student: Student) {
    return ((student as Student).activityHistory || []).filter((a: any) => {
      const act = allActivities.find((act) => act.id === a.activityId);
      return act && mySubjects.includes(act.category);
    });
  }
  return (
    <div class="profile-gradient-block" style={{ 'margin-top': '2.5rem' }}>
      <h3 style={{ color: '#2563eb', 'margin-bottom': '1.2em', 'font-size': '1.18em' }}>Успеваемость моих классов</h3>
      <div style={{ 'margin-bottom': '1.5em' }}>
        <b>Мои классы:</b>
        {myClasses.length === 0 && <span style={{ color: '#888', 'margin-left': '1em' }}>Нет классов</span>}
        <For each={myClasses}>{(group) => (
          <button style={{ 'margin-left': '1em', 'margin-bottom': '0.5em', padding: '0.5em 1.2em', 'border-radius': '8px', border: '1.5px solid #2563eb', background: selectedClass() === group ? '#2563eb' : '#fff', color: selectedClass() === group ? '#fff' : '#2563eb', cursor: 'pointer', 'font-weight': 600 }} onClick={() => setSelectedClass(group)}>{group}</button>
        )}</For>
      </div>
      <Show when={selectedClass()}>
        <div style={{ 'margin-top': '1.5em' }}>
          <h4 style={{ color: '#2563eb', 'margin-bottom': '0.7em' }}>Ученики класса {selectedClass()}</h4>
          <table style={{ width: '100%', 'border-collapse': 'collapse', background: 'var(--bg-primary)', 'border-radius': '8px', overflow: 'hidden', 'box-shadow': '0 2px 8px rgba(37,99,235,0.10)' }}>
            <thead>
              <tr style={{ background: '#e3eafc', color: '#213547' }}>
                <th style={{ padding: '0.6em', 'text-align': 'left' }}>ФИО</th>
                <th style={{ padding: '0.6em', 'text-align': 'left' }}>История заданий</th>
              </tr>
            </thead>
            <tbody>
              <For each={studentsInClass()}>{(student) => (
                <tr>
                  <td style={{ padding: '0.5em', 'font-weight': 600 }}>{student.surname} {student.name}</td>
                  <td style={{ padding: '0.5em' }}>
                    <For each={getStudentHistory(student)}>{(activity) => (
                      <div style={{ 'margin-bottom': '0.7em', background: '#f7fafd', 'border-radius': '8px', padding: '0.5em 1em', border: '1px solid #e3eafc' }}>
                        <b>{activity.title}</b> — <span style={{ color: activity.status === 'passed' ? '#2a9d8f' : activity.status === 'failed' ? '#e76f51' : '#888' }}>{activity.status === 'passed' ? 'Пройдено' : activity.status === 'failed' ? 'Не пройдено' : 'На проверке'}</span>
                        {activity.score !== undefined && (
                          <span style={{ 'margin-left': '1em', color: '#213547' }}>Оценка: {activity.score}{activity.maxScore ? ` / ${activity.maxScore}` : ''}</span>
                        )}
                        <button style={{ 'margin-left': '1em', background: '#2563eb', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.2em 0.8em', cursor: 'pointer', fontWeight: 500 }} onClick={() => setCommentEdit({ studentId: student.id, activityId: activity.activityId, date: activity.date, value: activity.comment || '' })}>Комментарий</button>
                        <Show when={commentEdit() && commentEdit().studentId === student.id && commentEdit().activityId === activity.activityId && commentEdit().date === activity.date}>
                          <form onSubmit={e => { e.preventDefault(); saveComment(student, activity, commentEdit().value); }} style={{ marginTop: '0.5em', display: 'flex', gap: '0.5em' }}>
                            <input type="text" value={commentEdit().value} onInput={e => setCommentEdit({ ...commentEdit(), value: e.target.value })} placeholder="Комментарий учителя" style={{ flex: 1 }} />
                            <button type="submit" style={{ background: '#2a9d8f', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.2em 0.8em', cursor: 'pointer', fontWeight: 500 }}>Сохранить</button>
                            <button type="button" onClick={() => setCommentEdit(null)} style={{ background: '#e76f51', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.2em 0.8em', cursor: 'pointer', fontWeight: 500 }}>Отмена</button>
                          </form>
                        </Show>
                        <Show when={(activity as any).comment}>
                          <div style={{ marginTop: '0.3em', color: '#213547', fontSize: '0.95em' }}><b>Комментарий учителя:</b> {(activity as any).comment}</div>
                        </Show>
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
  );
}

const userIcons = {
  surname: '👤',
  name: '📝',
  email: '✉️',
  role: '🎓',
};

// Компонент аватарки
function AvatarUploader({ user, onChange }: { user: any, onChange: (dataUrl: string) => void }) {
  const [preview, setPreview] = createSignal(user.avatar || '');
  const [hover, setHover] = createSignal(false);
  createEffect(() => { setPreview(user.avatar || ''); });
  function handleFile(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (ev) => {
        setPreview(ev.target?.result as string);
        onChange(ev.target?.result as string);
      };
      reader.readAsDataURL(file);
    }
  }
  return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', marginBottom: '1.5em' }}>
      <div
        class="usp"
        style={{
          width: '120px', height: '120px', borderRadius: '50%', overflow: 'hidden', background: 'var(--bg-secondary, #e3eafc)', marginBottom: '0.7em', border: '3px solid #2563eb', boxShadow: '0 4px 16px rgba(37,99,235,0.18)', position: 'relative', display: 'flex', alignItems: 'center', justifyContent: 'center', transition: 'box-shadow 0.2s', cursor: 'pointer',
        }}
        onMouseEnter={() => setHover(true)}
        onMouseLeave={() => setHover(false)}
      >
        <Show when={preview()} fallback={<span style={{ color: '#2563eb', fontSize: '3em' }}>👤</span>}>
          <img src={preview()} alt="avatar" style={{ width: '100%', height: '100%', objectFit: 'cover', borderRadius: '50%' }} />
        </Show>
        <Show when={hover()}>
          <div style={{ position: 'absolute', top: 0, left: 0, width: '100%', height: '100%', background: 'rgba(37,99,235,0.18)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', fontWeight: 700, fontSize: '1.2em', transition: 'background 0.2s', borderRadius: '50%' }}>
            <span style={{ fontSize: '2em', marginRight: '0.3em' }}>📷</span> Загрузить
          </div>
        </Show>
        <input type="file" accept="image/*" style={{ display: 'none' }} id="avatar-upload" onChange={handleFile} />
        <label for="avatar-upload" style={{ position: 'absolute', width: '100%', height: '100%', top: 0, left: 0, cursor: 'pointer' }}></label>
      </div>
    </div>
  );
}

const Profile = () => {
  const params = useParams();
  const allUsers = getAllUsers();
  const currentUser = getCurrentUser();
  let user = currentUser;
  let isAdminView = false;
  if (params.id) {
    const found = allUsers.find(u => u.id === params.id);
    if (found) {
      user = found;
      isAdminView = currentUser?.role === 'admin';
    }
  }
  if (!user) return <div style={{ color: '#e76f51', margin: '2em' }}>Пользователь не найден</div>;

  // Состояния для редактирования
  const [editField, setEditField] = createSignal<string | null>(null);
  const [editValue, setEditValue] = createSignal<string>('');
  function startEdit(field: string, value: string) {
    setEditField(field);
    setEditValue(value);
  }
  function saveEdit(field: string) {
    const all = getAllUsers();
    const idx = all.findIndex((u: any) => u.id === user.id);
    if (idx !== -1) {
      all[idx][field] = editValue();
    localStorage.setItem('users', JSON.stringify(all));
      if (user.id === currentUser?.id) {
        localStorage.setItem('user', JSON.stringify(all[idx]));
      }
    }
    user[field] = editValue();
    setEditField(null);
  }

  function saveAvatar(dataUrl: string) {
    const all = getAllUsers();
    const idx = all.findIndex((u: any) => u.id === user.id);
    if (idx !== -1) {
      all[idx].avatar = dataUrl;
        localStorage.setItem('users', JSON.stringify(all));
      localStorage.setItem('user', JSON.stringify(all[idx]));
    }
    user.avatar = dataUrl;
  }
  const icons = { success: '✅', fail: '❌', progress: '⏳' };
  const [selectedActivityIdx, setSelectedActivityIdx] = createSignal<number|null>(null);
  // Для редактирования групп и предметов
  const [editGroups, setEditGroups] = createSignal(false);
  const [editSubjects, setEditSubjects] = createSignal(false);
  const [selectedGroup, setSelectedGroup] = createSignal('');
  const [selectedSubject, setSelectedSubject] = createSignal('');
  const [showDeleteModal, setShowDeleteModal] = createSignal(false);

  function addGroupToTeacher(groupId: string) {
    const all = getAllUsers();
    const idx = all.findIndex((u: any) => u.id === user.id);
    if (idx !== -1) {
      if (!Array.isArray(all[idx].groups)) all[idx].groups = [];
      if (!all[idx].groups.includes(groupId)) all[idx].groups.push(groupId);
      localStorage.setItem('users', JSON.stringify(all));
      if (user.id === currentUser?.id) localStorage.setItem('user', JSON.stringify(all[idx]));
    }
    user.groups = all[idx].groups;
    setSelectedGroup('');
  }
  function removeGroupFromTeacher(groupId: string) {
    const all = getAllUsers();
    const idx = all.findIndex((u: any) => u.id === user.id);
    if (idx !== -1 && Array.isArray(all[idx].groups)) {
      all[idx].groups = all[idx].groups.filter((g: string) => g !== groupId);
      localStorage.setItem('users', JSON.stringify(all));
      if (user.id === currentUser?.id) localStorage.setItem('user', JSON.stringify(all[idx]));
    }
    user.groups = all[idx].groups;
  }
  function addSubjectToTeacher(subjectId: string) {
    const all = getAllUsers();
    const idx = all.findIndex((u: any) => u.id === user.id);
    if (idx !== -1) {
      if (!Array.isArray(all[idx].subjects)) all[idx].subjects = [];
      if (!all[idx].subjects.includes(subjectId)) all[idx].subjects.push(subjectId);
      localStorage.setItem('users', JSON.stringify(all));
      if (user.id === currentUser?.id) localStorage.setItem('user', JSON.stringify(all[idx]));
    }
    user.subjects = all[idx].subjects;
    setSelectedSubject('');
  }
  function removeSubjectFromTeacher(subjectId: string) {
    const all = getAllUsers();
    const idx = all.findIndex((u: any) => u.id === user.id);
    if (idx !== -1 && Array.isArray(all[idx].subjects)) {
      all[idx].subjects = all[idx].subjects.filter((s: string) => s !== subjectId);
      localStorage.setItem('users', JSON.stringify(all));
      if (user.id === currentUser?.id) localStorage.setItem('user', JSON.stringify(all[idx]));
    }
    user.subjects = all[idx].subjects;
  }
  function deleteUser() {
    const all = getAllUsers();
    const idx = all.findIndex((u: any) => u.id === user.id);
    if (idx !== -1) {
      all.splice(idx, 1);
      localStorage.setItem('users', JSON.stringify(all));
      if (user.id === currentUser?.id) localStorage.removeItem('user');
    }
    window.location.href = '/admin';
  }

  if (user.role === 'parent') {
    // Получаем детей
    const children = allUsers.filter(u => user.childrenIds.includes(u.id));
    // Возвращаем профиль родителя с блоком успеваемости детей
    return (
      <>
        <Header />
        <div class="main-shell">
          <div class="profile-gradient-block usp" style={{ maxWidth: '420px', margin: '0 auto', boxShadow: '0 4px 24px rgba(37,99,235,0.13)', padding: '2.2em 2em 2em 2em', background: 'var(--bg-primary, #181a1b)', color: 'var(--text-primary, #e3eafc)', border: '1.5px solid #2563eb', transition: 'background 0.2s, color 0.2s' }}>
            <h2 style={{ color: '#2563eb', fontWeight: 700, fontSize: '1.25rem', letterSpacing: '0.01em', marginBottom: '1.2em', textAlign: 'center' }}>Профиль родителя</h2>
            <div style={{ marginBottom: '2em', color: '#2563eb', fontWeight: 600 }}>Привязанные дети:</div>
            <For each={children}>{child => (
              <div style={{ marginBottom: '2em', background: 'rgba(37,99,235,0.07)', borderRadius: '14px', padding: '1.2em 1em', boxShadow: '0 2px 12px rgba(37,99,235,0.07)' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1em', marginBottom: '0.7em' }}>
                  <span style={{ fontWeight: 700, color: '#2563eb' }}>{child.surname} {child.name}</span>
                  <button style={{ background: 'none', color: '#2563eb', border: '1px solid #2563eb', borderRadius: '3px', padding: '0.2em 0.7em', cursor: 'pointer', fontWeight: 400, fontSize: '0.95em' }}>
                    <a href={`/admin/profile/${child.id}`} style={{ color: '#2563eb', textDecoration: 'none' }}>К данным</a>
                  </button>
                </div>
                <Show when={child.activityHistory && child.activityHistory.length > 0} fallback={<div style={{ color: '#888', fontSize: '1.1em' }}>Нет данных по успеваемости</div>}>
                  <table style={{ width: '100%', borderCollapse: 'collapse', background: 'var(--bg-primary)', borderRadius: '8px', overflow: 'hidden', boxShadow: '0 2px 8px rgba(37,99,235,0.10)' }}>
                    <thead>
                      <tr style={{ background: '#e3eafc', color: '#213547' }}>
                        <th style={{ padding: '0.6em', textAlign: 'left' }}>Задание</th>
                        <th style={{ padding: '0.6em', textAlign: 'left' }}>Дата</th>
                        <th style={{ padding: '0.6em', textAlign: 'left' }}>Результат</th>
                        <th style={{ padding: '0.6em', textAlign: 'left' }}>Баллы</th>
                      </tr>
                    </thead>
                    <tbody>
                      <For each={child.activityHistory}>{item => (
                        <tr>
                          <td style={{ padding: '0.5em' }}>{item.title}</td>
                          <td style={{ padding: '0.5em' }}>{new Date(item.date).toLocaleString('ru-RU')}</td>
                          <td style={{ padding: '0.5em', color: item.status === 'passed' ? '#2a9d8f' : item.status === 'failed' ? '#e76f51' : '#888' }}>
                            {item.status === 'passed' ? 'Пройдено' : item.status === 'failed' ? 'Не пройдено' : 'В процессе'}
                          </td>
                          <td style={{ padding: '0.5em' }}>{item.score !== undefined ? `${item.score}${item.maxScore ? ' / ' + item.maxScore : ''}` : '-'}</td>
                        </tr>
                      )}</For>
                    </tbody>
                  </table>
                </Show>
              </div>
            )}</For>
          </div>
        </div>
      </>
    );
  }

  return (
    <>
      <Header />
      <div class="main-shell">
        <div class="profile-gradient-block usp" style={{ maxWidth: '420px', margin: '0 auto', boxShadow: '0 4px 24px rgba(37,99,235,0.13)', padding: '2.2em 2em 2em 2em', background: 'var(--bg-primary, #181a1b)', color: 'var(--text-primary, #e3eafc)', border: '1.5px solid #2563eb', transition: 'background 0.2s, color 0.2s' }}>
          <AvatarUploader user={user} onChange={saveAvatar} />
          <h2 style={{ color: '#2563eb', fontWeight: 700, fontSize: '1.25rem', letterSpacing: '0.01em', marginBottom: '1.2em', textAlign: 'center' }}>Профиль пользователя</h2>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.7em' }}>
            {/* Фамилия */}
            <div style={{ display: 'flex', alignItems: 'center', background: 'var(--bg-secondary, #23272a)', borderRadius: '14px', padding: '0.7em 1em', fontSize: '1.08em', color: 'var(--text-primary, #e3eafc)' }}>
              <span style={{ fontSize: '1.2em', marginRight: '0.7em' }}>{userIcons.surname}</span>
              <b style={{ minWidth: '80px', color: '#2563eb' }}>Фамилия:</b>
              {isAdminView && editField() === 'surname' ? (
                <>
                  <input value={editValue()} onInput={e => setEditValue(e.currentTarget.value)} style={{ marginLeft: '0.7em', fontSize: '1em', padding: '0.2em 0.5em', borderRadius: '6px', border: '1px solid #2563eb' }} />
                  <button onClick={() => saveEdit('surname')} style={{ marginLeft: '0.5em', background: '#2563eb', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.2em 0.7em', cursor: 'pointer' }}>✔</button>
                </>
              ) : (
                <>
                  <span style={{ marginLeft: '0.7em' }}>{user.surname}</span>
                  {isAdminView && <span style={{ marginLeft: '0.7em', cursor: 'pointer', color: '#2563eb' }} onClick={() => startEdit('surname', user.surname)}>✏️</span>}
                </>
              )}
            </div>
            {/* Имя */}
            <div style={{ display: 'flex', alignItems: 'center', background: 'var(--bg-secondary, #23272a)', borderRadius: '14px', padding: '0.7em 1em', fontSize: '1.08em', color: 'var(--text-primary, #e3eafc)' }}>
              <span style={{ fontSize: '1.2em', marginRight: '0.7em' }}>{userIcons.name}</span>
              <b style={{ minWidth: '80px', color: '#2563eb' }}>Имя:</b>
              {isAdminView && editField() === 'name' ? (
                <>
                  <input value={editValue()} onInput={e => setEditValue(e.currentTarget.value)} style={{ marginLeft: '0.7em', fontSize: '1em', padding: '0.2em 0.5em', borderRadius: '6px', border: '1px solid #2563eb' }} />
                  <button onClick={() => saveEdit('name')} style={{ marginLeft: '0.5em', background: '#2563eb', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.2em 0.7em', cursor: 'pointer' }}>✔</button>
                </>
              ) : (
                <>
                  <span style={{ marginLeft: '0.7em' }}>{user.name}</span>
                  {isAdminView && <span style={{ marginLeft: '0.7em', cursor: 'pointer', color: '#2563eb' }} onClick={() => startEdit('name', user.name)}>✏️</span>}
                </>
              )}
            </div>
            {/* Email */}
            <div style={{ display: 'flex', alignItems: 'center', background: 'var(--bg-secondary, #23272a)', borderRadius: '14px', padding: '0.7em 1em', fontSize: '1.08em', color: 'var(--text-primary, #e3eafc)' }}>
              <span style={{ fontSize: '1.2em', marginRight: '0.7em' }}>{userIcons.email}</span>
              <b style={{ minWidth: '80px', color: '#2563eb' }}>Email:</b>
              {isAdminView && editField() === 'email' ? (
                <>
                  <input value={editValue()} onInput={e => setEditValue(e.currentTarget.value)} style={{ marginLeft: '0.7em', fontSize: '1em', padding: '0.2em 0.5em', borderRadius: '6px', border: '1px solid #2563eb' }} />
                  <button onClick={() => saveEdit('email')} style={{ marginLeft: '0.5em', background: '#2563eb', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.2em 0.7em', cursor: 'pointer' }}>✔</button>
                </>
              ) : (
                <>
                  <span style={{ marginLeft: '0.7em' }}>{user.email}</span>
                  {isAdminView && <span style={{ marginLeft: '0.7em', cursor: 'pointer', color: '#2563eb' }} onClick={() => startEdit('email', user.email)}>✏️</span>}
                </>
              )}
            </div>
            {/* Роль */}
            <div style={{ display: 'flex', alignItems: 'center', background: 'var(--bg-secondary, #23272a)', borderRadius: '14px', padding: '0.7em 1em', fontSize: '1.08em', color: 'var(--text-primary, #e3eafc)' }}>
              <span style={{ fontSize: '1.2em', marginRight: '0.7em' }}>{userIcons.role}</span>
              <b style={{ minWidth: '80px', color: '#2563eb' }}>Роль:</b>
              {isAdminView && editField() === 'role' ? (
                <>
                  <input value={editValue()} onInput={e => setEditValue(e.currentTarget.value)} style={{ marginLeft: '0.7em', fontSize: '1em', padding: '0.2em 0.5em', borderRadius: '6px', border: '1px solid #2563eb' }} />
                  <button onClick={() => saveEdit('role')} style={{ marginLeft: '0.5em', background: '#2563eb', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.2em 0.7em', cursor: 'pointer' }}>✔</button>
                </>
              ) : (
                <>
                  <span style={{ marginLeft: '0.7em' }}>{user.role}</span>
                  {isAdminView && <span style={{ marginLeft: '0.7em', cursor: 'pointer', color: '#2563eb' }} onClick={() => startEdit('role', user.role)}>✏️</span>}
                </>
              )}
            </div>
          </div>
        </div>
        {user.role === 'teacher' && isAdminView && (
  <div class="profile-gradient-block usp" style={{ background: 'rgba(37,99,235,0.07)', boxShadow: '0 2px 12px rgba(37,99,235,0.07)', transition: 'background 0.2s', marginBottom: '2em', padding: '1.5em' }}>
    <h3 style={{ color: '#2563eb', marginBottom: '1em' }}>Группы преподавателя</h3>
    <ul style={{ paddingLeft: 0, listStyle: 'none', marginBottom: '1em' }}>
      <For each={user.groups || []}>{(gid) => {
        const group = groups.find(g => g.id === gid);
        return (
          <li style={{ marginBottom: '0.5em', display: 'flex', alignItems: 'center' }}>
            <span style={{ marginRight: '1em' }}>{group ? group.name : gid}</span>
            <button style={{ background: '#e76f51', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.2em 0.8em', cursor: 'pointer', fontWeight: 500 }} onClick={() => removeGroupFromTeacher(gid)}>Удалить</button>
          </li>
        );
      }}</For>
    </ul>
    <div style={{ display: 'flex', alignItems: 'center', gap: '1em' }}>
      <select value={selectedGroup()} onInput={e => setSelectedGroup(e.currentTarget.value)} style={{ padding: '0.4em', borderRadius: '6px', border: '1px solid #2563eb' }}>
        <option value="">Выберите группу</option>
        <For each={groups.filter(g => !(user.groups || []).includes(g.id))}>{(g) => (
          <option value={g.id}>{g.name}</option>
        )}</For>
      </select>
      <button disabled={!selectedGroup()} style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.4em 1em', cursor: selectedGroup() ? 'pointer' : 'not-allowed', fontWeight: 500 }} onClick={() => addGroupToTeacher(selectedGroup())}>Добавить группу</button>
    </div>
    <h3 style={{ color: '#2563eb', margin: '1.5em 0 1em 0' }}>Предметы преподавателя</h3>
    <ul style={{ paddingLeft: 0, listStyle: 'none', marginBottom: '1em' }}>
      <For each={user.subjects || []}>{(sid) => {
        const subj = subjects.find(s => s.id === sid);
        return (
          <li style={{ marginBottom: '0.5em', display: 'flex', alignItems: 'center' }}>
            <span style={{ marginRight: '1em' }}>{subj ? subj.name : sid}</span>
            <button style={{ background: '#e76f51', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.2em 0.8em', cursor: 'pointer', fontWeight: 500 }} onClick={() => removeSubjectFromTeacher(sid)}>Удалить</button>
          </li>
        );
      }}</For>
    </ul>
    <div style={{ display: 'flex', alignItems: 'center', gap: '1em' }}>
      <select value={selectedSubject()} onInput={e => setSelectedSubject(e.currentTarget.value)} style={{ padding: '0.4em', borderRadius: '6px', border: '1px solid #2563eb' }}>
        <option value="">Выберите предмет</option>
        <For each={subjects.filter(s => !(user.subjects || []).includes(s.id))}>{(s) => (
          <option value={s.id}>{s.name}</option>
        )}</For>
      </select>
      <button disabled={!selectedSubject()} style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: '6px', padding: '0.4em 1em', cursor: selectedSubject() ? 'pointer' : 'not-allowed', fontWeight: 500 }} onClick={() => addSubjectToTeacher(selectedSubject())}>Добавить предмет</button>
    </div>
  </div>
)}
        {/* Блок успеваемости для ученика или если админ просматривает профиль ученика */}
        <Show when={user.role === 'student' || isAdminView}>
          <div class="profile-gradient-block usp" style={{ background: 'rgba(37,99,235,0.07)', boxShadow: '0 2px 12px rgba(37,99,235,0.07)', transition: 'background 0.2s' }}>
            <h3 style={{ color: '#2563eb', marginBottom: '1.2em', fontSize: '1.18em' }}>Успеваемость</h3>
            {user.activityHistory && user.activityHistory.length > 0 ? (
              (() => {
                const total = user.activityHistory.length;
                const passed = user.activityHistory.filter((item) => item.status === 'passed').length;
                const percent = Math.round((passed / total) * 100);
                return (
                  <>
                    <div style={{ marginBottom: '1.2em', display: 'flex', alignItems: 'center', gap: '1em', flexWrap: 'wrap' }}>
                      <div style={{ flex: 1, minWidth: '180px' }}>
                        <div style={{ height: '18px', background: '#e3eafc', borderRadius: '9px', overflow: 'hidden', position: 'relative' }}>
                          <div style={{ width: percent + '%', height: '100%', background: percent >= 60 ? '#2a9d8f' : '#e76f51', transition: 'width 0.5s', borderRadius: '9px' }}></div>
                        </div>
                        <div style={{ fontSize: '1.1em', marginTop: '0.4em', color: percent >= 60 ? '#2a9d8f' : '#e76f51', fontWeight: 600 }}>{percent}% успешно</div>
      </div>
                      <div style={{ fontSize: '1.2em', color: '#213547', fontWeight: 600 }}>
                        {icons.success} {passed} / {total} {icons.progress}
      </div>
      </div>
                    <table style={{ width: '100%', borderCollapse: 'collapse', background: 'var(--bg-primary)', borderRadius: '8px', overflow: 'hidden', boxShadow: '0 2px 8px rgba(37,99,235,0.10)' }}>
                      <thead>
                        <tr style={{ background: '#e3eafc', color: '#213547' }}>
                          <th style={{ padding: '0.6em', textAlign: 'left' }}>Задание</th>
                          <th style={{ padding: '0.6em', textAlign: 'left' }}>Дата</th>
                          <th style={{ padding: '0.6em', textAlign: 'left' }}>Результат</th>
                          <th style={{ padding: '0.6em', textAlign: 'left' }}>Баллы</th>
                        </tr>
                      </thead>
                      <tbody>
                        {user.activityHistory.map((item, idx) => (
                          <>
                            <tr onClick={() => setSelectedActivityIdx(selectedActivityIdx() === idx ? null : idx)} style={{ cursor: 'pointer', position: 'relative' }}>
                            <td style={{ padding: '0.5em' }}>{item.title}</td>
                            <td style={{ padding: '0.5em' }}>{new Date(item.date).toLocaleString('ru-RU')}</td>
                            <td style={{ padding: '0.5em', color: item.status === 'passed' ? '#2a9d8f' : item.status === 'failed' ? '#e76f51' : '#888' }}>
                              {item.status === 'passed' ? icons.success + ' Пройдено' : item.status === 'failed' ? icons.fail + ' Не пройдено' : icons.progress + ' В процессе'}
                            </td>
                            <td style={{ padding: '0.5em' }}>{item.score !== undefined ? `${item.score}${item.maxScore ? ' / ' + item.maxScore : ''}` : '-'}</td>
                          </tr>
                            {selectedActivityIdx() === idx && (
                              <tr>
                                <td colSpan={4} style={{ padding: 0, border: 'none', background: 'transparent' }}>
                                  <div style={{
                                    margin: '0.5em 0 1em 0',
                                    background: 'var(--bg-secondary)',
                                    'border-radius': '12px',
                                    boxShadow: '0 2px 12px rgba(37,99,235,0.10)',
                                    padding: '1.1em 1.5em',
                                    color: '#213547',
                                    fontSize: '1em',
                                    fontStyle: 'italic',
                                    position: 'relative',
                                    maxWidth: '600px',
                                    'margin-left': 'auto',
                                    'margin-right': 'auto',
                                    zIndex: 10
                                  }}>
                                    <button onClick={() => setSelectedActivityIdx(null)} style={{
                                      position: 'absolute',
                                      top: '0.7em',
                                      right: '1em',
                                      background: 'none',
                                      border: 'none',
                                      color: '#e76f51',
                                      fontSize: '1.3em',
                                      cursor: 'pointer',
                                      fontWeight: 700
                                    }} aria-label="Закрыть">×</button>
                                    <div style={{ marginBottom: '0.5em', fontWeight: 700, color: '#2563eb' }}>{item.title}</div>
                                    <div><b>Дата:</b> {new Date(item.date).toLocaleString('ru-RU')}</div>
                                    <div><b>Результат:</b> {item.status === 'passed' ? 'Пройдено' : item.status === 'failed' ? 'Не пройдено' : 'В процессе'}</div>
                                    <div><b>Баллы:</b> {item.score !== undefined ? `${item.score}${item.maxScore ? ' / ' + item.maxScore : ''}` : '-'}</div>
                                    {/* Детализация по тесту */}
                                    {item.type === 'quiz' && Array.isArray(item.answers) && (
                                      (() => {
                                        const quiz = quizzes.find(q => q.id === item.activityId);
                                        if (!quiz) return <div style={{ color: '#e76f51', marginTop: '1em' }}>Тест не найден</div>;
                                        return (
                                          <div style={{ marginTop: '1em' }}>
                                            <b>Детализация по вопросам:</b>
                                            <ul style={{ paddingLeft: '1.2em', marginTop: '0.5em' }}>
                                              {quiz.questions.map((q, qidx) => (
                                                <li style={{ marginBottom: '0.7em' }}>
                                                  <div><b>Вопрос {qidx + 1}:</b> {q.question}</div>
                                                  <div style={{ marginLeft: '0.7em' }}>
                                                    <b>Ответ ученика:</b> {typeof item.answers[qidx] === 'string' ? item.answers[qidx] : (q.options && typeof item.answers[qidx] === 'number' ? q.options[item.answers[qidx] as number] : '')}
                                                  </div>
                                                  {q.type === 'choice' && q.correct !== undefined && (
                                                    <div style={{ marginLeft: '0.7em' }}>
                                                      <b>Правильный ответ:</b> {q.options ? q.options[q.correct] : ''}
                                                      <span style={{ marginLeft: '1em', color: item.answers[qidx] === q.correct ? '#2a9d8f' : '#e76f51', fontWeight: 600 }}>
                                                        {item.answers[qidx] === q.correct ? '✔' : '✘'}
                                                      </span>
                                                    </div>
                                                  )}
                                                </li>
                                              ))}
                                            </ul>
                                          </div>
                                        );
                                      })()
                                    )}
                                    {/* Приложенные материалы (attachments или files) */}
                                    {Array.isArray((item as any).attachments) && (item as any).attachments.length > 0 && (
                                      <div style={{ 'margin-top': '0.7em', color: '#2563eb' }}>
                                        <b>Приложенные материалы:</b>
                                        <ul style={{ paddingLeft: '1.2em', marginTop: '0.3em' }}>
                                          {(item as any).attachments.map((att: any) => (
                                            <li>
                                              <span style={{ color: '#2563eb', fontWeight: 600, marginRight: '0.5em' }}>{att.type ? att.type.toUpperCase() : ''}</span>
                                              {att.url ? (
                                                <a href={att.url} target="_blank" rel="noopener noreferrer" style={{ color: '#e76f51', textDecoration: 'underline' }}>{att.url}</a>
                                              ) : (
                                                <span>{att.name || att}</span>
                                              )}
                                            </li>
                                          ))}
                                        </ul>
                                      </div>
                                    )}
                                    {Array.isArray((item as any).files) && (item as any).files.length > 0 && (
                                      <div style={{ 'margin-top': '0.7em', color: '#2563eb' }}>
                                        <b>Приложенные файлы:</b>
                                        <ul style={{ paddingLeft: '1.2em', marginTop: '0.3em' }}>
                                          {(item as any).files.map((file: string) => (
                                            <li>{file}</li>
                                          ))}
                                        </ul>
                                      </div>
                                    )}
                                    {/* Комментарий учителя */}
                                    {(item as any).comment && (
                                      <div style={{ 'margin-top': '1em', color: '#213547', background: 'var(--bg-primary)', 'border-radius': '8px', padding: '0.8em 1em', fontSize: '1em' }}><b>Комментарий учителя:</b> {(item as any).comment}</div>
                                    )}
                                  </div>
                                </td>
                              </tr>
                            )}
                          </>
                        ))}
                      </tbody>
                    </table>
            </>
                );
              })()
          ) : (
              <div style={{ color: '#888', fontSize: '1.1em', textAlign: 'center', margin: '2em 0' }}>Нет данных по успеваемости</div>
          )}
        </div>
      </Show>
      <Show when={user.role === 'teacher'}>
          <div class="usp" style={{ background: 'rgba(37,99,235,0.07)', boxShadow: '0 2px 12px rgba(37,99,235,0.07)', transition: 'background 0.2s', padding: '1.2em 0.5em', marginBottom: '2em' }}>
          <TeacherClassPerformance user={user.role === 'student' ? user as Student : user} />
          </div>
        </Show>
        {isAdminView && (
  <div style={{ margin: '2em 0', textAlign: 'center' }}>
    <button style={{ background: '#e76f51', color: '#fff', border: 'none', borderRadius: '8px', padding: '0.7em 2em', fontWeight: 700, fontSize: '1.1em', cursor: 'pointer' }} onClick={() => setShowDeleteModal(true)}>Удалить пользователя</button>
    <Show when={showDeleteModal()}>
      <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', background: 'rgba(0,0,0,0.25)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div style={{ background: '#fff', borderRadius: '12px', padding: '2em 2.5em', boxShadow: '0 4px 24px rgba(37,99,235,0.18)', textAlign: 'center' }}>
          <div style={{ fontSize: '1.2em', marginBottom: '1.2em', color: '#e76f51', fontWeight: 700 }}>Вы уверены, что хотите удалить пользователя?</div>
          <button style={{ background: '#e76f51', color: '#fff', border: 'none', borderRadius: '8px', padding: '0.7em 2em', fontWeight: 700, fontSize: '1.1em', cursor: 'pointer', marginRight: '1em' }} onClick={deleteUser}>Удалить</button>
          <button style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: '8px', padding: '0.7em 2em', fontWeight: 700, fontSize: '1.1em', cursor: 'pointer' }} onClick={() => setShowDeleteModal(false)}>Отмена</button>
        </div>
      </div>
    </Show>
              </div>
          )}
        </div>
    </>
  );
};

export default Profile; 