import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import { VantResolver } from '@vant/auto-import-resolver'
import { fileURLToPath, URL } from 'node:url'

// 检测是否在 Capacitor 原生环境
const isCapacitor = process.env.CAPACITOR_BUILD === 'true'

export default defineConfig({
  // 生产环境使用 /mobile/ 路径，开发环境使用相对路径
  base: process.env.NODE_ENV === 'production' ? '/mobile/' : './',
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
    // 监听所有接口，但优化配置避免崩溃
    host: '0.0.0.0',
    port: 5173,
    strictPort: false,
    proxy: {
      '/api': {
        target: 'http://localhost:8088',
        changeOrigin: true,
        secure: false,
        // 添加超时配置
        timeout: 30000,
        // 添加代理错误处理
        configure: (proxy, _options) => {
          proxy.on('error', (err, _req, res) => {
            console.log('Proxy error:', err.message);
          });
          proxy.on('proxyReq', (proxyReq, _req, _res) => {
            console.log('Proxy request:', proxyReq.method, proxyReq.path);
          });
        }
      }
    },
    // 优化文件监听，减少文件系统压力
    watch: {
      usePolling: false,
      interval: 1000
    },
    // 限制并发连接数
    middlewareMode: false
  }
})
