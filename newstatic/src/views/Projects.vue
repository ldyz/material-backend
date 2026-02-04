<template>
  <div class="projects-container">
    <!-- 项目树抽屉 -->
    <el-drawer
      v-model="treeDrawerVisible"
      title="项目层级结构"
      size="400px"
      direction="ltr"
    >
      <ProjectTree
        v-if="selectedProjectId"
        :project-id="selectedProjectId"
        @add-child="handleAddChild"
        @view-schedule="handleViewSchedule"
      />
    </el-drawer>

    <el-card shadow="never">
      <!-- 工具栏 -->
      <TableToolbar>
        <template #left>
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索项目名称、编码"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select
            v-model="searchForm.status"
            placeholder="项目状态"
            clearable
            style="width: 120px"
          >
            <el-option label="全部" value="" />
            <el-option label="筹备中" value="planning" />
            <el-option label="进行中" value="active" />
            <el-option label="已完成" value="closed" />
            <el-option label="已暂停" value="on_hold" />
          </el-select>
          <el-select
            v-model="searchForm.parent_id"
            placeholder="筛选父项目"
            clearable
            filterable
            style="width: 180px"
            @change="handleSearch"
          >
            <el-option label="全部项目" :value="null" />
            <el-option
              v-for="project in parentProjectList"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </el-select>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </template>
        <template #right>
          <el-button
            type="primary"
            :icon="Plus"
            @click="handleAdd"
            v-if="authStore.hasPermission('project_create')"
          >
            创建项目
          </el-button>
          <el-button
            type="warning"
            :icon="Download"
            @click="handleExport"
            v-if="authStore.hasPermission('project_export')"
          >
            导出
          </el-button>
        </template>
      </TableToolbar>

      <!-- 树形表格 -->
      <el-table
        v-loading="loading"
        :data="treeData"
        row-key="id"
        border
        stripe
        style="width: 100%"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        :expand-row-keys="expandedKeys"
        @expand-change="handleExpandChange"
        default-expand-all
      >
        <el-table-column prop="code" label="项目编码" width="150" fixed="left">
          <template #default="scope">
            <div class="tree-node">
              <el-link type="primary" @click="handleView(scope.row)">
                {{ scope.row.code }}
              </el-link>
              <el-tag v-if="scope.row.level !== undefined" :type="getLevelTagType(scope.row.level)" size="small" style="margin-left: 8px;">
                {{ getLevelText(scope.row.level) }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="项目名称" min-width="200" show-overflow-tooltip fixed="left" />
        <el-table-column prop="location" label="项目地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="start_date" label="开始日期" width="110" />
        <el-table-column prop="end_date" label="结束日期" width="110" />
        <el-table-column prop="progress_percentage" label="进度" width="150">
          <template #default="scope">
            <el-progress
              :percentage="scope.row.progress_percentage || 0"
              :status="getProgressStatus(scope.row.progress_percentage, scope.row.status)"
              :stroke-width="12"
            />
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="330" fixed="right">
          <template #default="scope">
            <el-button
              type="info"
              size="small"
              :icon="View"
              @click="handleView(scope.row)"
            >
              查看
            </el-button>
            <el-button
              type="primary"
              size="small"
              :icon="Edit"
              @click="handleEdit(scope.row)"
              v-if="canEdit(scope.row)"
            >
              编辑
            </el-button>
            <el-button
              type="success"
              size="small"
              @click="handleViewTree(scope.row)"
            >
              <el-icon><Rank /></el-icon>
              层级
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="handleAddChildDirect(scope.row)"
            >
              <el-icon><Plus /></el-icon>
              子项目
            </el-button>
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              @click="handleDelete(scope.row)"
              v-if="canDelete(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        class="mt-20"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <Dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="700px"
      :loading="dialogLoading"
      @confirm="handleSubmit"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="项目编码" prop="code">
              <el-input
                v-model="formData.code"
                placeholder="请输入项目编码"
                maxlength="50"
                :disabled="isViewMode"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="项目名称" prop="name">
              <el-autocomplete
                v-model="formData.name"
                :fetch-suggestions="searchProjectsByName"
                placeholder="请输入项目名称（可搜索已有项目）"
                :trigger-on-focus="false"
                clearable
                style="width: 100%"
                :disabled="isViewMode"
                @select="handleSelectProject"
              >
                <template #default="{ item }">
                  <div class="project-suggestion">
                    <div class="suggestion-name">{{ item.name }}</div>
                    <div class="suggestion-info">
                      <el-tag size="small" type="info">{{ item.code }}</el-tag>
                      <span class="suggestion-hint">点击使用此项目</span>
                    </div>
                  </div>
                </template>
              </el-autocomplete>
              <div class="text-gray text-sm">
                <el-icon><InfoFilled /></el-icon>
                输入名称搜索已有项目，或直接输入新项目名称
              </div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="父项目" prop="parent_id">
          <el-select
            v-model="formData.parent_id"
            placeholder="请选择父项目（可选，不选则为主项目）"
            clearable
            filterable
            style="width: 100%"
            :disabled="isViewMode"
          >
            <el-option
              v-for="project in parentProjectList"
              :key="project.id"
              :label="`${project.name} (层级${getLevelText(project.level)})`"
              :value="project.id"
              :disabled="project.id === formData.id"
            />
          </el-select>
          <div class="text-gray text-sm">选择父项目后，此项目将成为子项目。最多支持4级项目。</div>
        </el-form-item>

        <el-form-item label="项目地址" prop="location">
          <el-input
            v-model="formData.location"
            placeholder="请输入项目地址"
            maxlength="200"
            :disabled="isViewMode"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始日期" prop="start_date">
              <el-date-picker
                v-model="formData.start_date"
                type="date"
                placeholder="选择开始日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                :disabled="isViewMode"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束日期" prop="end_date">
              <el-date-picker
                v-model="formData.end_date"
                type="date"
                placeholder="选择结束日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                :disabled="isViewMode"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="项目状态" prop="status">
              <el-select
                v-model="formData.status"
                placeholder="请选择状态"
                style="width: 100%"
                :disabled="isViewMode"
              >
                <el-option label="筹备中" value="planning" />
                <el-option label="进行中" value="active" />
                <el-option label="已完成" value="closed" />
                <el-option label="已暂停" value="on_hold" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="进度" prop="progress">
              <el-input-number
                v-model="formData.progress"
                :min="0"
                :max="100"
                :step="5"
                :disabled="isViewMode"
                style="width: 100%"
              />
              <span style="margin-left: 10px">%</span>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="项目描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="4"
            placeholder="请输入项目描述"
            maxlength="1000"
            :disabled="isViewMode"
          />
        </el-form-item>

        <el-form-item label="关联成员" prop="member_ids">
          <div style="display: flex; gap: 8px; align-items: center; flex-wrap: wrap">
            <el-tag
              v-for="userId in formData.member_ids"
              :key="userId"
              :closable="!isViewMode"
              @close="removeMember(userId)"
              style="margin: 2px"
            >
              {{ getUserNameById(userId) }}
            </el-tag>
            <el-button
              v-if="!isViewMode"
              type="primary"
              size="small"
              @click="openMemberSelector"
              icon="Plus"
            >
              选择成员
            </el-button>
          </div>
          <div class="text-gray text-sm mt-1">
            已选择 {{ formData.member_ids?.length || 0 }} 个成员，点击"选择成员"按钮进行修改
          </div>
        </el-form-item>
      </el-form>
    </Dialog>

    <!-- 成员选择对话框 -->
    <el-dialog
      v-model="memberSelectorVisible"
      title="选择项目成员"
      width="800px"
      :close-on-click-modal="false"
    >
      <div style="margin-bottom: 16px">
        <el-input
          v-model="memberSearchText"
          placeholder="搜索用户名、姓名或邮箱"
          prefix-icon="Search"
          clearable
          style="width: 300px"
        />
      </div>

      <div style="max-height: 500px; overflow-y: auto">
        <div v-for="group in filteredGroupedUsers" :key="group.label">
          <div style="
            background: #f5f7fa;
            padding: 10px 16px;
            font-weight: bold;
            color: #606266;
            border-radius: 4px;
            margin-bottom: 8px;
            position: sticky;
            top: 0;
            z-index: 10;
          ">
            <el-checkbox
              v-model="group.checked"
              :indeterminate="group.indeterminate"
              @change="handleGroupCheck(group)"
            >
              {{ group.label }} ({{ group.users.length }}人，已选{{ getGroupSelectedCount(group) }}人)
            </el-checkbox>
          </div>

          <el-checkbox-group v-model="selectedMembersTemp">
            <div
              v-for="user in group.users"
              :key="user.id"
              style="
                padding: 8px 16px;
                border-bottom: 1px solid #ebeef5;
                display: flex;
                align-items: center;
              "
            >
              <el-checkbox :label="user.id" style="flex: 1">
                <div style="display: flex; align-items: center; gap: 12px; flex: 1">
                  <div style="flex: 1">
                    <div style="font-weight: 500">{{ user.username }}</div>
                    <div style="font-size: 12px; color: #909399">
                      {{ user.full_name || '未设置姓名' }} | {{ user.email }}
                    </div>
                  </div>
                  <el-tag v-if="user.is_active" type="success" size="small">活跃</el-tag>
                  <el-tag v-else type="info" size="small">离线</el-tag>
                </div>
              </el-checkbox>
            </div>
          </el-checkbox-group>
        </div>

        <el-empty
          v-if="filteredGroupedUsers.length === 0"
          description="没有找到匹配的用户"
          :image-size="100"
        />
      </div>

      <template #footer>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span style="color: #909399; font-size: 14px">
            已选择 {{ selectedMembersTemp.length }} 个成员
          </span>
          <div>
            <el-button @click="memberSelectorVisible = false">取消</el-button>
            <el-button type="primary" @click="confirmMemberSelection">确定</el-button>
          </div>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { projectApi, userApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  Download,
  Edit,
  Delete,
  View,
  Rank,
  InfoFilled
} from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import TableToolbar from '@/components/common/TableToolbar.vue'
import ProjectTree from '@/components/progress/ProjectTree.vue'

const authStore = useAuthStore()
const router = useRouter()

// 项目树相关
const treeDrawerVisible = ref(false)
const selectedProjectId = ref(null)
const parentProjectList = ref([])

// 树形数据
const treeData = ref([])
// 展开的行
const expandedKeys = ref([])

// 列表数据
const loading = ref(false)
const tableData = ref([])
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: '',
  parent_id: null
})

