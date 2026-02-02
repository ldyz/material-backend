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
        <!-- 用户管理 -->
        <el-tab-pane label="用户管理" name="users">
          <TableToolbar>
            <template #left>
              <el-input
                v-model="userSearch.keyword"
                placeholder="搜索用户名、姓名"
                clearable
                style="width: 250px"
                @keyup.enter="fetchUsers"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-button type="primary" :icon="Search" @click="fetchUsers">
                搜索
              </el-button>
            </template>
            <template #right>
              <el-button
                type="primary"
                :icon="Plus"
                @click="handleAddUser"
                v-if="authStore.hasPermission('user_create')"
              >
                添加用户
              </el-button>
            </template>
          </TableToolbar>

          <el-table
            v-loading="userLoading"
            :data="userList"
            border
            stripe
            style="width: 100%; margin-top: 20px"
          >
            <!-- 序号列已移除 -->
            <el-table-column prop="username" label="用户名" width="150" />
            <el-table-column prop="full_name" label="姓名" width="120" />
            <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
            <el-table-column prop="role" label="角色" width="120" />
            <el-table-column prop="is_active" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.is_active ? 'success' : 'danger'" size="small">
                  {{ scope.row.is_active ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <el-button
                  type="primary"
                  size="small"
                  :icon="Edit"
                  @click="handleEditUser(scope.row)"
                  v-if="authStore.hasPermission('user_edit')"
                >
                  编辑
                </el-button>
                <el-button
                  type="info"
                  size="small"
                  @click="handleResetPassword(scope.row)"
                  v-if="authStore.hasPermission('user_reset_password')"
                >
                  重置密码
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="userPagination.page"
            v-model:page-size="userPagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="userPagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleUserSizeChange"
            @current-change="handleUserPageChange"
            style="margin-top: 20px; display: flex; justify-content: flex-end"
          />
        </el-tab-pane>

        <!-- 角色管理 -->
        <el-tab-pane label="角色管理" name="roles">
          <TableToolbar>
            <template #right>
              <el-button
                type="primary"
                :icon="Plus"
                @click="handleAddRole"
                v-if="authStore.hasPermission('role_create')"
              >
                添加角色
              </el-button>
            </template>
          </TableToolbar>

          <el-table
            v-loading="roleLoading"
            :data="roleList"
            border
            stripe
            style="width: 100%; margin-top: 20px"
          >
            <!-- 序号列已移除 -->
            <el-table-column prop="name" label="角色名称" width="180" />
            <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
            <el-table-column prop="created_at" label="创建时间" width="160" />
            <el-table-column label="操作" width="280" fixed="right">
              <template #default="scope">
                <el-button
                  type="primary"
                  size="small"
                  :icon="Edit"
                  @click="handleEditRole(scope.row)"
                  v-if="authStore.hasPermission('role_edit')"
                >
                  编辑
                </el-button>
                <el-button
                  type="success"
                  size="small"
                  :icon="Setting"
                  @click="handleRolePermissions(scope.row)"
                  v-if="authStore.hasPermission('role_assign_permissions')"
                >
                  权限配置
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  :icon="Delete"
                  @click="handleDeleteRole(scope.row)"
                  v-if="authStore.hasPermission('role_delete')"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 操作日志 -->
        <el-tab-pane label="操作日志" name="logs">
          <TableToolbar>
            <template #left>
              <el-input
                v-model="logSearch.keyword"
                placeholder="搜索操作内容"
                clearable
                style="width: 250px"
                @keyup.enter="fetchLogs"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-select
                v-model="logSearch.module"
                placeholder="模块"
                clearable
                style="width: 150px"
                @change="fetchLogs"
              >
                <el-option label="全部" value="" />
                <el-option label="物资管理" value="material" />
                <el-option label="库存管理" value="stock" />
                <el-option label="出库单" value="requisition" />
                <el-option label="入库单" value="inbound" />
                <el-option label="项目管理" value="project" />
                <el-option label="系统管理" value="system" />
              </el-select>
              <el-date-picker
                v-model="logSearch.date_range"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="YYYY-MM-DD"
                style="width: 240px"
                @change="fetchLogs"
              />
              <el-button type="primary" :icon="Search" @click="fetchLogs">
                搜索
              </el-button>
              <el-button :icon="Refresh" @click="handleResetLogSearch">重置</el-button>
            </template>
          </TableToolbar>

          <el-table
            v-loading="logLoading"
            :data="logList"
            border
            stripe
            style="width: 100%; margin-top: 20px"
          >
            <!-- 序号列已移除 -->
            <el-table-column prop="user" label="操作人" width="120" />
            <el-table-column prop="module" label="模块" width="120">
              <template #default="scope">
                {{ getModuleText(scope.row.module) }}
              </template>
            </el-table-column>
            <el-table-column prop="message" label="操作" width="120" show-overflow-tooltip />
            <el-table-column prop="ip_address" label="IP地址" width="130" />
            <el-table-column prop="created_at" label="操作时间" width="160" />
          </el-table>

          <el-pagination
            v-model:current-page="logPagination.page"
            v-model:page-size="logPagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="logPagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleLogSizeChange"
            @current-change="handleLogPageChange"
            style="margin-top: 20px; display: flex; justify-content: flex-end"
          />
        </el-tab-pane>

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

    <!-- 用户对话框 -->
    <Dialog
      v-model="userDialogVisible"
      :title="userDialogTitle"
      width="600px"
      :loading="userDialogLoading"
      @confirm="handleSubmitUser"
    >
      <el-form
        ref="userFormRef"
        :model="userForm"
        :rules="userFormRules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="userForm.username"
            placeholder="请输入用户名"
            maxlength="50"
            :disabled="isEditUser"
          />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEditUser">
          <el-input
            v-model="userForm.password"
            type="password"
            placeholder="请输入密码"
            maxlength="50"
            show-password
          />
        </el-form-item>
        <el-form-item label="姓名" prop="full_name">
          <el-input
            v-model="userForm.full_name"
            placeholder="请输入姓名"
            maxlength="50"
          />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="userForm.email"
            placeholder="请输入邮箱"
            maxlength="100"
          />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select
            v-model="userForm.role"
            placeholder="请选择角色"
            style="width: 100%"
          >
            <el-option
              v-for="role in roleList"
              :key="role.id"
              :label="role.name"
              :value="role.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="is_active">
          <el-radio-group v-model="userForm.is_active">
            <el-radio :label="true">启用</el-radio>
            <el-radio :label="false">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </Dialog>

    <!-- 角色对话框 -->
    <Dialog
      v-model="roleDialogVisible"
      :title="roleDialogTitle"
      width="600px"
      :loading="roleDialogLoading"
      @confirm="handleSubmitRole"
    >
      <el-form
        ref="roleFormRef"
        :model="roleForm"
        :rules="roleFormRules"
        label-width="100px"
      >
        <el-form-item label="角色名称" prop="name">
          <el-input
            v-model="roleForm.name"
            placeholder="请输入角色名称"
            maxlength="50"
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="roleForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入角色描述"
            maxlength="200"
          />
        </el-form-item>
      </el-form>
    </Dialog>

    <!-- 权限配置对话框 -->
    <Dialog
      v-model="permissionDialogVisible"
      title="权限配置"
      width="800px"
      :loading="permissionDialogLoading"
      @confirm="handleSubmitPermissions"
    >
      <el-tree
        ref="permissionTreeRef"
        :data="permissionTree"
        :props="treeProps"
        show-checkbox
        node-key="id"
        :default-checked-keys="checkedPermissions"
        style="margin-bottom: 20px"
      />
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { userApi, roleApi, systemApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  Edit,
  Delete,
  Download,
  Upload,
  Setting
} from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import TableToolbar from '@/components/common/TableToolbar.vue'

const authStore = useAuthStore()

// 当前激活的标签
const activeTab = ref('users')

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

// ============ 用户管理 ============
const userLoading = ref(false)
const userList = ref([])
const userPagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})
const userSearch = reactive({
  keyword: ''
})

