<template>
  <div class="system-container">
    <el-card shadow="never">
      <!-- 权限刷新提示区 -->
      <el-alert
        title="权限更新提示"
        type="info"
        :closable="false"
        style="margin-bottom: 16px"
      >
        <template #default>
          <div style="display: flex; align-items: center; justify-content: space-between;">
            <span>如果刚更新了角色权限但页面没有生效，请点击下方按钮刷新权限</span>
            <el-button
              type="primary"
              size="small"
              :icon="Refresh"
              @click="handleRefreshPermissions"
              :loading="refreshingPermissions"
            >
              刷新我的权限
            </el-button>
          </div>
        </template>
      </el-alert>

      <el-tabs v-model="activeTab" type="border-card">
        <!-- 数据备份 -->
        <el-tab-pane label="数据备份" name="backup">
          <div class="backup-section">
            <el-alert
              title="数据备份"
              type="info"
              :closable="false"
              style="margin-bottom: 20px"
            >
              定期备份数据库数据，确保数据安全。建议每天备份一次。
            </el-alert>

            <el-form :model="backupForm" label-width="120px" style="max-width: 600px">
              <el-form-item label="备份名称">
                <el-input
                  v-model="backupForm.name"
                  placeholder="自动生成或手动输入"
                  maxlength="100"
                />
              </el-form-item>
              <el-form-item label="备份说明">
                <el-input
                  v-model="backupForm.description"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入备份说明"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="backupLoading" @click="handleCreateBackup">
                  立即备份
                </el-button>
              </el-form-item>
            </el-form>

            <el-divider />

            <h3>备份历史</h3>
            <el-table
              v-loading="backupListLoading"
              :data="backupList"
              border
              stripe
              style="width: 100%; margin-top: 20px"
            >
              <!-- 序号列已移除 -->
              <el-table-column prop="filename" label="备份名称" min-width="200" show-overflow-tooltip />
              <el-table-column prop="description" label="说明" min-width="200" show-overflow-tooltip />
              <el-table-column prop="size" label="文件大小" width="120">
                <template #default="scope">
                  {{ scope.row.size !== null && scope.row.size !== undefined ? formatFileSize(scope.row.size) : '-' }}
                </template>
              </el-table-column>
              <el-table-column prop="created_at" label="备份时间" width="160" />
              <el-table-column label="操作" width="320" fixed="right">
                <template #default="scope">
                  <el-button
                    type="primary"
                    size="small"
                    :icon="Upload"
                    @click="handleRestoreBackup(scope.row)"
                    v-if="authStore.hasPermission('backup_restore')"
                  >
                    恢复
                  </el-button>
                  <el-button
                    type="success"
                    size="small"
                    :icon="Download"
                    @click="handleDownloadBackup(scope.row)"
                  >
                    下载
                  </el-button>
                  <el-button
                    type="danger"
                    size="small"
                    :icon="Delete"
                    @click="handleDeleteBackup(scope.row)"
                    v-if="authStore.hasPermission('backup_delete')"
                  >
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <!-- 系统设置 -->
        <el-tab-pane label="系统设置" name="settings">
          <div class="settings-section">
            <el-form :model="systemSettings" label-width="150px" style="max-width: 800px">
              <el-divider content-position="left">基本信息</el-divider>
              <el-form-item label="系统名称">
                <el-input v-model="systemSettings.system_name" maxlength="50" />
              </el-form-item>
              <el-form-item label="系统简称">
                <el-input v-model="systemSettings.system_short_name" maxlength="20" />
              </el-form-item>

              <el-divider content-position="left">安全设置</el-divider>
              <el-form-item label="密码最小长度">
                <el-input-number v-model="systemSettings.password_min_length" :min="6" :max="20" />
              </el-form-item>
              <el-form-item label="Token 有效期（小时）">
                <el-input-number v-model="systemSettings.token_expiry" :min="1" :max="168" />
                <span style="margin-left: 10px; color: #909399">默认 72 小时</span>
              </el-form-item>
              <el-form-item label="启用验证码">
                <el-switch v-model="systemSettings.enable_captcha" />
              </el-form-item>

              <el-divider content-position="left">上传设置</el-divider>
              <el-form-item label="上传目录" prop="upload_directory">
                <el-input
                  v-model="systemSettings.upload_directory"
                  placeholder="文件上传目录路径"
                  style="width: 400px"
                />
                <div class="form-tip">相对于服务器根目录的路径，例如: static/uploads</div>
              </el-form-item>
              <el-form-item label="最大上传大小（MB）" prop="max_file_size">
                <el-input-number v-model="systemSettings.max_file_size" :min="1" :max="100" />
              </el-form-item>
              <el-form-item label="最多上传数量" prop="max_upload_count">
                <el-input-number v-model="systemSettings.max_upload_count" :min="1" :max="50" />
                <div class="form-tip">单次最多可上传的文件数量</div>
              </el-form-item>
              <el-form-item label="允许的文件类型" prop="allowed_file_types">
                <el-input
                  v-model="systemSettings.allowed_file_types"
                  placeholder="例如: jpg,jpeg,png,gif,bmp,webp,svg"
                  style="width: 400px"
                />
                <div class="form-tip">用逗号分隔的文件扩展名列表</div>
              </el-form-item>

              <el-form-item>
                <el-button type="primary" :loading="settingsLoading" @click="handleSaveSettings">
                  保存设置
                </el-button>
                <el-button @click="fetchSystemSettings">重置</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { systemApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Refresh,
  Delete,
  Download,
  Upload
} from '@element-plus/icons-vue'

