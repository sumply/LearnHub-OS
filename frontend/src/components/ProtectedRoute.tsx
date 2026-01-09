/*
import type { Component, JSX } from 'solid-js';
import { createSignal, onMount, Show } from 'solid-js';
import { isAuthenticated, getCurrentUser } from '../utils/api';

interface ProtectedRouteProps {
  children: JSX.Element;
  fallback?: JSX.Element;
  role?: string | string[]; // допустимая роль или массив ролей
}

const ProtectedRoute: Component<ProtectedRouteProps> = (props) => {
  const [isLoading, setIsLoading] = createSignal(true);
  const [isAuthorized, setIsAuthorized] = createSignal(false);
  const [roleError, setRoleError] = createSignal(false);

  onMount(() => {
    if (!isAuthenticated()) {
      setIsAuthorized(false);
      setIsLoading(false);
      return;
    }
    const user = getCurrentUser();
    if (!user) {
      setIsAuthorized(false);
      setIsLoading(false);
      return;
    }
    if (props.role) {
      if (Array.isArray(props.role)) {
        if (!props.role.includes(user.role)) {
          setRoleError(true);
          setIsAuthorized(false);
          setIsLoading(false);
          return;
        }
      } else {
        if (user.role !== props.role) {
          setRoleError(true);
          setIsAuthorized(false);
          setIsLoading(false);
          return;
        }
      }
    }
    setIsAuthorized(true);
    setIsLoading(false);
  });

  return (
    <Show
      when={!isLoading()}
      fallback={<div style={{'text-align': 'center', 'margin-top': '2rem'}}>Загрузка...</div>}
    >
      <Show
        when={isAuthorized()}
        fallback={
          roleError() ? (
            <div style={{'display': 'flex', 'flex-direction': 'column', 'align-items': 'center', 'margin-top': '2rem'}}>
              <div class="error-message" style={{'max-width': '320px'}}>Нет доступа к этой странице</div>
              <button class="nav-btn" style={{'margin-top': '1.2rem'}} onClick={() => window.location.href = '/'}>На главную</button>
            </div>
          ) : (
            props.fallback || (
              <div style={{'display': 'flex', 'flex-direction': 'column', 'align-items': 'center', 'margin-top': '2rem'}}>
                <div class="error-message" style={{'max-width': '320px'}}>Требуется авторизация</div>
                <button class="nav-btn" style={{'margin-top': '1.2rem'}} onClick={() => window.location.href = '/login'}>Войти</button>
              </div>
            )
          )
        }
      >
        {props.children}
      </Show>
    </Show>
  );
};

export default ProtectedRoute;
*/

import { type Component, Show } from 'solid-js';
import { getCurrentUser } from '../utils/api';
import { useNavigate } from '@solidjs/router';

interface ProtectedRouteProps {
  children: any;
  role?: string;
}

const ProtectedRoute: Component<ProtectedRouteProps> = (props) => {
  const user = getCurrentUser();
  const navigate = useNavigate();

  if (!user) {
    navigate('/login');
    return null;
  }
  if (props.role && user.role !== props.role) {
    navigate('/');
    return null;
  }
  return <Show when={!!user}>{props.children}</Show>;
};

export default ProtectedRoute; 