const userDialogVisible = ref(false)
const isEditUser = ref(false)
const userDialogTitle = computed(() => isEditUser.value ? '编辑用户' : '添加用户')
const userDialogLoading = ref(false)
const userFormRef = ref(null)

const userForm = reactive({
  id: null,
  username: '',
  password: '',
  full_name: '',
  email: '',
  role: '',
  is_active: true
})

const userFormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]{3,20}$/, message: '用户名格式不正确（3-20位字母数字下划线）', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  full_name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// ============ 角色管理 ============
const roleLoading = ref(false)
const roleList = ref([])

const roleDialogVisible = ref(false)
const isEditRole = ref(false)
const roleDialogTitle = computed(() => isEditRole.value ? '编辑角色' : '添加角色')
const roleDialogLoading = ref(false)
const roleFormRef = ref(null)

const roleForm = reactive({
  id: null,
  name: '',
  description: ''
})

const roleFormRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' }
  ]
}

// ============ 权限配置 ============
const permissionDialogVisible = ref(false)
const permissionDialogLoading = ref(false)
const permissionTreeRef = ref(null)
const currentRole = ref(null)
const checkedPermissions = ref([])

const treeProps = {
  children: 'children',
  label: 'label'
}

const permissionTree = [
  {
    id: 'material',
    label: '物资管理',
    children: [
      { id: 'material_view', label: '查看物资' },
      { id: 'material_create', label: '创建物资' },
      { id: 'material_edit', label: '编辑物资' },
      { id: 'material_delete', label: '删除物资' },
      { id: 'material_import', label: '导入物资' },
      { id: 'material_export', label: '导出物资' }
    ]
  },
  {
    id: 'stock',
    label: '库存管理',
    children: [
      { id: 'stock_view', label: '查看库存' },
      { id: 'stock_in', label: '入库' },
      { id: 'stock_out', label: '出库' },
      { id: 'stock_edit', label: '编辑' },
      { id: 'stock_delete', label: '删除' },
      { id: 'stocklog_view', label: '查看日志' },
      { id: 'stock_export', label: '导出' }
    ]
  },
  {
    id: 'requisition',
    label: '出库单管理',
    children: [
      { id: 'requisition_view', label: '查看出库单' },
      { id: 'requisition_create', label: '创建出库单' },
      { id: 'requisition_edit', label: '编辑出库单' },
      { id: 'requisition_delete', label: '删除出库单' },
      { id: 'requisition_approve', label: '审核出库单' },
      { id: 'requisition_issue', label: '发货' },
      { id: 'requisition_export', label: '导出' }
    ]
  },
  {
    id: 'inbound',
    label: '入库单管理',
    children: [
      { id: 'inbound_view', label: '查收入库单' },
      { id: 'inbound_create', label: '创建入库单' },
      { id: 'inbound_edit', label: '编辑入库单' },
      { id: 'inbound_delete', label: '删除入库单' },
      { id: 'inbound_approve', label: '审核入库单' },
      { id: 'inbound_export', label: '导出' }
    ]
  },
  {
    id: 'project',
    label: '项目管理',
    children: [
      { id: 'project_view', label: '查看项目' },
      { id: 'project_create', label: '创建项目' },
      { id: 'project_edit', label: '编辑项目' },
      { id: 'project_delete', label: '删除项目' },
      { id: 'project_export', label: '导出' }
    ]
  },
  {
    id: 'construction_log',
    label: '施工日志',
    children: [
      { id: 'construction_log_view', label: '查看日志' },
      { id: 'construction_log_create', label: '创建日志' },
      { id: 'construction_log_edit', label: '编辑日志' },
      { id: 'construction_log_delete', label: '删除日志' },
      { id: 'construction_log_export', label: '导出' }
    ]
  },
  {
    id: 'progress',
    label: '进度管理',
    children: [
      { id: 'progress_view', label: '查看进度' },
      { id: 'progress_create', label: '创建任务' },
      { id: 'progress_edit', label: '编辑任务' },
      { id: 'progress_delete', label: '删除任务' },
      { id: 'progress_export', label: '导出' }
    ]
  },
  {
    id: 'system',
    label: '系统管理',
    children: [
      { id: 'user_view', label: '查看用户' },
      { id: 'user_create', label: '创建用户' },
      { id: 'user_edit', label: '编辑用户' },
      { id: 'user_delete', label: '删除用户' },
      { id: 'user_reset_password', label: '重置密码' },
      { id: 'role_view', label: '查看角色' },
      { id: 'role_create', label: '创建角色' },
      { id: 'role_edit', label: '编辑角色' },
      { id: 'role_delete', label: '删除角色' },
      { id: 'role_assign_permissions', label: '分配权限' },
      { id: 'log_view', label: '查看日志' },
      { id: 'backup_create', label: '创建备份' },
      { id: 'backup_delete', label: '删除备份' },
      { id: 'settings_edit', label: '系统设置' }
    ]
  }
]

