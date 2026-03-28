import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8088',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true,
    assetsDir: 'assets',
    sourcemap: false,
    chunkSizeWarningLimit: 1000, // Increase limit to 1000KB
    rollupOptions: {
      output: {
        manualChunks: {
          // Element Plus (large UI library)
          'element-plus': ['element-plus'],

          // Vue ecosystem (core framework)
          'vue-vendor': ['vue', 'vue-router', 'pinia'],

          // Charts library
          'charts': ['chart.js', 'vue-chartjs'],

          // PDF export libraries
          'pdf-export': ['jspdf', 'html2canvas'],

          // Editor
          'editor': ['quill', '@vueup/vue-quill'],

          // Other major libraries
          'vendor': ['axios', 'lodash-es', 'date-fns', 'vue-virtual-scroller'],
        },
        chunkFileNames: 'assets/[name]-[hash].js',
        entryFileNames: 'assets/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash].[ext]'
      },
    },
  },
})
