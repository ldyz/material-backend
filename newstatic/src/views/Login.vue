<template>
  <div class="login-container">
    <div class="login-background">
      <div class="background-shape shape-1"></div>
      <div class="background-shape shape-2"></div>
      <div class="background-shape shape-3"></div>
    </div>
    <div class="login-card">
      <div class="login-header">
        <div class="login-icon">
          <el-icon :size="64"><User /></el-icon>
        </div>
        <h2 class="login-title">用户登录</h2>
        <p class="login-subtitle">欢迎回来，请登录您的账号</p>
      </div>
      <div class="login-body">
        <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" @submit.prevent="handleLogin">
          <el-form-item prop="username">
            <div class="login-input-group">
              <div class="input-icon">
                <el-icon><User /></el-icon>
              </div>
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名"
                autocomplete="username"
                size="large"
              />
            </div>
          </el-form-item>
          <el-form-item prop="password">
            <div class="login-input-group">
              <div class="input-icon">
                <el-icon><Lock /></el-icon>
              </div>
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                autocomplete="current-password"
                size="large"
                show-password
                @keyup.enter="handleLogin"
              />
            </div>
          </el-form-item>
          <div class="login-actions">
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              class="login-button"
              @click="handleLogin"
            >
              <span class="button-text">登 录</span>
              <el-icon class="button-icon"><ArrowRight /></el-icon>
            </el-button>
          </div>
        </el-form>
        <div class="login-tips">
          <el-icon><InfoFilled /></el-icon>
          <span>默认管理员账号：admin / admin</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { User, Lock, ArrowRight, InfoFilled } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
    loading.value = true

    await authStore.login({
      username: loginForm.username,
      password: loginForm.password
    })

    // 登录成功后跳转
    const redirect = route.query.redirect || '/dashboard'
    router.push(redirect)
  } catch (error) {
    console.error('登录失败:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  z-index: 9999;
  overflow: hidden;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.background-shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  animation: float 20s infinite ease-in-out;
}

.shape-1 {
  width: 400px;
  height: 400px;
  background: white;
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

.shape-2 {
  width: 300px;
  height: 300px;
  background: white;
  bottom: -50px;
  right: -50px;
  animation-delay: 5s;
}

.shape-3 {
  width: 200px;
  height: 200px;
  background: white;
  top: 50%;
  right: 10%;
  animation-delay: 10s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  33% {
    transform: translate(30px, -30px) rotate(120deg);
  }
  66% {
    transform: translate(-20px, 20px) rotate(240deg);
  }
}

.login-card {
  position: relative;
  background: white;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  width: 90%;
  max-width: 420px;
  padding: 0;
  overflow: hidden;
  animation: slideUp 0.5s ease-out;
  z-index: 1;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40px 30px 30px;
  text-align: center;
  color: white;
}

.login-icon {
  font-size: 64px;
  margin-bottom: 15px;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 8px 0;
  letter-spacing: 1px;
}

.login-subtitle {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
  font-weight: 400;
}

.login-body {
  padding: 40px 30px 30px;
}

.login-input-group {
  display: flex;
  align-items: center;
  background: #f5f7fa;
  border-radius: 12px;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  padding: 0 15px;
}

.login-input-group:focus-within {
  background: white;
  border-color: #667eea;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
}

.login-input-group :deep(.el-input__wrapper) {
  background: transparent;
  box-shadow: none;
  padding: 16px 15px;
}

.login-input-group :deep(.el-input__inner) {
  font-size: 15px;
  color: #333;
}

.input-icon {
  color: #667eea;
  font-size: 18px;
  margin-right: 10px;
}

.login-actions {
  margin-top: 35px;
}

.login-button {
  width: 100%;
  height: 52px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  color: white;
  font-size: 16px;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
  letter-spacing: 1px;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.login-button:active {
  transform: translateY(0);
}

.button-text {
  flex: 1;
  text-align: center;
}

.button-icon {
  transition: transform 0.3s ease;
}

.login-button:hover .button-icon {
  transform: translateX(5px);
}

.login-tips {
  margin-top: 25px;
  padding: 12px 16px;
  background: #f0f4ff;
  border-radius: 8px;
  border-left: 4px solid #667eea;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: #555;
}

.login-tips .el-icon {
  color: #667eea;
  font-size: 14px;
}

@media (max-width: 768px) {
  .login-card {
    max-width: 90%;
  }

  .login-header {
    padding: 30px 20px 20px;
  }

  .login-title {
    font-size: 24px;
  }

  .login-body {
    padding: 30px 20px 20px;
  }

  .background-shape {
    display: none;
  }
}

@media (max-width: 480px) {
  .login-card {
    width: 95%;
  }

  .login-icon {
    font-size: 48px;
  }

  .login-title {
    font-size: 20px;
  }
}
</style>
