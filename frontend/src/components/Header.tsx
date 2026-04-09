import type { Component } from 'solid-js';
import { createSignal, onMount } from 'solid-js';
import NavBar from './NavBar';
import AccessibilitySettings from './AccessibilitySettings';
import { getCurrentUser, logout } from '../utils/api';
import { useTheme } from '../contexts/ThemeContext';
import { A } from '@solidjs/router';

const Header: Component = () => {
  const [user, setUser] = createSignal<any>(null);
  const { theme } = useTheme();

  onMount(() => {
    const currentUser = getCurrentUser();
    setUser(currentUser);
  });

  const handleLogout = async () => {
    try {
      await logout();
      window.location.href = '/login';
    } catch (error) {
      console.error('Ошибка при выходе:', error);
      window.location.href = '/login';
    }
  };

  const getRoleText = (role: string) => {
    switch (role) {
      case 'admin': return 'Администратор';
      case 'teacher': return 'Учитель';
      case 'parent': return 'Родитель';
      case 'student': return 'Ученик';
      default: return 'Пользователь';
    }
  };

  return (
    <>
      <header class="header">
        <div class="header-left">
          <A href="/" class="header-logo" aria-label="Мир Квизов — на главную" style={{ display: 'flex', 'align-items': 'center', gap: 'clamp(0.3rem, 1vw, 0.7em)', 'text-decoration': 'none' }}>
            <span class="logo-icon-wrap" aria-hidden="true">
              <img
                class="logo-icon"
                src={theme() === 'dark' ? '/mir_quiz_log1_black.webp' : '/mir_quiz_log1.webp'}
                alt=""
                decoding="async"
                loading="eager"
              />
            </span>
            <span style={{ 'font-size': 'var(--font-size-xl)', 'font-weight': 700, color: '#2563eb', 'white-space': 'nowrap', 'letter-spacing': '0.01em', 'line-height': 1 }}>Мир Квизов</span>
          </A>
        </div>
        
        <div class="header-center">
          <NavBar />
          <div class="header-actions">
            <AccessibilitySettings />
          </div>
        </div>
        
        <div class="header-right">
          {user() && (
            <div class="user-section">
              <div class="user-info">
                <div class="user-name">
                  {user().name} {user().surname}
                </div>
                <div class="user-role">
                  {getRoleText(user().role)}
                  {user().group && <span class="user-group"> • {user().group}</span>}
                </div>
              </div>
              
              <button
                class="logout-btn"
                onClick={handleLogout}
                aria-label="Выйти из системы"
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                  <polyline points="16,17 21,12 16,7"/>
                  <line x1="21" y1="12" x2="9" y2="12"/>
                </svg>
                <span>Выйти</span>
              </button>
            </div>
          )}
        </div>
      </header>
    </>
  );
};

export default Header; 