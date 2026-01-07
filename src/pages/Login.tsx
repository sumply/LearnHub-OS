import type { Component } from 'solid-js';
import { createSignal } from 'solid-js';
import Header from '../components/Header';
import { login, register, handleAuthError } from '../utils/api';
import type { UserRole } from '../config/users';

const Login: Component = () => {
  const [isRegister, setIsRegister] = createSignal(false);
  const [role, setRole] = createSignal('student');
  const [email, setEmail] = createSignal('');
  const [password, setPassword] = createSignal('');
  const [regEmail, setRegEmail] = createSignal('');
  const [regEmailConfirm, setRegEmailConfirm] = createSignal('');
  const [regPassword, setRegPassword] = createSignal('');
  const [regPasswordConfirm, setRegPasswordConfirm] = createSignal('');
  const [regFirstname, setRegFirstname] = createSignal('');
  const [regSecondname, setRegSecondname] = createSignal('');
  const [regLastname, setRegLastname] = createSignal('');
  const [regGroup, setRegGroup] = createSignal('');
  const [isLoading, setIsLoading] = createSignal(false);
  const [error, setError] = createSignal('');

  // Обработчик входа
  const handleLogin = async (e: Event) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      // Используем API для входа
      const result = await login({
        email: email(),
        password: password(),
      });

      if (result.success) {
        // Перенаправляем на главную страницу
        window.location.href = '/';
      } else {
        setError(result.message || 'Ошибка входа');
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка входа';
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  // Обработчик регистрации
  const handleRegister = async (e: Event) => {
    e.preventDefault();
    setError('');

    // Валидация
    if (regEmail() !== regEmailConfirm()) {
      setError('Email не совпадает!');
      return;
    }
    if (regPassword() !== regPasswordConfirm()) {
      setError('Пароли не совпадают!');
      return;
    }
    if (regPassword().length < 6) {
      setError('Пароль должен содержать минимум 6 символов');
      return;
    }
    if (!regFirstname() || !regLastname()) {
      setError('Имя и фамилия обязательны для заполнения');
      return;
    }

    setIsLoading(true);

    try {
      // Используем API для регистрации
      const result = await register({
        email: regEmail(),
        password: regPassword(),
        name: regFirstname(),
        surname: regLastname(),
        role: role() as UserRole,
        group: regGroup() || undefined,
      });

      if (result.success) {
        alert('Регистрация успешна! Логин и пароль будут сгенерированы сервером. Ожидайте письмо с данными для входа.');
        setIsRegister(false);
        // Очищаем форму регистрации
        setRegEmail('');
        setRegEmailConfirm('');
        setRegPassword('');
        setRegPasswordConfirm('');
        setRegFirstname('');
        setRegSecondname('');
        setRegLastname('');
        setRegGroup('');
      } else {
        setError(result.message || 'Ошибка регистрации');
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка регистрации';
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <>
      <Header />
      <div class="auth-container">
        <div class="auth-card">
          <h2>
            {isRegister() ? 'Регистрация' : 'Вход'}
          </h2>
          
          {error() && (
            <div class="error-message">
              {error()}
            </div>
          )}

          {!isRegister() ? (
            <form class="auth-form" onSubmit={handleLogin}>
              <label>Email / Логин</label>
              <input 
                type="text" 
                placeholder="Введите email или логин" 
                // required 
                value={email()} 
                onInput={e => setEmail(e.currentTarget.value)}
                disabled={isLoading()}
              />
              <label>Пароль</label>
              <input 
                type="password" 
                placeholder="Введите пароль" 
                // required 
                value={password()} 
                onInput={e => setPassword(e.currentTarget.value)}
                disabled={isLoading()}
              />
              <button 
                type="submit" 
                class="auth-btn" 
                disabled={isLoading()}
              >
                {isLoading() ? 'Вход...' : 'Войти'}
              </button>
              <div class="auth-switch">
                Нет аккаунта?{' '}
                <a href="#" onClick={e => { e.preventDefault(); setIsRegister(true); setError(''); }}>Зарегистрироваться</a>
              </div>
            </form>
          ) : (
            <form class="auth-form" onSubmit={handleRegister}>
              <label>Email</label>
              <input 
                type="email" 
                placeholder="Введите email" 
                required 
                value={regEmail()} 
                onInput={e => setRegEmail(e.currentTarget.value)}
                disabled={isLoading()}
              />
              <label>Подтвердите Email</label>
              <input 
                type="email" 
                placeholder="Повторите email" 
                required 
                value={regEmailConfirm()} 
                onInput={e => setRegEmailConfirm(e.currentTarget.value)}
                disabled={isLoading()}
              />
              <label>Пароль</label>
              <input 
                type="password" 
                placeholder="Введите пароль (минимум 6 символов)" 
                required 
                value={regPassword()} 
                onInput={e => setRegPassword(e.currentTarget.value)}
                disabled={isLoading()}
              />
              <label>Подтвердите пароль</label>
              <input 
                type="password" 
                placeholder="Повторите пароль" 
                required 
                value={regPasswordConfirm()} 
                onInput={e => setRegPasswordConfirm(e.currentTarget.value)}
                disabled={isLoading()}
              />
              <label>Фамилия (lastname)</label>
              <input 
                type="text" 
                placeholder="Ваша фамилия" 
                required 
                value={regLastname()} 
                onInput={e => setRegLastname(e.currentTarget.value)}
                disabled={isLoading()}
              />
              <label>Имя (firstname)</label>
              <input 
                type="text" 
                placeholder="Ваше имя" 
                required 
                value={regFirstname()} 
                onInput={e => setRegFirstname(e.currentTarget.value)}
                disabled={isLoading()}
              />
              <label>Отчество (secondname)</label>
              <input 
                type="text" 
                placeholder="Ваше отчество (необязательно)" 
                value={regSecondname()} 
                onInput={e => setRegSecondname(e.currentTarget.value)}
                disabled={isLoading()}
              />
              <label>Я —</label>
              <select 
                value={role()} 
                onInput={e => setRole(e.currentTarget.value)} 
                required
                disabled={isLoading()}
              >
                <option value="student">Ученик</option>
                <option value="parent">Родитель</option>
                <option value="teacher">Учитель</option>
                <option value="admin">Администратор</option>
              </select>
              {role() !== 'teacher' && (
                <>
                  <label>Класс/группа</label>
                  <input 
                    type="text" 
                    placeholder="Например, 11А" 
                    required 
                    value={regGroup()} 
                    onInput={e => setRegGroup(e.currentTarget.value)}
                    disabled={isLoading()}
                  />
                </>
              )}
              <button 
                type="submit" 
                class="auth-btn" 
                disabled={isLoading()}
              >
                {isLoading() ? 'Регистрация...' : 'Зарегистрироваться'}
              </button>
              <div class="auth-switch">
                Уже есть аккаунт?{' '}
                <a href="#" onClick={e => { e.preventDefault(); setIsRegister(false); setError(''); }}>Войти</a>
              </div>
            </form>
          )}
        </div>
      </div>
    </>
  );
};

export default Login; 