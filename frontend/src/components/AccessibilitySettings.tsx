import type { Component } from 'solid-js';
import { createSignal, Show } from 'solid-js';
import { useTheme } from '../contexts/ThemeContext';

const ThemeToggle = () => {
  const { theme, toggleTheme } = useTheme();
  return (
    <button
      class={`theme-toggle ${theme() === 'dark' ? 'dark' : 'light'}`}
      aria-label={theme() === 'dark' ? 'Включить светлую тему' : 'Включить тёмную тему'}
      onClick={toggleTheme}
      type="button"
    >
      <span class="toggle-track">
        <span class="toggle-thumb">
          {theme() === 'dark' ? (
            <svg width="22" height="22" viewBox="0 0 22 22" fill="none"><path d="M16.5 14.5C13.5 16.5 9.5 15.5 8 12.5C6.5 9.5 8.5 5.5 12 4.5C11.5 6 12 8.5 14 10.5C16 12.5 18.5 13 16.5 14.5Z" fill="#fff"/><circle cx="16" cy="7" r="1" fill="#fff"/><circle cx="18" cy="10" r="0.5" fill="#fff"/></svg>
          ) : (
            <svg width="22" height="22" viewBox="0 0 22 22" fill="none"><circle cx="11" cy="11" r="5" fill="#fff"/><g stroke="#fff" stroke-width="1.5"><line x1="11" y1="2" x2="11" y2="5"/><line x1="11" y1="17" x2="11" y2="20"/><line x1="2" y1="11" x2="5" y2="11"/><line x1="17" y1="11" x2="20" y2="11"/><line x1="4.22" y1="4.22" x2="6.34" y2="6.34"/><line x1="15.66" y1="15.66" x2="17.78" y2="17.78"/><line x1="4.22" y1="17.78" x2="6.34" y2="15.66"/><line x1="15.66" y1="6.34" x2="17.78" y2="4.22"/></g></svg>
          )}
        </span>
      </span>
    </button>
  );
};

const AccessibilitySettings: Component = () => {
  const { fontSize, setFontSize, scale, setScale, highContrast, toggleHighContrast } = useTheme();
  const [isOpen, setIsOpen] = createSignal(false);

  const toggleSettings = () => setIsOpen(!isOpen());

  return (
    <div class="accessibility-settings">
      <button 
        class="accessibility-toggle"
        onClick={toggleSettings}
        aria-label="Настройки доступности"
        title="Настройки доступности"
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/>
          <path d="M12 1v6m0 6v6m11-7h-6m-6 0H1"/>
        </svg>
      </button>
      <ThemeToggle />
      <Show when={isOpen()}>
        <div class="accessibility-panel">
          <div class="accessibility-header">
            <h3>Настройки доступности</h3>
            <button 
              class="close-btn"
              onClick={toggleSettings}
              aria-label="Закрыть настройки"
            >
              ×
            </button>
          </div>

          <div class="accessibility-content">
            {/* Размер текста */}
            <div class="setting-group">
              <label for="font-size">Размер текста: {fontSize()}px</label>
              <input
                id="font-size"
                type="range"
                min="12"
                max="24"
                step="1"
                value={fontSize()}
                onInput={(e) => setFontSize(parseInt(e.currentTarget.value))}
              />
              <div class="range-labels">
                <span>12px</span>
                <span>24px</span>
              </div>
            </div>

            {/* Масштаб */}
            <div class="setting-group">
              <label for="scale">Масштаб: {Math.round(scale() * 100)}%</label>
              <input
                id="scale"
                type="range"
                min="0.8"
                max="1.5"
                step="0.1"
                value={scale()}
                onInput={(e) => setScale(parseFloat(e.currentTarget.value))}
              />
              <div class="range-labels">
                <span>80%</span>
                <span>150%</span>
              </div>
            </div>

            {/* Высокий контраст */}
            <div class="setting-group">
              <label class="toggle-label">
                <input
                  type="checkbox"
                  checked={highContrast()}
                  onChange={toggleHighContrast}
                />
                <span class="toggle-slider"></span>
                Высокий контраст
              </label>
            </div>

            {/* Быстрые кнопки */}
            <div class="quick-actions">
              <button 
                class="quick-btn"
                onClick={() => setFontSize(16)}
              >
                Сброс текста
              </button>
              <button 
                class="quick-btn"
                onClick={() => setScale(1)}
              >
                Сброс масштаба
              </button>
            </div>
          </div>
        </div>
      </Show>
    </div>
  );
};

export default AccessibilitySettings; 