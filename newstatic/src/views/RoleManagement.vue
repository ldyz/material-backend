<template>
  <div class="role-management">
    <el-card shadow="never">
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
        v-loading="loading"
        :data="roleList"
        border
        stripe
        style="width: 100%; margin-top: 20px"
      >
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
    </el-card>

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
      <div style="margin-bottom: 20px">
        <el-text>正在为角色 <strong>{{ currentRole?.name }}</strong> 配置权限</el-text>
      </div>
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
import { roleApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Edit,
  Delete,
  Setting
} from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import TableToolbar from '@/components/common/TableToolbar.vue'

const authStore = useAuthStore()

// 表格数据
const loading = ref(false)
const roleList = ref([])

// 角色对话框
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

// 权限配置对话框
const permissionDialogVisible = ref(false)
const permissionDialogLoading = ref(false)
const permissionTreeRef = ref(null)
const currentRole = ref(null)
const checkedPermissions = ref([])

const treeProps = {
  children: 'children',
  label: 'label'
}

// 权限树
const permissionTree = [
  {
    id: 'user',
    label: '用户管理',
    children: [
      { id: 'user_view', label: '查看用户' },
      { id: 'user_create', label: '创建用户' },
      { id: 'user_edit', label: '编辑用户' },
      { id: 'user_delete', label: '删除用户' }
    ]
  },
  {
    id: 'role',
    label: '角色管理',
    children: [
      { id: 'role_view', label: '查看角色' },
      { id: 'role_create', label: '创建角色' },
      { id: 'role_edit', label: '编辑角色' },
      { id: 'role_delete', label: '删除角色' },
      { id: 'role_assign_permissions', label: '分配权限' }
    ]
  },
  {
    id: 'project',
    label: '项目管理',
    children: [
      { id: 'project_view', label: '查看项目' },
      { id: 'project_create', label: '创建项目' },
      { id: 'project_edit', label: '编辑项目' },
      { id: 'project_delete', label: '删除项目' }
    ]
  },
  {
    id: 'material',
    label: '物资管理',
    children: [
      { id: 'material_view', label: '查看物资' },
      { id: 'material_create', label: '创建物资' },
      { id: 'material_edit', label: '编辑物资' },
      { id: 'material_delete', label: '删除物资' },
      { id: 'material_import', label: '导入物资' }
    ]
  },
  {
    id: 'material_plan',
    label: '物资计划',
    children: [
      { id: 'material_plan_view', label: '查看物资计划' },
      { id: 'material_plan_create', label: '创建物资计划' },
      { id: 'material_plan_edit', label: '编辑物资计划' },
      { id: 'material_plan_delete', label: '删除物资计划' },
      { id: 'material_plan_approve', label: '审核物资计划' }
    ]
  },
  {
    id: 'stock',
    label: '库存管理',
    children: [
      { id: 'stock_view', label: '查看库存' },
      { id: 'stock_create', label: '创建库存' },
      { id: 'stock_edit', label: '编辑' },
      { id: 'stock_delete', label: '删除' },
      { id: 'stock_in', label: '入库' },
      { id: 'stock_out', label: '出库' },
      { id: 'stock_export', label: '导出' },
      { id: 'stock_alerts', label: '库存预警' }
    ]
  },
  {
    id: 'stocklog',
    label: '库存日志',
    children: [
      { id: 'stocklog_view', label: '查看日志' },
      { id: 'stocklog_delete', label: '删除日志' }
    ]
  },
  {
    id: 'inbound',
    label: '入库管理',
    children: [
      { id: 'inbound_view', label: '查看入库单' },
      { id: 'inbound_create', label: '创建入库单' },
      { id: 'inbound_edit', label: '编辑入库单' },
      { id: 'inbound_delete', label: '删除入库单' },
      { id: 'inbound_approve', label: '审核入库单' },
      { id: 'inbound_export', label: '导出' }
    ]
  },
  {
    id: 'requisition',
    label: '出库管理',
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
    id: 'construction_log',
    label: '施工日志',
    children: [
      { id: 'construction_log_view', label: '查看日志' },
      { id: 'construction_log_create', label: '创建日志' },
      { id: 'construction_log_edit', label: '编辑日志' },
      { id: 'construction_log_delete', label: '删除日志' },
      { id: 'construction_log_export', label: '导出日志' }
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
    id: 'audit',
    label: '审计日志',
    children: [
      { id: 'audit_view', label: '查看审计日志' }
    ]
  },
  {
    id: 'ai_agent',
    label: 'AI 智能体',
    children: [
      { id: 'ai_agent_view', label: '查看 AI' },
      { id: 'ai_agent_query', label: 'AI 查询' },
      { id: 'ai_agent_operate', label: 'AI 操作' },
      { id: 'ai_agent_workflow', label: 'AI 工作流' },
      { id: 'ai_agent_logs', label: 'AI 日志' }
    ]
  },
  {
    id: 'system',
    label: '系统管理',
    children: [
      { id: 'system_log', label: '查看系统日志' },
      { id: 'system_backup', label: '数据备份' },
      { id: 'system_config', label: '系统配置' },
      { id: 'system_statistics', label: '系统统计' },
      { id: 'system_activities', label: '系统动态' }
    ]
  },
  {
    id: 'workflow',
    label: '工作流管理',
    children: [
      { id: 'workflow_view', label: '查看工作流' },
      { id: 'workflow_create', label: '创建工作流' },
      { id: 'workflow_edit', label: '编辑工作流' },
      { id: 'workflow_delete', label: '删除工作流' },
      { id: 'workflow_activate', label: '激活工作流' },
      { id: 'workflow_instance_view', label: '查看实例' },
      { id: 'workflow_instance_resubmit', label: '重新提交' },
      { id: 'workflow_task_view', label: '查看任务' },
      { id: 'workflow_task_approve', label: '审批任务' },
      { id: 'workflow_task_reject', label: '拒绝任务' },
      { id: 'workflow_task_delegate', label: '委派任务' },
      { id: 'workflow_log_view', label: '查看流程日志' }
    ]
  },
  {
    id: 'appointment',
    label: '施工预约',
    children: [
      { id: 'appointment_view', label: '查看预约单' },
      { id: 'appointment_create', label: '创建预约单' },
      { id: 'appointment_edit', label: '编辑预约单' },
      { id: 'appointment_delete', label: '删除预约单' },
      { id: 'appointment_submit', label: '提交审批' },
      { id: 'appointment_approve', label: '审批预约单' },
      { id: 'appointment_assign', label: '分配作业人员' },
      { id: 'appointment_execute', label: '执行作业' },
      { id: 'appointment_cancel', label: '取消预约单' },
      { id: 'appointment_export', label: '导出数据' }
    ]
  }
]

// 获取角色列表
const fetchRoles = async () => {
  loading.value = true
  try {
    const { data } = await roleApi.getList({ pageSize: 1000 })
    roleList.value = data || []
  } catch (error) {
    console.error('获取角色列表失败:', error)
    ElMessage.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

// 添加角色
const handleAddRole = () => {
  resetRoleForm()
  isEditRole.value = false
  roleDialogVisible.value = true
}

// 编辑角色
const handleEditRole = (row) => {
  Object.assign(roleForm, {
    id: row.id,
    name: row.name,
    description: row.description
  })
  isEditRole.value = true
  roleDialogVisible.value = true
}

// 删除角色
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
    if (error !== 'cancel') {
      console.error('操作失败:', error)
    }
  })
}

// 提交角色
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
    ElMessage.error(error?.message || '操作失败，请重试')
  } finally {
    roleDialogLoading.value = false
  }
}

// 配置权限
const handleRolePermissions = (row) => {
  currentRole.value = row
  checkedPermissions.value = row.permissions || []
  permissionDialogVisible.value = true
}

// 提交权限配置
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
    ElMessage.error(error?.message || '权限配置失败，请重试')
  } finally {
    permissionDialogLoading.value = false
  }
}

// 重置角色表单
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

// 初始化
onMounted(() => {
  fetchRoles()
})
</script>

<style scoped>
.role-management {
  padding: 0;
}
</style>
