import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import { VantResolver } from '@vant/auto-import-resolver'
import { fileURLToPath, URL } from 'node:url'

// 检测是否在 Capacitor 原生环境
const isCapacitor = process.env.CAPACITOR_BUILD === 'true'

export default defineConfig({
  // Capacitor 环境使用相对路径，Web环境使用 /mobile/ 前缀
  base: isCapacitor ? '' : '/mobile/',
  plugins: [
    vue(),
    Components({
      resolvers: [VantResolver()],
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8088',
        changeOrigin: true,
        secure: false,
      }
    }
  }
})