// ============ 操作日志 ============
const logLoading = ref(false)
const logList = ref([])
const logPagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})
const logSearch = reactive({
  keyword: '',
  module: '',
  date_range: []
})

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

// ============ 用户管理方法 ============
// 适配统一响应格式
const fetchUsers = async () => {
  userLoading.value = true
  try {
    const params = {
      page: userPagination.page,
      page_size: userPagination.pageSize,
      keyword: userSearch.keyword || undefined
    }
    const { data, pagination: pag } = await userApi.getList(params)
    userList.value = data || []
    userPagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取用户列表失败:', error)
  } finally {
    userLoading.value = false
  }
}

// 用户管理分页处理
const handleUserPageChange = (page) => {
  userPagination.page = page
  fetchUsers()
}

const handleUserSizeChange = (size) => {
  userPagination.pageSize = size
  userPagination.page = 1
  fetchUsers()
}

const handleAddUser = () => {
  resetUserForm()
  isEditUser.value = false
  userDialogVisible.value = true
  fetchRoles()
}

const handleEditUser = (row) => {
  Object.assign(userForm, {
    id: row.id,
    username: row.username,
    full_name: row.full_name,
    email: row.email,
    role: row.role,
    is_active: row.is_active
  })
  isEditUser.value = true
  userDialogVisible.value = true
  fetchRoles()
}

