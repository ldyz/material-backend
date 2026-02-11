<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside :width="sidebarWidth" class="layout-aside">
      <div class="logo-container">
        <el-icon class="logo-icon"><Box /></el-icon>
        <span v-if="!appStore.sidebarCollapsed" class="logo-text">{{ appStore.systemName }}</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        :collapse="appStore.sidebarCollapsed"
        :unique-opened="true"
        router
        class="layout-menu"
      >
        <template v-for="menu in visibleMenus" :key="menu.path">
          <el-menu-item :index="menu.path" v-if="!menu.children">
            <el-icon>
              <component :is="menu.icon" />
            </el-icon>
            <template #title>{{ menu.title }}</template>
          </el-menu-item>
          <el-sub-menu :index="menu.path" v-else>
            <template #title>
              <el-icon>
                <component :is="menu.icon" />
              </el-icon>
              <span>{{ menu.title }}</span>
            </template>
            <el-menu-item
              v-for="child in menu.children"
              :key="child.path"
              :index="child.path"
            >
              <el-icon>
                <component :is="child.icon" />
              </el-icon>
              <template #title>{{ child.title }}</template>
            </el-menu-item>
          </el-sub-menu>
        </template>
      </el-menu>
    </el-aside>

    <el-container class="layout-main">
      <!-- 顶部导航栏 -->
      <el-header class="layout-header">
        <div class="header-left">
          <el-button
            :icon="appStore.sidebarCollapsed ? 'Expand' : 'Fold'"
            circle
            @click="appStore.toggleSidebar"
          />
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentRoute.meta.title">
              {{ currentRoute.meta.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <!-- Notification Bell -->
          <NotificationBell class="header-notification" />

          <el-dropdown @command="handleCommand">
            <div class="user-info">
              <el-avatar
                :size="40"
                :src="authStore.user?.avatar || undefined"
                :style="{ backgroundColor: getAvatarColor() }"
              >
                {{ authStore.displayName?.charAt(0) || '?' }}
              </el-avatar>
              <span class="username">{{ authStore.displayName }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="uploadAvatar">
                  <el-icon><Picture /></el-icon>
                  上传头像
                </el-dropdown-item>
                <el-dropdown-item command="resetPassword">
                  <el-icon><Key /></el-icon>
                  修改密码
                </el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <!-- 隐藏的文件上传input -->
          <input
            ref="avatarInputRef"
            type="file"
            accept="image/*"
            style="display: none"
            @change="handleAvatarChange"
          />

          <!-- 头像裁剪对话框 -->
          <AvatarCropperDialog
            v-model="showCropperDialog"
            :image-file="selectedAvatarFile"
            @success="handleAvatarSuccess"
          />
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main class="layout-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { useNotificationStore } from '@/stores/notificationStore'
import { createVisibleMenus } from '@/utils/permissions'
import {
  Box,
  UserFilled,
  Key,
  SwitchButton,
  Expand,
  Fold,
  Odometer,
  Management,
  Document,
  DocumentCopy,
  List,
  ShoppingCart,
  Goods,
  TrendCharts,
  Setting,
  DataAnalysis,
  Grid,
  Connection,
  Clock,
  Calendar,
  Picture
} from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import NotificationBell from '@/components/Notification/NotificationBell.vue'
import AvatarCropperDialog from '@/components/common/AvatarCropperDialog.vue'
import { authApi } from '@/api'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()
const notificationStore = useNotificationStore()

// 当前激活的菜单
const activeMenu = computed(() => route.path)

// 当前路由
const currentRoute = computed(() => route)

// 侧边栏宽度
const sidebarWidth = computed(() => {
  return appStore.sidebarCollapsed ? '64px' : '200px'
})

// 菜单配置（按业务流程和使用频率排序）
const menuConfig = [
  {
    path: '/dashboard',
    title: '仪表板',
    icon: TrendCharts,
    permissions: [] // 所有人可见
  },
  {
    path: '/projects',
    title: '项目管理',
    icon: Odometer,
    permissions: ['project_view'] // 项目查看权限
  },
  {
    path: '/progress',
    title: '进度管理',
    icon: DataAnalysis,
    permissions: ['progress_view'] // 进度查看权限
  },
  {
    path: '/construction-log',
    title: '施工日志',
    icon: Document,
    permissions: ['constructionlog_view', 'constructionlog_create', 'constructionlog_edit', 'constructionlog_delete'] // 施工日志权限
  },
  {
    path: '/materials-group',
    title: '物资管理',
    icon: Goods,
    permissions: ['material_view', 'material_plan_view'], // 物资或计划查看权限
    children: [
      {
        path: '/material-plans',
        title: '物资计划',
        icon: List,
        permissions: ['material_plan_view'] // 物资计划查看权限
      },
      {
        path: '/materials',
        title: '物资浏览',
        icon: Goods,
        permissions: ['material_view'] // 物资查看权限
      },
      {
        path: '/material-categories',
        title: '物资分类',
        icon: Grid,
        permissions: ['material_view'] // 使用物资查看权限
      }
    ]
  },
  {
    path: '/stock-group',
    title: '库存管理',
    icon: Management,
    permissions: ['stock_view', 'inbound_view', 'requisition_view'], // 库存、入库或出库权限
    children: [
      {
        path: '/stock',
        title: '库存浏览',
        icon: Management,
        permissions: ['stock_view'] // 库存查看权限
      },
      {
        path: '/inbound',
        title: '入库管理',
        icon: ShoppingCart,
        permissions: ['inbound_view'] // 入库单查看权限
      },
      {
        path: '/requisitions',
        title: '出库管理',
        icon: DocumentCopy,
        permissions: ['requisition_view'] // 出库单查看权限
      }
    ]
  },
  {
    path: '/appointments',
    title: '施工预约',
    icon: Calendar,
    permissions: ['appointment_view'] // 预约查看权限
  },
  {
    path: '/workflows',
    title: '工作流管理',
    icon: Connection,
    permissions: ['system_config'] // 使用系统配置权限
  },
  {
    path: '/operation-logs',
    title: '操作日志',
    icon: Clock,
    permissions: ['audit_view'] // 操作日志查看权限
  },
  {
    path: '/system',
    title: '系统管理',
    icon: Setting,
    permissions: [
      'user_view',
      'role_view',
      'system_log',
      'system_backup',
      'system_config',
      'system_report'
    ],
    children: [
      {
        path: '/system/users',
        title: '用户管理',
        icon: UserFilled,
        permissions: ['user_view']
      },
      {
        path: '/system/roles',
        title: '角色管理',
        icon: Management,
        permissions: ['role_view']
      },
      {
        path: '/system',
        title: '系统设置',
        icon: Setting,
        permissions: ['system_log', 'system_backup', 'system_config', 'system_report']
      }
    ]
  }
]

// 可见菜单 - 使用 ref 而非 computed，避免响应式循环
// 在组件挂载时一次性计算，不响应权限变化
const visibleMenus = ref([])

// 初始化菜单
const initMenus = () => {
  visibleMenus.value = createVisibleMenus(menuConfig, authStore)
}

// 组件挂载时初始化菜单
onMounted(() => {
  initMenus()
  handleResize()
  window.addEventListener('resize', handleResize)

  // 初始化通知 WebSocket
  if (authStore.isAuthenticated) {
    notificationStore.fetchUnreadCount()
    notificationStore.initWebSocket()
  }
})

// 组件卸载时清理
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  // 关闭通知 WebSocket
  notificationStore.closeWebSocket()
})

// 处理用户下拉菜单命令
const handleCommand = async (command) => {
  if (command === 'logout') {
    try {
      await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      await authStore.logout()
      router.push('/login')
    } catch (error) {
      // 用户取消
    }
  } else if (command === 'resetPassword') {
    router.push('/reset-password')
  } else if (command === 'uploadAvatar') {
    avatarInputRef.value?.click()
  }
}

// 头像上传
const avatarInputRef = ref(null)
const showCropperDialog = ref(false)
const selectedAvatarFile = ref(null)

const handleAvatarChange = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return

  // 验证文件大小（5MB - 裁剪前可以大一些）
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error('图片文件大小不能超过5MB')
    return
  }

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('只支持图片格式的文件')
    return
  }

  // 打开裁剪对话框
  selectedAvatarFile.value = file
  showCropperDialog.value = true

  // 清空 input
  if (avatarInputRef.value) {
    avatarInputRef.value.value = ''
  }
}

const handleAvatarSuccess = async () => {
  // 刷新用户信息
  await authStore.refreshUserInfo()
  selectedAvatarFile.value = null
}

const getAvatarColor = () => {
  const colors = ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399']
  const userId = authStore.user?.id || 0
  return colors[userId % colors.length]
}


// 响应式处理
const handleResize = () => {
  const width = window.innerWidth
  if (width < 768) {
    appStore.setDevice('mobile')
  } else {
    appStore.setDevice('desktop')
  }
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.layout-aside {
  background: #fff;
  border-right: 1px solid #e6e6e6;
  transition: width 0.3s;
  overflow: hidden;
}

.logo-container {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
}

.logo-icon {
  font-size: 24px;
  color: #409eff;
  margin-right: 10px;
}

.logo-text {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.layout-menu {
  border-right: none;
  height: calc(100vh - 60px);
  overflow-y: auto;
}

.layout-main {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.layout-header {
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.header-notification {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background 0.3s;
}

.user-info:hover {
  background: #f5f7fa;
}

.username {
  font-size: 14px;
  color: #333;
}

.layout-content {
  background: #f5f7fa;
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
