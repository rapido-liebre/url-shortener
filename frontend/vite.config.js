import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/links': 'http://localhost:8080',
      '/u': 'http://localhost:8080'
    }
  }
})