const handleSubmitUser = async () => {
  if (!userFormRef.value) return

  try {
    await userFormRef.value.validate()
    userDialogLoading.value = true

    const data = {
      username: userForm.username,
      full_name: userForm.full_name,
      email: userForm.email,
      role: userForm.role,
      is_active: userForm.is_active
    }

    if (isEditUser.value) {
      await userApi.update(userForm.id, data)
      ElMessage.success('更新成功')
    } else {
      data.password = userForm.password
      await userApi.create(data)
      ElMessage.success('创建成功')
    }

    userDialogVisible.value = false
    fetchUsers()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    userDialogLoading.value = false
  }
}

const handleResetPassword = (row) => {
  ElMessageBox.prompt('请输入新密码', `重置 ${row.full_name} 的密码`, {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /^.{6,}$/,
    inputErrorMessage: '密码至少6位'
  }).then(async ({ value }) => {
    try {
      await userApi.resetPassword(row.id, { password: value })
      ElMessage.success('密码重置成功')
    } catch (error) {
      console.error('重置密码失败:', error)
      ElMessage.error(error?.message || '密码重置失败，请重试')
    }
  }).catch((error) => {
    // 用户取消操作，不需要提示
    if (error !== 'cancel') {
      console.error('操作失败:', error)
    }
  })
}

const resetUserForm = () => {
  Object.assign(userForm, {
    id: null,
    username: '',
    password: '',
    full_name: '',
    email: '',
    role: '',
    is_active: true
  })
  if (userFormRef.value) {
    userFormRef.value.clearValidate()
  }
}

// ============ 角色管理方法 ============
// 适配统一响应格式
const fetchRoles = async () => {
  roleLoading.value = true
  try {
    const { data } = await roleApi.getList({ pageSize: 1000 })
    roleList.value = data || []
  } catch (error) {
    console.error('获取角色列表失败:', error)
  } finally {
    roleLoading.value = false
  }
}

