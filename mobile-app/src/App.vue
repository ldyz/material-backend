<template>
  <router-view />
</template>

<script setup>
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  min-height: 100vh;
  background-color: #f7f8fa;
  position: relative;
}

/* 修复 Vant Toast 样式 */
:deep(.van-toast) {
  background-color: rgba(0, 0, 0, 0.8) !important;
  color: #fff !important;
}

:deep(.van-toast__text) {
  color: #fff !important;
}

:deep(.van-toast--success) {
  background-color: rgba(25, 137, 250, 0.9) !important;
}

:deep(.van-toast--fail) {
  background-color: rgba(255, 73, 73, 0.9) !important;
}

:deep(.van-toast--loading) {
  background-color: rgba(0, 0, 0, 0.8) !important;
}

:deep(.van-loading__text) {
  color: #fff !important;
}

:deep(.van-loading__spinner) {
  color: #fff !important;
}
</style>