// 用户列表
const userList = ref([])

// 对话框
const dialogVisible = ref(false)
const isViewMode = ref(false)
const dialogTitle = computed(() => {
  if (isViewMode.value) return '查看项目'
  return formData.id ? '编辑项目' : '创建项目'
})
const dialogLoading = ref(false)
const formRef = ref(null)

// 表单数据
const formData = reactive({
  id: null,
  code: '',
  name: '',
  location: '',
  start_date: '',
  end_date: '',
  status: 'planning',
  progress: 0,
  progress_percentage: 0,
  description: '',
  member_ids: [],
  parent_id: null
})

// 成员选择器相关
const memberSelectorVisible = ref(false)
const memberSearchText = ref('')
const selectedMembersTemp = ref([])
const originalMemberIds = ref([])

// 打开成员选择器
const openMemberSelector = () => {
  // 保存原始选中的成员ID
  originalMemberIds.value = [...(formData.member_ids || [])]
  // 初始化临时选择
  selectedMembersTemp.value = [...(formData.member_ids || [])]
  // 清空搜索
  memberSearchText.value = ''
  // 打开对话框
  memberSelectorVisible.value = true
}

// 移除成员
const removeMember = (userId) => {
  formData.member_ids = formData.member_ids.filter(id => id !== userId)
}

