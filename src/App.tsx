import type { Component } from 'solid-js';
import { ThemeProvider } from './contexts/ThemeContext';
import Routes from './routes';

const App: Component = () => {
  return (
    <ThemeProvider>
      <Routes />
    </ThemeProvider>
  );
};

export default App;
