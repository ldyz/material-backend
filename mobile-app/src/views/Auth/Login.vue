<template>
  <div class="login-page">
    <div class="login-header">
      <h1 class="app-title">{{ appTitle }}</h1>
      <p class="app-subtitle">材料管理系统</p>
    </div>

    <van-form @submit="onSubmit" class="login-form">
      <van-cell-group inset>
        <van-field
          v-model="username"
          name="username"
          label="用户名"
          placeholder="请输入用户名"
          :rules="[{ required: true, message: '请输入用户名' }]"
          clearable
          autocomplete="username"
        />
        <van-field
          v-model="password"
          type="password"
          name="password"
          label="密码"
          placeholder="请输入密码"
          :rules="[{ required: true, message: '请输入密码' }]"
          clearable
          autocomplete="current-password"
        />
      </van-cell-group>

      <div class="login-button">
        <van-button
          round
          block
          type="primary"
          native-type="submit"
          :loading="loading"
        >
          登录
        </van-button>
      </div>
    </van-form>

    <div class="login-footer">
      <p>© {{ currentYear }} {{ appTitle }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const route = useRoute()
const { login } = useAuth()

const username = ref('')
const password = ref('')
const loading = ref(false)

const appTitle = computed(() => import.meta.env.VITE_APP_TITLE || '材料管理')
const currentYear = computed(() => new Date().getFullYear())

async function onSubmit(values) {
  try {
    loading.value = true
    await login(values.username, values.password)

    // 跳转到重定向页面或首页
    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (error) {
    showToast({
      type: 'fail',
      message: error.message || '登录失败',
    })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-header {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
}

.app-title {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 10px;
}

.app-subtitle {
  font-size: 16px;
  opacity: 0.9;
}

.login-form {
  margin-bottom: 30px;
}

.login-button {
  margin-top: 30px;
  padding: 0 16px;
}

.login-footer {
  text-align: center;
  color: rgba(255, 255, 255, 0.8);
  font-size: 12px;
  padding: 20px 0;
}

.login-footer p {
  margin: 0;
}
</style>
