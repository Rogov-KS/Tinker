import { defineConfig } from 'vite'

const backendTarget =
  process.env.VITE_BACKEND_URL ?? 'http://localhost:3000'

export default defineConfig({
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: backendTarget,
        changeOrigin: true,
      },
    },
  },
  build: {
    minify: false,
    sourcemap: true,
  },
})
