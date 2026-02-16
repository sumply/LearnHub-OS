import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'

export default defineConfig({
  plugins: [solid()],
  server: {
    proxy: {
      '/api': {
        target: 'http://85.239.55.179:8000',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/api/, '') // Убираем /api из пути
      }
    },
    port: 5173, // Порт по умолчанию для Vite
    host: true, // Доступен на всех сетевых интерфейсах
  },
  publicDir: 'public',
  build: {
    rollupOptions: {
      input: {
        main: 'index.html'
      }
    }
  }
})
