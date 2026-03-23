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

/* 修复 Vant Toast 样式 - 使用最强选择器确保样式生效 */
.van-toast,
body > .van-toast,
#app > .van-toast,
.van-toast__container,
.van-notify {
  background-color: rgba(0, 0, 0, 0.85) !important;
  color: #fff !important;
}

.van-toast .van-toast__text,
.van-toast__text,
.van-notify__text {
  color: #fff !important;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3) !important;
}

/* 成功提示 */
.van-toast--success,
.van-toast.van-toast--success,
.van-toast__container .van-toast--success {
  background-color: rgba(25, 137, 250, 0.9) !important;
  color: #fff !important;
}

.van-toast--success .van-toast__text {
  color: #fff !important;
}

/* 失败提示 */
.van-toast--fail,
.van-toast.van-toast--fail,
.van-toast__container .van-toast--fail {
  background-color: rgba(255, 73, 73, 0.9) !important;
  color: #fff !important;
}

.van-toast--fail .van-toast__text {
  color: #fff !important;
}

/* 加载提示 */
.van-toast--loading,
.van-toast.van-toast--loading,
.van-toast__container .van-toast--loading {
  background-color: rgba(0, 0, 0, 0.85) !important;
  color: #fff !important;
}

.van-toast--loading .van-toast__text,
.van-loading__text {
  color: #fff !important;
}

.van-loading__spinner,
.van-loading__spinner .circular {
  color: #fff !important;
  fill: #fff !important;
}

/* 确保 loading 图标可见 */
.van-loading__spinner circle {
  stroke: #fff !important;
}
</style>