// 根据ID获取用户名
const getUserNameById = (userId) => {
  const user = userList.value.find(u => u.id === userId)
  if (user) {
    return user.full_name || user.username
  }
  return `用户${userId}`
}

// 确认成员选择
const confirmMemberSelection = () => {
  formData.member_ids = [...selectedMembersTemp.value]
  memberSelectorVisible.value = false
}

// 获取分组中已选择的数量
const getGroupSelectedCount = (group) => {
  return group.users.filter(u => selectedMembersTemp.value.includes(u.id)).length
}

// 处理分组全选/取消
const handleGroupCheck = (group) => {
  const userIds = group.users.map(u => u.id)
  if (group.checked) {
    // 全选：添加该组所有用户
    userIds.forEach(id => {
      if (!selectedMembersTemp.value.includes(id)) {
        selectedMembersTemp.value.push(id)
      }
    })
  } else {
    // 取消全选：移除该组所有用户
    selectedMembersTemp.value = selectedMembersTemp.value.filter(id => !userIds.includes(id))
  }
}

// 过滤并分组的用户（支持搜索）
const filteredGroupedUsers = computed(() => {
  const searchText = memberSearchText.value.toLowerCase().trim()
  const groups = groupedUsers.value.map(group => {
    // 过滤用户
    const filteredUsers = group.users.filter(user => {
      if (!searchText) return true
      return (
        (user.username && user.username.toLowerCase().includes(searchText)) ||
        (user.full_name && user.full_name.toLowerCase().includes(searchText)) ||
        (user.email && user.email.toLowerCase().includes(searchText))
      )
    })

    // 计算该组是否全选或半选
    const checked = filteredUsers.length > 0 && filteredUsers.every(u => selectedMembersTemp.value.includes(u.id))
    const indeterminate = filteredUsers.some(u => selectedMembersTemp.value.includes(u.id)) && !checked

    return {
      ...group,
      users: filteredUsers,
      checked,
      indeterminate
    }
  })

  // 只返回有用户的分组
  return groups.filter(g => g.users.length > 0)
})

