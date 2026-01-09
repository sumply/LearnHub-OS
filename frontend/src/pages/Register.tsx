import type { Component } from 'solid-js';
import { createSignal } from 'solid-js';
import Header from '../components/Header';
import { addUser, isEmailTaken } from '../config/users';
import type { UserRole, User } from '../config/users';
import { useNavigate } from '@solidjs/router';

const Register: Component = () => {
  const [role, setRole] = createSignal('student');
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
  const navigate = useNavigate();

  const handleRegister = async (e: Event) => {
    e.preventDefault();
    setError('');
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
    if (isEmailTaken(regEmail())) {
      setError('Пользователь с таким email уже существует!');
      return;
    }
    setIsLoading(true);
    const base = {
      id: Date.now().toString(),
      email: regEmail(),
      password: regPassword(),
      name: regFirstname(),
      surname: regLastname(),
      role: role() as UserRole,
    };
    let newUser: User;
    if (role() === 'admin') newUser = { ...base, role: 'admin' };
    else if (role() === 'teacher') {
      newUser = { ...base, role: 'teacher', subjects: [] };
    } else if (role() === 'parent') {
      newUser = { ...base, role: 'parent', childrenIds: [] };
    } else {
      newUser = { ...base, role: 'student', group: regGroup() };
    }
    try {
      addUser(newUser);
      alert('Регистрация успешна! Теперь вы можете войти в систему.');
      navigate('/login');
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
          <h2>Регистрация</h2>
          {error() && (
            <div class="error-message">{error()}</div>
          )}
          <form class="auth-form" onSubmit={handleRegister}>
            <label>Email</label>
            <input type="email" placeholder="Введите email" required value={regEmail()} onInput={e => setRegEmail(e.currentTarget.value)} disabled={isLoading()} />
            <label>Подтвердите Email</label>
            <input type="email" placeholder="Повторите email" required value={regEmailConfirm()} onInput={e => setRegEmailConfirm(e.currentTarget.value)} disabled={isLoading()} />
            <label>Пароль</label>
            <input type="password" placeholder="Введите пароль (минимум 6 символов)" required value={regPassword()} onInput={e => setRegPassword(e.currentTarget.value)} disabled={isLoading()} />
            <label>Подтвердите пароль</label>
            <input type="password" placeholder="Повторите пароль" required value={regPasswordConfirm()} onInput={e => setRegPasswordConfirm(e.currentTarget.value)} disabled={isLoading()} />
            <label>Фамилия (lastname)</label>
            <input type="text" placeholder="Ваша фамилия" required value={regLastname()} onInput={e => setRegLastname(e.currentTarget.value)} disabled={isLoading()} />
            <label>Имя (firstname)</label>
            <input type="text" placeholder="Ваше имя" required value={regFirstname()} onInput={e => setRegFirstname(e.currentTarget.value)} disabled={isLoading()} />
            <label>Отчество (secondname)</label>
            <input type="text" placeholder="Ваше отчество (необязательно)" value={regSecondname()} onInput={e => setRegSecondname(e.currentTarget.value)} disabled={isLoading()} />
            <label>Я —</label>
            <select value={role()} onInput={e => setRole(e.currentTarget.value)} required disabled={isLoading()}>
              <option value="student">Ученик</option>
              <option value="parent">Родитель</option>
              <option value="teacher">Учитель</option>
              <option value="admin">Администратор</option>
            </select>
            {role() !== 'teacher' && (
              <>
                <label>Класс/группа</label>
                <input type="text" placeholder="Например, 11А" required value={regGroup()} onInput={e => setRegGroup(e.currentTarget.value)} disabled={isLoading()} />
              </>
            )}
            <button type="submit" class="auth-btn" disabled={isLoading()}>
              {isLoading() ? 'Регистрация...' : 'Зарегистрироваться'}
            </button>
            <div class="auth-switch">
              Уже есть аккаунт?{' '}
              <a href="/login">Войти</a>
            </div>
          </form>
        </div>
      </div>
    </>
  );
};

export default Register; 