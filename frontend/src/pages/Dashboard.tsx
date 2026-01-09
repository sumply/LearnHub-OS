import type { Component } from 'solid-js';
import Header from '../components/Header';
import { getCurrentUser, removeAuthToken, removeRefreshToken, removeCurrentUser } from '../utils/api';

const Dashboard: Component = () => {
  const user = getCurrentUser();

  const handleLogout = () => {
    removeAuthToken();
    removeRefreshToken();
    removeCurrentUser();
    window.location.href = '/login';
  };

  if (!user) {
    return (
      <div style="text-align:center; margin-top:3rem; color:#e76f51; font-size:1.2em;">Пользователь не найден. Пожалуйста, войдите в систему.</div>
    );
  }

  return (
    <>
      <Header />
      <div class="main-shell">
        {/* основной контент Dashboard */}
        <div style="max-width:800px; margin:2rem auto; padding:2rem; background:#fff; border-radius:16px; box-shadow:0 2px 12px rgba(0,0,0,0.07);">
          <h1 style="margin-bottom:1.5rem; color:#e76f51;">Личный кабинет</h1>
          <section style="margin-bottom:2.5rem;">
            <h2 style="font-size:1.3em; margin-bottom:0.7rem;">Профиль</h2>
            <div><b>ФИО:</b> {user.surname} {user.name}</div>
            <div><b>Email:</b> {user.email}</div>
            <div><b>Роль:</b> {user.role === 'student' ? 'Ученик' : user.role === 'parent' ? 'Родитель' : user.role === 'teacher' ? 'Учитель' : user.role === 'admin' ? 'Администратор' : user.role}</div>
            {user.group && <div><b>Класс/группа:</b> {user.group}</div>}
            <button style="margin-top:1rem; padding:0.5rem 1.2rem; border-radius:8px; border:none; background:#f7c873; color:#213547; font-weight:600; cursor:pointer; margin-right:1rem;" disabled>Редактировать профиль</button>
            <button style="margin-top:1rem; padding:0.5rem 1.2rem; border-radius:8px; border:none; background:#2563eb; color:#fff; font-weight:600; cursor:pointer;" onClick={handleLogout}>Выйти</button>
          </section>
          {/* Здесь можно добавить отображение выполненных заданий, материалов и т.д. */}
          {user.role === 'student' && user.activityHistory && user.activityHistory.length > 0 && (
            <section style={{ marginTop: '2.5rem' }}>
              <h2 style={{ fontSize: '1.2em', marginBottom: '0.7rem', color: '#2563eb' }}>История прохождения заданий</h2>
              <table style={{ width: '100%', borderCollapse: 'collapse', background: '#f9fafb', borderRadius: '8px', overflow: 'hidden' }}>
                <thead>
                  <tr style={{ background: '#e3eafc', color: '#213547' }}>
                    <th style={{ padding: '0.6em', textAlign: 'left' }}>Задание</th>
                    <th style={{ padding: '0.6em', textAlign: 'left' }}>Дата</th>
                    <th style={{ padding: '0.6em', textAlign: 'left' }}>Результат</th>
                    <th style={{ padding: '0.6em', textAlign: 'left' }}>Баллы</th>
                  </tr>
                </thead>
                <tbody>
                  {user.activityHistory.map((item: any) => (
                    <tr>
                      <td style={{ padding: '0.5em' }}>{item.title}</td>
                      <td style={{ padding: '0.5em' }}>{new Date(item.date).toLocaleString('ru-RU')}</td>
                      <td style={{ padding: '0.5em', color: item.status === 'passed' ? '#2a9d8f' : item.status === 'failed' ? '#e76f51' : '#888' }}>
                        {item.status === 'passed' ? 'Пройдено' : item.status === 'failed' ? 'Не пройдено' : 'В процессе'}
                      </td>
                      <td style={{ padding: '0.5em' }}>{item.score !== undefined ? `${item.score}${item.maxScore ? ' / ' + item.maxScore : ''}` : '-'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </section>
          )}
        </div>
      </div>
    </>
  );
};

export default Dashboard; 