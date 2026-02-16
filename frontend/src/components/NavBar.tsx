import { getCurrentUser } from '../utils/api';
import { createSignal, onMount } from 'solid-js';
import { A } from '@solidjs/router';

const NavBar = () => {
  const [user, setUser] = createSignal<any>(null);

  onMount(() => {
    setUser(getCurrentUser());
  });

  const role = () => user()?.role || null;

  return (
    <nav class="navbar" style={{ background: 'transparent' }}>
      {!user() && (
        <>
          <A class="nav-btn" href="/login" style={{ 'min-width': '120px', 'margin-right': '0.5rem' }}>Войти</A>
          {/* Регистрация временно отключена */}
          {/* <A class="nav-btn" href="/register" style={{ 'min-width': '120px', 'margin-right': '0.5rem' }}>Зарегистрироваться</A> */}
        </>
      )}
      {user() && (
        <>
          <A class="nav-btn" href="/library" style={{ 'min-width': '120px', 'margin-right': '0.5rem' }}>Библиотека</A>
          <A class="nav-btn" href="/tasks" style={{ 'min-width': '120px', 'margin-right': '0.5rem' }}>Интерактивные задания</A>
          {role() === 'student' && <A class="nav-btn" href="/progress" style={{ 'min-width': '120px', 'margin-right': '0.5rem' }}>Моя успеваемость</A>}
          {role() === 'parent' && <A class="nav-btn" href="/child-progress" style={{ 'min-width': '120px', 'margin-right': '0.5rem' }}>Успеваемость ребёнка</A>}
          {/* {role() === 'teacher' && <A class="nav-btn" href="/teacher/groups" style={{ 'min-width': '120px', 'margin-right': '0.5rem' }}>Мои классы</A>} */}
          {/* {role() === 'teacher' && <A class="nav-btn" href="/teacher/materials" style={{ 'min-width': '120px', 'margin-right': '0.5rem' }}>Мои материалы</A>} */}
          {role() === 'admin' && <A class="nav-btn" href="/admin" style={{ 'min-width': '120px', 'margin-right': '0.5rem' }}>Админка</A>}
          {role() !== 'admin' && <A class="nav-btn" href="/profile" style={{ 'min-width': '120px', 'margin-right': '0.5rem' }}>Личный кабинет</A>}
        </>
      )}
    </nav>
  );
};

export default NavBar; 