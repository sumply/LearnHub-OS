import { type Component } from 'solid-js';
import { A } from '@solidjs/router';
import { getCurrentUser } from '../utils/api';

const StudentDashboard: Component = () => {
  const user = getCurrentUser();
  const isStudent = user?.role === 'student';

  return (
    <div class="dashboard">
      <h2>Добро пожаловать, ученик!</h2>
      <ul>
        <li><A href="/materials">Учебные материалы</A></li>
        <li><A href="/tasks">Интерактивные задания</A></li>
        <li><A href="/progress">Моя успеваемость</A></li>
        <li><A href="/flashcards">Карточки</A></li>
        <li><A href="/quiz">Тесты</A></li>
      </ul>
      {isStudent && (
        <div style={{'margin-top': '2rem'}}>
          <button style={{padding: '10px 20px', background: '#2a9d8f', color: '#fff', border: 'none', 'border-radius': '5px', cursor: 'pointer', 'margin-right': '1rem'}}>Я прошёл урок</button>
          <button style={{padding: '10px 20px', background: '#e76f51', color: '#fff', border: 'none', 'border-radius': '5px', cursor: 'pointer', 'margin-right': '1rem'}}>Начать тест</button>
          <button style={{padding: '10px 20px', background: '#264653', color: '#fff', border: 'none', 'border-radius': '5px', cursor: 'pointer'}}>Задать вопрос учителю</button>
        </div>
      )}
    </div>
  );
};

export default StudentDashboard; 