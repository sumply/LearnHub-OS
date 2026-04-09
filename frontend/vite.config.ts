import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import solid from 'vite-plugin-solid'

const frontendRoot = path.dirname(fileURLToPath(import.meta.url))

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, frontendRoot, '')
  const apiTarget = env.VITE_API_URL || 'http://185.152.92.245:8000'

  return {
    root: frontendRoot,
    envDir: frontendRoot,
    plugins: [solid()],
    server: {
      proxy: {
        '/api': {
          target: apiTarget,
          changeOrigin: true,
          secure: false,
          rewrite: (p) => p.replace(/^\/api/, ''),
        },
      },
      port: 5173,
      host: true,
    },
    publicDir: 'public',
    build: {
      rollupOptions: {
        input: {
          main: 'index.html',
        },
      },
    },
  }
})
