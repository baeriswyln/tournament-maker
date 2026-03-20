import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, './src')
    }
  },
  server: {
    watch : {
      usePolling: true,
      interval: 1000
    },
    port: 3000,
    proxy: {
      '/api': {
        target: 'https://tournament-maker.baeriswyl.dev',
        changeOrigin: true
      }
    }
  }
})
