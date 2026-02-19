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
        manualChunks(id) {
          // Vendor chunks
          if (id.includes('node_modules')) {
            // Element Plus
            if (id.includes('element-plus')) {
              return 'element-plus'
            }

            // Vue ecosystem
            if (id.includes('vue/') || id.includes('@vue') || id.includes('pinia') || id.includes('vue-router')) {
              return 'vue-vendor'
            }

            // Editor (Quill)
            if (id.includes('quill')) {
              return 'editor'
            }

            // Charts
            if (id.includes('chart.js') || id.includes('vue-chartjs')) {
              return 'charts'
            }

            // PDF export libraries - split separately
            if (id.includes('jspdf')) {
              return 'pdf-export'
            }

            if (id.includes('html2canvas')) {
              return 'pdf-export'
            }

            // Date utilities
            if (id.includes('date-fns')) {
              return 'date-utils'
            }

            // Virtual scroller
            if (id.includes('vue-virtual-scroller')) {
              return 'virtual-scroller'
            }

            // Other utilities
            if (id.includes('lodash-es')) {
              return 'lodash'
            }

            if (id.includes('axios')) {
              return 'http'
            }

            // Node_modules polyfills
            if (id.includes('core-js') || id.includes('regenerator-runtime')) {
              return 'polyfills'
            }

            return 'vendor'
          }
        },
        chunkFileNames: 'assets/[name]-[hash].js',
        entryFileNames: 'assets/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash].[ext]'
      },
    },
  },
})
