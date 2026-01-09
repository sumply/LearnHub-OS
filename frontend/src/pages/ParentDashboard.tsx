import { type Component } from 'solid-js';
import { A } from '@solidjs/router';
import { getCurrentUser } from '../utils/api';

const ParentDashboard: Component = () => {
  const user = getCurrentUser();
  const isParent = user?.role === 'parent';

  return (
    <div class="dashboard">
      <h2>Добро пожаловать, родитель!</h2>
      <ul>
        <li><A href="/child-progress">Успеваемость ребёнка</A></li>
        <li><A href="/materials">Учебные материалы</A></li>
        <li><A href="/flashcards">Карточки</A></li>
        <li><A href="/quiz">Тесты</A></li>
      </ul>
      {isParent && (
        <div style={{'margin-top': '2rem'}}>
          <button style={{padding: '10px 20px', background: '#2a9d8f', color: '#fff', border: 'none', 'border-radius': '5px', cursor: 'pointer', 'margin-right': '1rem'}}>Посмотреть комментарии учителя</button>
          <button style={{padding: '10px 20px', background: '#e76f51', color: '#fff', border: 'none', 'border-radius': '5px', cursor: 'pointer'}}>Связаться с учителем</button>
        </div>
      )}
    </div>
  );
};

export default ParentDashboard; 