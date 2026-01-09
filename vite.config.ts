import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'

export default defineConfig({
  plugins: [solid()],
  server: {
    proxy: {
      '/api': {
        target: 'http://188.225.24.208:8000',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/api/, '') // Убираем /api из пути
      }
    }
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
