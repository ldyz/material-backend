<template>
  <view class="login-page">
    <view class="login-header">
      <image class="logo" src="/static/images/logo.png" mode="aspectFit" />
      <text class="title">项目管理系统</text>
      <text class="subtitle">微信小程序版</text>
    </view>

    <!-- 微信一键登录 -->
    <view class="login-content" v-if="loginMode === 'wechat'">
      <button class="wechat-btn" @click="handleWechatLogin">
        <text class="btn-text">微信一键登录</text>
      </button>
      <view class="switch-mode" @click="switchToAccount">
        <text>使用账号密码登录</text>
      </view>
    </view>

    <!-- 账号密码登录 -->
    <view class="login-content" v-else>
      <view class="form-item">
        <input
          class="input"
          type="text"
          v-model="form.username"
          placeholder="请输入用户名"
          placeholder-class="placeholder"
        />
      </view>
      <view class="form-item">
        <input
          class="input"
          type="password"
          v-model="form.password"
          placeholder="请输入密码"
          placeholder-class="placeholder"
        />
      </view>
      <button class="login-btn" :loading="loading" @click="handleLogin">
        <text class="btn-text">登录</text>
      </button>
      <view class="switch-mode" @click="switchToWechat">
        <text>使用微信登录</text>
      </view>
    </view>

    <!-- 绑定账号弹窗 -->
    <view class="bind-modal" v-if="showBindModal">
      <view class="modal-content">
        <view class="modal-header">
          <text class="modal-title">绑定账号</text>
        </view>
        <view class="modal-body">
          <text class="modal-tips">该微信未绑定账号，请输入已有账号密码进行绑定</text>
          <view class="form-item">
            <input
              class="input"
              type="text"
              v-model="bindForm.username"
              placeholder="请输入用户名"
              placeholder-class="placeholder"
            />
          </view>
          <view class="form-item">
            <input
              class="input"
              type="password"
              v-model="bindForm.password"
              placeholder="请输入密码"
              placeholder-class="placeholder"
            />
          </view>
        </view>
        <view class="modal-footer">
          <button class="cancel-btn" @click="showBindModal = false">取消</button>
          <button class="confirm-btn" :loading="bindLoading" @click="handleBindWechat">绑定</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import { useAuthStore } from '@/stores/auth'

export default {
  data() {
    return {
      loginMode: 'wechat', // wechat | account
      loading: false,
      bindLoading: false,
      showBindModal: false,
      openid: '',
      form: {
        username: '',
        password: ''
      },
      bindForm: {
        username: '',
        password: ''
      }
    }
  },
  methods: {
    switchToAccount() {
      this.loginMode = 'account'
    },
    switchToWechat() {
      this.loginMode = 'wechat'
    },
    async handleWechatLogin() {
      const authStore = useAuthStore()
      this.loading = true

      try {
        await authStore.wechatLogin()
        this.loginSuccess()
      } catch (error) {
        if (error.needBind) {
          this.openid = error.openid
          this.showBindModal = true
        } else {
          uni.showToast({
            title: error.message || '登录失败',
            icon: 'none'
          })
        }
      } finally {
        this.loading = false
      }
    },
    async handleLogin() {
      if (!this.form.username) {
        uni.showToast({ title: '请输入用户名', icon: 'none' })
        return
      }
      if (!this.form.password) {
        uni.showToast({ title: '请输入密码', icon: 'none' })
        return
      }

      const authStore = useAuthStore()
      this.loading = true

      try {
        await authStore.login(this.form.username, this.form.password)
        this.loginSuccess()
      } catch (error) {
        uni.showToast({
          title: error.message || '登录失败',
          icon: 'none'
        })
      } finally {
        this.loading = false
      }
    },
    async handleBindWechat() {
      if (!this.bindForm.username) {
        uni.showToast({ title: '请输入用户名', icon: 'none' })
        return
      }
      if (!this.bindForm.password) {
        uni.showToast({ title: '请输入密码', icon: 'none' })
        return
      }

      const authStore = useAuthStore()
      this.bindLoading = true

      try {
        await authStore.bindWechat(this.bindForm.username, this.bindForm.password)
        this.showBindModal = false
        this.loginSuccess()
      } catch (error) {
        uni.showToast({
          title: error.message || '绑定失败',
          icon: 'none'
        })
      } finally {
        this.bindLoading = false
      }
    },
    loginSuccess() {
      uni.showToast({
        title: '登录成功',
        icon: 'success'
      })
      setTimeout(() => {
        uni.switchTab({
          url: '/pages/index/index'
        })
      }, 1000)
    }
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #1989fa 0%, #07c160 100%);
  padding: 0 48rpx;
  display: flex;
  flex-direction: column;
}

.login-header {
  padding-top: 160rpx;
  text-align: center;
  margin-bottom: 80rpx;
}

.logo {
  width: 160rpx;
  height: 160rpx;
  margin-bottom: 32rpx;
}

.title {
  display: block;
  font-size: 48rpx;
  font-weight: 600;
  color: #ffffff;
  margin-bottom: 16rpx;
}

.subtitle {
  display: block;
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
}

.login-content {
  background-color: #ffffff;
  border-radius: 24rpx;
  padding: 48rpx 32rpx;
}

.form-item {
  margin-bottom: 32rpx;
}

.input {
  width: 100%;
  height: 96rpx;
  background-color: #f7f8fa;
  border-radius: 16rpx;
  padding: 0 32rpx;
  font-size: 32rpx;
}

.placeholder {
  color: #c8c9cc;
}

.wechat-btn, .login-btn {
  width: 100%;
  height: 96rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 16rpx;
}

.wechat-btn {
  background-color: #07c160;
}

.login-btn {
  background-color: #1989fa;
}

.btn-text {
  color: #ffffff;
  font-size: 32rpx;
  font-weight: 500;
}

.switch-mode {
  text-align: center;
  margin-top: 32rpx;
  font-size: 28rpx;
  color: #969799;
}

.bind-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  width: 600rpx;
  background-color: #ffffff;
  border-radius: 24rpx;
  overflow: hidden;
}

.modal-header {
  padding: 32rpx;
  text-align: center;
  border-bottom: 1rpx solid #ebedf0;
}

.modal-title {
  font-size: 36rpx;
  font-weight: 500;
  color: #323233;
}

.modal-body {
  padding: 32rpx;
}

.modal-tips {
  display: block;
  font-size: 28rpx;
  color: #969799;
  margin-bottom: 32rpx;
}

.modal-footer {
  display: flex;
  border-top: 1rpx solid #ebedf0;
}

.cancel-btn, .confirm-btn {
  flex: 1;
  height: 100rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  border-radius: 0;
  background-color: transparent;
}

.cancel-btn {
  color: #646566;
  border-right: 1rpx solid #ebedf0;
}

.confirm-btn {
  color: #1989fa;
}
</style>