// 表单验证规则
const formRules = {
  code: [
    { required: true, message: '请输入项目编码', trigger: 'blur' },
    { pattern: /^[A-Z0-9-]+$/, message: '项目编码格式不正确（大写字母、数字、连字符）', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入项目名称', trigger: 'blur' }
  ],
  location: [
    { required: true, message: '请输入项目地址', trigger: 'blur' }
  ],
  start_date: [
    { required: true, message: '请选择开始日期', trigger: 'change' }
  ],
  end_date: [
    { required: true, message: '请选择结束日期', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择项目状态', trigger: 'change' }
  ]
}

// 获取列表数据
const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      search: searchForm.keyword || undefined, // 使用正确的参数名
      status: searchForm.status || undefined,
      parent_id: searchForm.parent_id || undefined
    }
    const { data, pagination: pag } = await projectApi.getList(params)
    tableData.value = data || []
    pagination.total = pag?.total || 0

    // 构建树形数据
    if (!searchForm.parent_id && !searchForm.keyword) {
      // 只在没有筛选条件时构建完整树
      treeData.value = buildProjectTree(data || [])
    } else {
      // 有筛选条件时，只显示筛选结果（平铺）
      treeData.value = (data || []).map(item => ({ ...item, children: [] }))
    }

    // 同时获取所有项目用于父项目显示
    await fetchParentProjects()
  } catch (error) {
    console.error('获取项目列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索项目名称（用于自动完成）
const searchProjectsByName = async (queryString, callback) => {
  if (!queryString || queryString.trim().length === 0) {
    callback([])
    return
  }

  try {
    const { data } = await projectApi.getList({
      search: queryString.trim(), // 使用正确的参数名
      page: 1,
      page_size: 20
    })

    // 过滤：只显示没有父项目的项目（顶级项目）
    // 子项目不应该再成为其他项目的子项目
    const suggestions = (data || [])
      .filter(project => !project.parent_id) // 只显示顶级项目
      .map(project => ({
        id: project.id,
        name: project.name,
        code: project.code,
        location: project.location,
        start_date: project.start_date,
        end_date: project.end_date,
        status: project.status,
        parent_id: project.parent_id
      }))

    callback(suggestions)
  } catch (error) {
    console.error('搜索项目失败:', error)
    callback([])
  }
}

// 选择项目时的处理
const handleSelectProject = (item) => {
  // 填充项目信息
  formData.name = item.name
  formData.code = item.code
  formData.location = item.location
  formData.start_date = item.start_date
  formData.end_date = item.end_date
  formData.status = item.status

  // 如果已有ID，说明是选择已有项目
  if (item.id) {
    formData.id = item.id

    // 如果已经选择了父项目，说明要将此项目添加为子项目
    if (formData.parent_id) {
      const parentName = parentProjectList.value.find(p => p.id === formData.parent_id)?.name || '未知项目'
      ElMessage.success(`已选择项目"${item.name}"，将添加为"${parentName}"的子项目`)
    } else {
      ElMessage.info(`已选择已有项目"${item.name}"，可以直接保存或选择父项目`)
    }
  }
}

// 构建项目树形数据
const buildProjectTree = (projects) => {
  if (!projects || projects.length === 0) return []

  // 创建项目映射
  const projectMap = new Map()
  projects.forEach(project => {
    projectMap.set(project.id, { ...project, children: [] })
  })

  const roots = []

  // 构建树形结构
  projects.forEach(project => {
    const node = projectMap.get(project.id)
    if (!project.parent_id) {
      // 没有父项目，作为根节点
      roots.push(node)
    } else {
      // 有父项目，添加到父项目的children中
      const parent = projectMap.get(project.parent_id)
      if (parent) {
        parent.children.push(node)
      } else {
        // 父项目不在列表中，作为根节点
        roots.push(node)
      }
    }
  })

  // 递归设置hasChildren属性
  const setHasChildren = (nodes) => {
    nodes.forEach(node => {
      node.hasChildren = node.children && node.children.length > 0
      if (node.hasChildren) {
        setHasChildren(node.children)
      }
    })
  }
  setHasChildren(roots)

  return roots
}

// 展开/折叠行变化
const handleExpandChange = (row, expandedRows) => {
  expandedKeys.value = expandedRows.map(r => r.id)
}

// 加载用户列表
// 适配统一响应格式
const fetchUsers = async () => {
  try {
    const { data } = await userApi.getList({ pageSize: 1000 })
    userList.value = data || []
  } catch (error) {
    console.error('获取用户列表失败:', error)
  }
}

// 按角色分组用户
const groupedUsers = computed(() => {
  const groups = {}
  const users = userList.value || []

  users.forEach(user => {
    // 获取用户角色（优先使用roles数组中的第一个角色名称，否则使用role字段）
    let roleName = '其他用户'
    if (user.roles && user.roles.length > 0 && user.roles[0].name) {
      roleName = user.roles[0].name
    } else if (user.role) {
      roleName = user.role
    } else if (user.full_name) {
      roleName = user.full_name
    }

    // 初始化分组
    if (!groups[roleName]) {
      groups[roleName] = []
    }

    // 添加用户到分组
    groups[roleName].push(user)
  })

  // 转换为数组格式
  return Object.keys(groups).map(roleName => ({
    label: roleName,
    users: groups[roleName]
  })).sort((a, b) => {
    // 将"系统管理员"和"项目经理"排在前面
    const priority = { '系统管理员': 1, '项目经理': 2, '其他用户': 999 }
    const aPriority = priority[a.label] ?? 999
    const bPriority = priority[b.label] ?? 999
    return aPriority - bPriority
  })
})

// 加载项目成员
// 适配统一响应格式
const fetchProjectMembers = async (projectId) => {
  try {
    const { data } = await projectApi.getMembers(projectId)
    const members = data || []
    // 兼容不同的ID字段名称（id, user_id, ID等）
    formData.member_ids = members.map(m => m.id || m.user_id || m.ID).filter(Boolean)
  } catch (error) {
    console.error('获取项目成员失败:', error)
    formData.member_ids = []
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

// 防抖搜索（输入时延迟执行，避免频繁请求）
let searchTimer = null
const debouncedSearch = () => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  searchTimer = setTimeout(() => {
    pagination.page = 1
    fetchData()
  }, 300) // 300毫秒延迟
}

// 监听搜索关键词变化，实现即时搜索
watch(() => searchForm.keyword, () => {
  debouncedSearch()
})

// 监听状态筛选变化，实现即时搜索
watch(() => searchForm.status, () => {
  pagination.page = 1
  fetchData()
})

// 监听临时选择变化，更新分组的checked/indeterminate状态
watch(() => selectedMembersTemp.value, () => {
  // filteredGroupedUsers computed会自动重新计算
}, { deep: true })

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  searchForm.parent_id = null
  pagination.page = 1
  fetchData()
}

// 查看
const handleView = async (row) => {
  Object.assign(formData, {
    id: row.id,
    code: row.code,
    name: row.name,
    location: row.location,
    start_date: row.start_date,
    end_date: row.end_date,
    status: row.status,
    progress: row.progress || 0,
    progress_percentage: row.progress_percentage || 0,
    description: row.description || '',
    member_ids: [],
    parent_id: row.parent_id
  })
  // 加载项目成员
  await fetchProjectMembers(row.id)
  isViewMode.value = true
  dialogVisible.value = true
}

// 新增
const handleAdd = async () => {
  resetForm()
  isViewMode.value = false
  // 加载父项目列表
  await fetchParentProjects()
  dialogVisible.value = true
}

// 编辑
const handleEdit = async (row) => {
  Object.assign(formData, {
    id: row.id,
    code: row.code,
    name: row.name,
    location: row.location,
    start_date: row.start_date,
    end_date: row.end_date,
    status: row.status,
    progress: row.progress || 0,
    progress_percentage: row.progress_percentage || 0,
    description: row.description || '',
    member_ids: [],
    parent_id: row.parent_id
  })
  // 加载项目成员
  await fetchProjectMembers(row.id)
  // 加载父项目列表
  await fetchParentProjects()
  isViewMode.value = false
  dialogVisible.value = true
}

// 查看项目树
const handleViewTree = (row) => {
  selectedProjectId.value = row.id
  treeDrawerVisible.value = true
}

// 直接添加子项目
const handleAddChildDirect = (row) => {
  // 检查层级限制
  const currentLevel = row.level || 0
  if (currentLevel >= 3) {
    ElMessage.warning(`"${row.name}" 已是第 ${getLevelText(currentLevel)}，无法再创建子项目（最多支持4级项目）`)
    return
  }

  handleAdd()
  formData.parent_id = row.id
  fetchParentProjects()

  // 显示提示：说明可以搜索已有项目
  ElMessage({
    message: `正在为"${row.name}"添加子项目。可以在项目名称中搜索已有项目，或直接输入新项目名称`,
    type: 'info',
    duration: 4000,
    showClose: true
  })
}

// 从树中添加子项目
const handleAddChild = (parentProject) => {
  // 检查层级限制
  const currentLevel = parentProject.level || 0
  if (currentLevel >= 3) {
    ElMessage.warning(`"${parentProject.name}" 已是第 ${getLevelText(currentLevel)}，无法再创建子项目（最多支持4级项目）`)
    return
  }

  handleAdd()
  formData.parent_id = parentProject.id
  fetchParentProjects()
  treeDrawerVisible.value = false

  // 显示提示：说明可以搜索已有项目
  ElMessage({
    message: `正在为"${parentProject.name}"添加子项目。可以在项目名称中搜索已有项目，或直接输入新项目名称`,
    type: 'info',
    duration: 4000,
    showClose: true
  })
}

// 查看进度计划
const handleViewSchedule = (project) => {
  router.push({
    name: 'Progress',
    query: { projectId: project.id }
  })
}

// 加载父项目列表
const fetchParentProjects = async () => {
  try {
    const { data } = await projectApi.getList({ pageSize: 1000 })
    parentProjectList.value = (data || []).filter(p => {
      // 只显示可以作为父项目的项目（level < 3）
      return (p.level || 0) < 3 && p.id !== formData.id
    })
  } catch (error) {
    console.error('获取父项目列表失败:', error)
    parentProjectList.value = []
  }
}

// 获取层级标签类型
const getLevelTagType = (level) => {
  const types = {
    0: '',
    1: 'success',
    2: 'warning',
    3: 'info'
  }
  return types[level] || 'info'
}

// 获取层级文本
const getLevelText = (level) => {
  const texts = {
    0: '主项目',
    1: '分部分项',
    2: '工作包',
    3: '活动'
  }
  return texts[level] || '未知'
}

// 获取父项目名称
const getParentProjectName = (parentId) => {
  const parent = parentProjectList.value.find(p => p.id === parentId)
  return parent ? parent.name : '未知'
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除项目"${row.name}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await projectApi.delete(row.id)
      ElMessage.success('删除成功')
      fetchData()
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

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    // 日期验证
    if (formData.start_date && formData.end_date) {
      if (new Date(formData.start_date) > new Date(formData.end_date)) {
        ElMessage.error('开始日期不能晚于结束日期')
        return
      }
    }

    dialogLoading.value = true

    const data = {
      code: formData.code,
      name: formData.name,
      location: formData.location,
      start_date: formData.start_date,
      end_date: formData.end_date,
      status: formData.status,
      progress: formData.progress,
      progress_percentage: formData.progress_percentage,
      description: formData.description,
      parent_id: formData.parent_id,
      member_ids: formData.member_ids || [] // 明确传递成员列表
    }

    let projectId = formData.id

    if (formData.id) {
      await projectApi.update(formData.id, data)
      ElMessage.success('更新成功')
    } else {
      const result = await projectApi.create(data)
      projectId = result.id || result.data?.id
      ElMessage.success('创建成功')
    }

    // 保存成员关联（只有在明确指定了成员时才调用此API）
    // 如果没有指定成员（空数组或未提供），则由后端的成员继承逻辑处理：
    // - 创建时：有父项目则继承父项目成员，否则添加创建者
    // - 编辑时：保持现有成员不变
    if (projectId && formData.member_ids && formData.member_ids.length > 0) {
      try {
        // 后端API使用批量替换模式：发送 user_ids 数组替换所有成员
        await projectApi.addMember(projectId, { user_ids: formData.member_ids })
      } catch (error) {
        console.error('保存成员关联失败:', error)
        ElMessage.warning('项目已保存，但成员关联保存失败')
      }
    }

    dialogVisible.value = false

    // 刷新数据（包括项目树、表格和父项目列表）
    await fetchData()
    await fetchParentProjects()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    dialogLoading.value = false
  }
}

// 导出
const handleExport = async () => {
  try {
    const response = await projectApi.export(searchForm)
    const blob = new Blob([response], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `项目列表_${new Date().getTime()}.xlsx`
    a.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    id: null,
    code: '',
    name: '',
    location: '',
    start_date: '',
    end_date: '',
    status: 'planning',
    progress: 0,
    progress_percentage: 0,
    description: '',
    member_ids: [],
    parent_id: null
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const types = {
    planning: 'info',
    active: 'primary',
    closed: 'success',
    on_hold: 'warning'
  }
  return types[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    planning: '筹备中',
    active: '进行中',
    closed: '已完成',
    on_hold: '已暂停'
  }
  return texts[status] || status
}

// 获取进度条状态
const getProgressStatus = (progress, status) => {
  if (status === 'closed') return 'success'
  if (status === 'on_hold') return 'exception'
  if (progress >= 100) return 'success'
  if (progress >= 80) return 'warning'
  return null
}

// 判断是否可编辑
const canEdit = (row) => {
  if (!authStore.hasPermission('project_edit')) return false
  return row.status !== 'completed'
}

// 判断是否可删除
const canDelete = (row) => {
  if (!authStore.hasPermission('project_delete')) return false
  return row.status === 'planning'
}

// 分页处理
const handlePageChange = (page) => {
  pagination.page = page
  fetchData()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchData()
}

onMounted(() => {
  fetchData()
  fetchUsers()
})
</script>

<style scoped>
.projects-container {
  padding: 0;
}

/* 项目名称自动完成样式 */
.project-suggestion {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 4px 0;
}

.suggestion-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.suggestion-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.suggestion-hint {
  font-size: 12px;
  color: #909399;
}

.view-switcher {
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mt-20 {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.text-gray {
  color: #909399;
}

.text-sm {
  font-size: 12px;
  margin-top: 4px;
}
</style>
