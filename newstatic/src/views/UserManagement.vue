<template>
  <div class="user-management">
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

      <TableToolbar>
        <template #left>
          <el-input
            v-model="searchKeyword"
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
        v-loading="loading"
        :data="userList"
        border
        stripe
        style="width: 100%; margin-top: 20px"
      >
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
        <el-table-column prop="created_at" label="创建时间" width="160" />
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
              v-if="authStore.hasPermission('user_edit')"
            >
              重置密码
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        style="margin-top: 20px; display: flex; justify-content: flex-end"
      />
    </el-card>

    <!-- 用户对话框 -->
    <Dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      :loading="dialogLoading"
      @confirm="handleSubmit"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            maxlength="50"
            :disabled="isEdit"
          />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            maxlength="50"
            show-password
          />
        </el-form-item>
        <el-form-item label="姓名" prop="full_name">
          <el-input
            v-model="form.full_name"
            placeholder="请输入姓名"
            maxlength="50"
          />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="form.email"
            placeholder="请输入邮箱"
            maxlength="100"
          />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select
            v-model="form.role"
            placeholder="请选择角色"
            style="width: 100%"
          >
            <el-option
              v-for="role in allRoles"
              :key="role.id"
              :label="role.name"
              :value="role.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="is_active">
          <el-radio-group v-model="form.is_active">
            <el-radio :label="true">启用</el-radio>
            <el-radio :label="false">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { userApi, roleApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  Edit
} from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import TableToolbar from '@/components/common/TableToolbar.vue'

const authStore = useAuthStore()

// 搜索关键词
const searchKeyword = ref('')

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

// 表格数据
const loading = ref(false)
const userList = ref([])
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 所有角色列表（用于下拉框）
const allRoles = ref([])

// 对话框
const dialogVisible = ref(false)
const isEdit = ref(false)
const dialogTitle = computed(() => isEdit.value ? '编辑用户' : '添加用户')
const dialogLoading = ref(false)
const formRef = ref(null)

// 表单数据
const form = reactive({
  id: null,
  username: '',
  password: '',
  full_name: '',
  email: '',
  role: '',
  is_active: true
})

const formRules = {
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
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 获取用户列表
const fetchUsers = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchKeyword.value || undefined
    }
    const { data, pagination: pag } = await userApi.getList(params)
    userList.value = data || []
    pagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 获取角色列表
const fetchRoles = async () => {
  try {
    const { data } = await roleApi.getList({ pageSize: 1000 })
    allRoles.value = data || []
  } catch (error) {
    console.error('获取角色列表失败:', error)
  }
}

// 添加用户
const handleAddUser = () => {
  resetForm()
  isEdit.value = false
  dialogVisible.value = true
  fetchRoles()
}

// 编辑用户
const handleEditUser = (row) => {
  Object.assign(form, {
    id: row.id,
    username: row.username,
    full_name: row.full_name,
    email: row.email,
    role: row.role,
    is_active: row.is_active
  })
  isEdit.value = true
  dialogVisible.value = true
  fetchRoles()
}

// 重置密码
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
    if (error !== 'cancel') {
      console.error('操作失败:', error)
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    dialogLoading.value = true

    const data = {
      username: form.username,
      full_name: form.full_name,
      email: form.email,
      role: form.role,
      is_active: form.is_active
    }

    if (form.password) {
      data.password = form.password
    }

    if (isEdit.value) {
      await userApi.update(form.id, data)
      ElMessage.success('更新成功')
    } else {
      await userApi.create(data)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchUsers()
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error(error?.message || '操作失败，请重试')
  } finally {
    dialogLoading.value = false
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    id: null,
    username: '',
    password: '',
    full_name: '',
    email: '',
    role: '',
    is_active: true
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 分页处理
const handlePageChange = (page) => {
  pagination.page = page
  fetchUsers()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchUsers()
}

// 初始化
onMounted(() => {
  fetchUsers()
  fetchRoles()
})
</script>

<style scoped>
.user-management {
  padding: 0;
}
</style>
