import { createContext, createSignal, useContext, onMount, createEffect } from 'solid-js';
import type { ParentComponent } from 'solid-js';

interface ThemeContextType {
  theme: () => 'light' | 'dark';
  toggleTheme: () => void;
  fontSize: () => number;
  setFontSize: (size: number) => void;
  scale: () => number;
  setScale: (scale: number) => void;
  highContrast: () => boolean;
  toggleHighContrast: () => void;
}

const ThemeContext = createContext<ThemeContextType>();

export const ThemeProvider: ParentComponent = (props) => {
  const [theme, setTheme] = createSignal<'light' | 'dark'>(
    localStorage.getItem('theme') as 'light' | 'dark' || 'light'
  );
  
  const [fontSize, setFontSize] = createSignal(
    parseInt(localStorage.getItem('fontSize') || '16')
  );
  
  const [scale, setScale] = createSignal(
    parseFloat(localStorage.getItem('scale') || '1')
  );
  
  const [highContrast, setHighContrast] = createSignal(
    localStorage.getItem('highContrast') === 'true'
  );

  const toggleTheme = () => {
    const newTheme = theme() === 'light' ? 'dark' : 'light';
    setTheme(newTheme);
    localStorage.setItem('theme', newTheme);
  };

  const updateFontSize = (size: number) => {
    setFontSize(size);
    localStorage.setItem('fontSize', size.toString());
    document.documentElement.style.fontSize = `${size}px`;
  };

  const updateScale = (newScale: number) => {
    setScale(newScale);
    localStorage.setItem('scale', newScale.toString());
    document.documentElement.style.setProperty('--scale', newScale.toString());
  };

  const toggleHighContrast = () => {
    const newValue = !highContrast();
    setHighContrast(newValue);
    localStorage.setItem('highContrast', newValue.toString());
  };

  // Гарантируем, что тема всегда применяется к <html>
  onMount(() => {
    document.documentElement.setAttribute('data-theme', theme());
    document.documentElement.style.fontSize = `${fontSize()}px`;
    document.documentElement.style.setProperty('--scale', scale().toString());
    document.documentElement.setAttribute('data-high-contrast', highContrast().toString());
  });

  createEffect(() => {
    document.documentElement.setAttribute('data-theme', theme());
  });

  createEffect(() => {
    document.documentElement.setAttribute('data-high-contrast', highContrast().toString());
  });

  createEffect(() => {
    document.documentElement.style.fontSize = `${fontSize()}px`;
  });

  createEffect(() => {
    document.documentElement.style.setProperty('--scale', scale().toString());
  });

  const value: ThemeContextType = {
    theme,
    toggleTheme,
    fontSize,
    setFontSize: updateFontSize,
    scale,
    setScale: updateScale,
    highContrast,
    toggleHighContrast,
  };

  return (
    <ThemeContext.Provider value={value}>
      {props.children}
    </ThemeContext.Provider>
  );
};

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}; 