const handleAddRole = () => {
  resetRoleForm()
  isEditRole.value = false
  roleDialogVisible.value = true
}

const handleEditRole = (row) => {
  Object.assign(roleForm, {
    id: row.id,
    name: row.name,
    description: row.description
  })
  isEditRole.value = true
  roleDialogVisible.value = true
}

const handleSubmitRole = async () => {
  if (!roleFormRef.value) return

  try {
    await roleFormRef.value.validate()
    roleDialogLoading.value = true

    const data = {
      name: roleForm.name,
      description: roleForm.description
    }

    if (isEditRole.value) {
      await roleApi.update(roleForm.id, data)
      ElMessage.success('更新成功')
    } else {
      await roleApi.create(data)
      ElMessage.success('创建成功')
    }

    roleDialogVisible.value = false
    fetchRoles()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    roleDialogLoading.value = false
  }
}

const handleDeleteRole = (row) => {
  ElMessageBox.confirm(
    `确定要删除角色"${row.name}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await roleApi.delete(row.id)
      ElMessage.success('删除成功')
      fetchRoles()
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

const handleRolePermissions = (row) => {
  currentRole.value = row
  checkedPermissions.value = row.permissions || []
  permissionDialogVisible.value = true
}

const handleSubmitPermissions = async () => {
  if (!permissionTreeRef.value) return

  try {
    permissionDialogLoading.value = true
    const checkedKeys = permissionTreeRef.value.getCheckedKeys()
    const halfCheckedKeys = permissionTreeRef.value.getHalfCheckedKeys()
    const allPermissions = [...checkedKeys, ...halfCheckedKeys]

    await roleApi.assignPermissions(currentRole.value.id, {
      permissions: allPermissions
    })
    ElMessage.success('权限配置成功')

    // 检查当前登录用户是否使用了被修改的角色
    const currentUser = authStore.user
    const currentRoleIds = currentUser?.role_ids || []
    if (currentRoleIds.includes(currentRole.value.id)) {
      // 刷新当前用户的权限信息
      await authStore.refreshUserInfo()
      ElMessage.info('您的权限已更新，页面将自动刷新')
      setTimeout(() => {
        location.reload()
      }, 1500)
    }

    permissionDialogVisible.value = false
    fetchRoles()
  } catch (error) {
    console.error('权限配置失败:', error)
  } finally {
    permissionDialogLoading.value = false
  }
}

const resetRoleForm = () => {
  Object.assign(roleForm, {
    id: null,
    name: '',
    description: ''
  })
  if (roleFormRef.value) {
    roleFormRef.value.clearValidate()
  }
}

// ============ 操作日志方法 ============
// 适配统一响应格式
const fetchLogs = async () => {
  logLoading.value = true
  try {
    const params = {
      page: logPagination.page,
      page_size: logPagination.pageSize,
      keyword: logSearch.keyword || undefined,
      module: logSearch.module || undefined,
      start_date: logSearch.date_range?.[0] || undefined,
      end_date: logSearch.date_range?.[1] || undefined
    }
    const { data, pagination: pag } = await systemApi.getLogs(params)
    logList.value = data || []
    logPagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取操作日志失败:', error)
  } finally {
    logLoading.value = false
  }
}

// 操作日志分页处理
const handleLogPageChange = (page) => {
  logPagination.page = page
  fetchLogs()
}

const handleLogSizeChange = (size) => {
  logPagination.pageSize = size
  logPagination.page = 1
  fetchLogs()
}

const handleResetLogSearch = () => {
  logSearch.keyword = ''
  logSearch.module = ''
  logSearch.date_range = []
  fetchLogs()
}

const getModuleText = (module) => {
  const texts = {
    material: '物资管理',
    stock: '库存管理',
    requisition: '出库单',
    inbound: '入库单',
    project: '项目管理',
    system: '系统管理'
  }
  return texts[module] || module
}

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
  fetchUsers()
  fetchRoles()
  fetchLogs()
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
