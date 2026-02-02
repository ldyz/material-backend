<template>
  <div class="profile-page">
    <!-- 用户信息卡片 -->
    <div class="user-card">
      <div class="user-avatar">
        <van-icon name="user-circle-o" size="64" color="#666" />
      </div>
      <div class="user-info">
        <h3 class="username">{{ userInfo?.username || '用户' }}</h3>
        <p class="user-roles">{{ userRolesText }}</p>
      </div>
    </div>

    <!-- 功能菜单 -->
    <van-cell-group inset title="功能">
      <van-cell
        title="我的权限"
        icon="shield-o"
        is-link
        @click="showPermissions = true"
      />
    </van-cell-group>

    <!-- 系统设置 -->
    <van-cell-group inset title="设置">
      <van-cell
        title="清除缓存"
        icon="delete-o"
        is-link
        @click="clearCache"
      />
      <van-cell
        title="关于"
        icon="info-o"
        :value="systemVersion"
        @click="showAbout"
      />
    </van-cell-group>

    <!-- 退出登录 -->
    <div class="logout-section">
      <van-button
        block
        type="danger"
        @click="handleLogout"
      >
        退出登录
      </van-button>
    </div>

    <!-- 权限弹窗 -->
    <van-popup
      v-model:show="showPermissions"
      position="bottom"
      round
      :style="{ height: '60%' }"
    >
      <div class="permissions-popup">
        <div class="popup-header">
          <h3>我的权限</h3>
          <van-icon
            name="cross"
            size="20"
            @click="showPermissions = false"
          />
        </div>
        <div class="permissions-list">
          <van-tag
            v-for="perm in displayPermissions"
            :key="perm.key"
            type="primary"
            size="medium"
            style="margin: 4px"
          >
            {{ perm.label }}
          </van-tag>
          <van-empty
            v-if="displayPermissions.length === 0"
            description="暂无权限"
          />
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { showDialog, showToast } from 'vant'
import { useAuth } from '@/composables/useAuth'
import { useUserStore } from '@/stores/user'
import { PERMISSIONS } from '@/utils/constants'

const { logout } = useAuth()
const userStore = useUserStore()

const showPermissions = ref(false)
const systemVersion = ref('1.0.0')

const userInfo = computed(() => userStore.userInfo)
const roles = computed(() => userStore.roles)
const permissions = computed(() => userStore.permissions)

// 用户角色文本
const userRolesText = computed(() => {
  const roleLabels = {
    admin: '管理员',
    project_manager: '项目经理',
    warehouse_manager: '仓库管理员',
    worker: '施工人员',
  }

  return roles.value
    .map(role => roleLabels[role] || role)
    .join('、') || '普通用户'
})

// 权限显示列表
const displayPermissions = computed(() => {
  const permissionLabels = {
    [PERMISSIONS.INBOUND_VIEW]: '查看入库',
    [PERMISSIONS.INBOUND_CREATE]: '创建入库',
    [PERMISSIONS.INBOUND_APPROVE]: '审批入库',
    [PERMISSIONS.INBOUND_DELETE]: '删除入库',
    [PERMISSIONS.REQUISITION_VIEW]: '查看出库',
    [PERMISSIONS.REQUISITION_CREATE]: '创建出库',
    [PERMISSIONS.REQUISITION_APPROVE]: '审批出库',
    [PERMISSIONS.REQUISITION_ISSUE]: '发料',
    [PERMISSIONS.REQUISITION_DELETE]: '删除出库',
    [PERMISSIONS.STOCK_VIEW]: '查看库存',
    [PERMISSIONS.STOCK_ADJUST]: '调整库存',
    [PERMISSIONS.CONSTRUCTION_LOG_VIEW]: '查看日志',
    [PERMISSIONS.CONSTRUCTION_LOG_CREATE]: '创建日志',
    [PERMISSIONS.AI_ANALYZE]: 'AI 分析',
  }

  return permissions.value
    .filter(p => permissionLabels[p])
    .map(p => ({
      key: p,
      label: permissionLabels[p],
    }))
})

// 退出登录
function handleLogout() {
  showDialog({
    title: '提示',
    message: '确定要退出登录吗？',
    showCancelButton: true,
  })
    .then(() => {
      logout()
    })
    .catch(() => {
      // 取消
    })
}

// 清除缓存
function clearCache() {
  showDialog({
    title: '提示',
    message: '确定要清除缓存吗？',
    showCancelButton: true,
  })
    .then(() => {
      localStorage.clear()
      showToast('缓存已清除')
    })
    .catch(() => {
      // 取消
    })
}

// 关于
function showAbout() {
  showDialog({
    title: '关于',
    message: `${import.meta.env.VITE_APP_TITLE || '材料管理系统'}\n版本: ${systemVersion.value}`,
  })
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding: 0 0 60px 0;
}

.user-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 30px 20px;
  display: flex;
  align-items: center;
  color: white;
}

.user-avatar {
  margin-right: 16px;
}

.user-info {
  flex: 1;
}

.username {
  font-size: 20px;
  font-weight: bold;
  margin: 0 0 8px 0;
}

.user-roles {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
}

.logout-section {
  padding: 16px;
  margin-top: 16px;
}

.permissions-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.popup-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #ebedf0;
}

.popup-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: bold;
}

.permissions-list {
  flex: 1;
  padding: 16px 20px;
  overflow-y: auto;
}
</style>