const authStore = useAuthStore()

// 当前激活的标签
const activeTab = ref('backup')

// 权限刷新状态
const refreshingPermissions = ref(false)

// 刷新当前用户的权限
const handleRefreshPermissions = async () => {
  try {
    refreshingPermissions.value = true
    await authStore.refreshUserInfo()
    ElMessage.success('权限刷新成功，页面将自动重新加载')
    setTimeout(() => {
      location.reload()
    }, 1500)
  } catch (error) {
    console.error('刷新权限失败:', error)
    ElMessage.error('刷新权限失败，请重新登录')
  } finally {
    refreshingPermissions.value = false
  }
}

// ============ 数据备份 ============
const backupLoading = ref(false)
const backupListLoading = ref(false)
const backupList = ref([])
const backupForm = reactive({
  name: '',
  description: ''
})

// ============ 系统设置 ============
const settingsLoading = ref(false)
const systemSettings = reactive({
  system_name: '材料管理系统',
  system_short_name: 'MMS',
  password_min_length: 6,
  token_expiry: 72,
  enable_captcha: false,
  upload_directory: 'static/uploads',
  max_file_size: 5,
  max_upload_count: 10,
  allowed_file_types: 'jpg,jpeg,png,gif,bmp,webp,svg'
})

// ============ 数据备份方法 ============
// 适配统一响应格式
const fetchBackups = async () => {
  backupListLoading.value = true
  try {
    const { data } = await systemApi.getBackups()
    backupList.value = data || []
  } catch (error) {
    console.error('获取备份列表失败:', error)
  } finally {
    backupListLoading.value = false
  }
}

const handleCreateBackup = async () => {
  try {
    backupLoading.value = true
    const data = {
      name: backupForm.name || `backup_${Date.now()}`,
      description: backupForm.description
    }
    await systemApi.createBackup(data)
    ElMessage.success('备份创建成功')
    backupForm.name = ''
    backupForm.description = ''
    fetchBackups()
  } catch (error) {
    console.error('创建备份失败:', error)
  } finally {
    backupLoading.value = false
  }
}

const handleDownloadBackup = (row) => {
  const a = document.createElement('a')
  a.href = row.filepath || row.file_url
  a.download = row.filename
  a.click()
}

const handleRestoreBackup = (row) => {
  ElMessageBox.confirm(
    `确定要从备份"${row.filename}"恢复数据吗？\n\n警告：此操作将清空当前数据库并恢复备份数据，请谨慎操作！`,
    '确认恢复数据库',
    {
      confirmButtonText: '确定恢复',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await systemApi.restoreBackup({ backup_name: row.filename, confirm: true })
      ElMessage.success('数据恢复成功，请刷新页面')
      setTimeout(() => {
        window.location.reload()
      }, 1500)
    } catch (error) {
      console.error('恢复失败:', error)
      ElMessage.error(error?.message || '恢复失败，请重试')
    }
  }).catch(() => {
    // 用户取消操作，不需要提示
  })
}

const handleDeleteBackup = (row) => {
  ElMessageBox.confirm(
    `确定要删除备份"${row.filename}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await systemApi.deleteBackup(row.id)
      ElMessage.success('删除成功')
      fetchBackups()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error(error?.message || '删除失败，请重试')
    }
  }).catch((error) => {
    // 用户取消操作，不需要提示
    if (error !== 'cancel') {
      console.error('操作失败:', error)
    }
  })
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

// ============ 系统设置方法 ============
// 适配统一响应格式
const fetchSystemSettings = async () => {
  try {
    const { data } = await systemApi.getSettings()
    Object.assign(systemSettings, data)
  } catch (error) {
    console.error('获取系统设置失败:', error)
  }
}

const handleSaveSettings = async () => {
  try {
    settingsLoading.value = true
    await systemApi.saveSettings(systemSettings)
    ElMessage.success('设置保存成功')
    authStore.setSystemName(systemSettings.system_name)
  } catch (error) {
    console.error('保存设置失败:', error)
  } finally {
    settingsLoading.value = false
  }
}

onMounted(() => {
  fetchBackups()
  fetchSystemSettings()
})
</script>

<style scoped>
.system-container {
  padding: 0;
}

.backup-section,
.settings-section {
  padding: 20px 0;
}

:deep(.el-tabs__content) {
  padding: 20px;
}

.text-gray {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}
</style>
