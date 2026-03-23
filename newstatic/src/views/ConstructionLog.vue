<template>
  <div class="construction-log-container">
    <el-card shadow="never">
      <!-- 工具栏 -->
      <TableToolbar>
        <template #left>
          <ProjectSelector
            v-model="searchForm.project_id"
            :projects="projectList"
            placeholder="选择项目（支持层级显示）"
            width="300px"
          />
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索日志标题、内容"
            clearable
            style="width: 220px"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 240px"
          />
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </template>
        <template #right>
          <el-button
            type="primary"
            :icon="Plus"
            @click="handleAdd"
            v-if="authStore.hasPermission('constructionlog_create')"
          >
            创建日志
          </el-button>
          <el-button
            type="warning"
            :icon="Download"
            @click="handleExport"
            v-if="authStore.hasPermission('constructionlog_export')"
          >
            导出
          </el-button>
        </template>
      </TableToolbar>

      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
      >
        <!-- 序号列已移除 -->
        <el-table-column prop="title" label="日志标题" min-width="200" show-overflow-tooltip fixed="left">
          <template #default="scope">
            <el-link type="primary" @click="handleView(scope.row)">
              {{ scope.row.title }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="项目名称" width="150" show-overflow-tooltip />
        <el-table-column prop="log_date" label="日志日期" width="110" align="center" sortable />
        <el-table-column prop="weather" label="天气" width="80" align="center">
          <template #default="scope">
            <span v-if="scope.row.weather === 'sunny'">晴</span>
            <span v-else-if="scope.row.weather === 'cloudy'">多云</span>
            <span v-else-if="scope.row.weather === 'overcast'">阴</span>
            <span v-else-if="scope.row.weather === 'rainy'">雨</span>
            <span v-else-if="scope.row.weather === 'snowy'">雪</span>
            <span v-else-if="scope.row.weather === 'windy'">大风</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="temperature" label="温度(°C)" width="90" align="center">
          <template #default="scope">
            {{ scope.row.temperature ? `${scope.row.temperature}°C` : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right" align="center">
          <template #default="scope">
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
      width="900px"
      :loading="dialogLoading"
      :show-cancel="!isViewMode"
      :confirm-text="isViewMode ? '关闭' : '确认'"
      @confirm="isViewMode ? dialogVisible = false : handleSubmit()"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="项目" prop="project_id">
              <el-select
                v-model="formData.project_id"
                placeholder="请选择项目"
                filterable
                style="width: 100%"
                :disabled="isViewMode"
              >
                <el-option
                  v-for="item in projectOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="日志日期" prop="log_date">
              <el-date-picker
                v-model="formData.log_date"
                type="date"
                placeholder="选择日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                :disabled="isViewMode"
              />
              <div class="form-tip">系统自动获取当前日期，可修改</div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="日志标题" prop="title">
              <el-input
                v-model="formData.title"
                placeholder="请输入日志标题"
                maxlength="100"
                :disabled="isViewMode"
              />
              <div class="form-tip">系统自动根据日期生成，可修改</div>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="天气" prop="weather">
              <el-select
                v-model="formData.weather"
                placeholder="自动获取"
                style="width: 100%"
                :disabled="isViewMode"
              >
                <el-option label="晴" value="sunny" />
                <el-option label="多云" value="cloudy" />
                <el-option label="阴" value="overcast" />
                <el-option label="雨" value="rainy" />
                <el-option label="雪" value="snowy" />
                <el-option label="大风" value="windy" />
              </el-select>
              <div class="form-tip">系统根据位置自动获取，可修改</div>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="温度(°C)" prop="temperature">
              <el-input-number
                v-model="formData.temperature"
                :min="-50"
                :max="60"
                :disabled="isViewMode"
                style="width: 100%"
              />
              <div class="form-tip">系统自动获取，可修改</div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="施工内容" prop="content">
          <RichTextEditor
            :key="`content-${editorKey}`"
            v-model="formData.content"
            :min-height="300"
            :max-length="10000"
            :read-only="isViewMode"
          />
        </el-form-item>
      </el-form>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { constructionLogApi, projectApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  Download,
  Edit,
  Delete,
  View
} from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import TableToolbar from '@/components/common/TableToolbar.vue'
import RichTextEditor from '@/components/common/RichTextEditor.vue'
import ProjectSelector from '@/components/common/ProjectSelector.vue'

const authStore = useAuthStore()

// 列表数据
const loading = ref(false)
const tableData = ref([])
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 项目选项
const projectOptions = ref([])

// 项目列表（树形结构）
const projectList = ref([])

// 搜索表单
const searchForm = reactive({
  keyword: '',
  project_id: null,
  date_range: []
})

// 对话框
const dialogVisible = ref(false)
const isViewMode = ref(false)
const dialogTitle = computed(() => {
  if (isViewMode.value) return '查看施工日志'
  return formData.id ? '编辑施工日志' : '创建施工日志'
})
const dialogLoading = ref(false)
const formRef = ref(null)
const editorKey = ref(0)

// 表单数据
const formData = reactive({
  id: null,
  project_id: null,
  title: '',
  log_date: new Date().toISOString().split('T')[0],
  weather: 'sunny',
  temperature: 25,
  content: '',
  progress: '',
  issues: '',
  remark: ''
})

// 表单验证规则
const formRules = {
  project_id: [
    { required: true, message: '请选择项目', trigger: 'change' }
  ],
  title: [
    { required: true, message: '请输入日志标题', trigger: 'blur' }
  ],
  log_date: [
    { required: true, message: '请选择日志日期', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入施工内容', trigger: 'blur' }
  ]
}

// 获取列表数据
// 适配统一响应格式
const fetchData = async () => {
  loading.value = true
  try {
    // 收集项目ID（包含子项目）
    let projectIds = []
    if (searchForm.project_id) {
      projectIds = collectProjectIds(searchForm.project_id, projectList.value)
    }

    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      search: searchForm.keyword || undefined,
      project_ids: projectIds.length > 0 ? projectIds.join(',') : undefined,
      start_date: searchForm.date_range?.[0] || undefined,
      end_date: searchForm.date_range?.[1] || undefined
    }
    const { data, pagination: pag } = await constructionLogApi.getList(params)
    tableData.value = data || []
    pagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取施工日志列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 收集项目及其所有子项目的ID
const collectProjectIds = (projectId, projectTree) => {
  const ids = [projectId]

  // 在树形结构中查找项目并收集子项目ID
  const findAndCollectChildren = (nodes) => {
    for (const node of nodes) {
      if (node.id === projectId) {
        // 找到目标项目，递归收集所有子项目ID
        const collectAllChildren = (project) => {
          if (project.children && project.children.length > 0) {
            for (const child of project.children) {
              ids.push(child.id)
              collectAllChildren(child)
            }
          }
        }
        collectAllChildren(node)
        return true
      }
      if (node.children && node.children.length > 0) {
        if (findAndCollectChildren(node.children)) {
          return true
        }
      }
    }
    return false
  }

  findAndCollectChildren(projectTree)
  return ids
}

// 构建项目树形结构
const buildProjectTree = (projects) => {
  if (!projects || projects.length === 0) return []

  const projectMap = new Map()
  projects.forEach(project => {
    projectMap.set(project.id, { ...project, children: [] })
  })

  const roots = []
  projects.forEach(project => {
    const node = projectMap.get(project.id)
    if (!project.parent_id) {
      roots.push(node)
    } else {
      const parent = projectMap.get(project.parent_id)
      if (parent) {
        parent.children.push(node)
      } else {
        roots.push(node)
      }
    }
  })

  return roots
}

// 获取项目列表
// 适配统一响应格式
const fetchProjects = async () => {
  try {
    const { data } = await projectApi.getList({ pageSize: 1000 })
    projectOptions.value = data || []
    projectList.value = buildProjectTree(data || [])
  } catch (error) {
    console.error('获取项目列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

// 防抖搜索
let searchTimer = null
const debouncedSearch = () => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  searchTimer = setTimeout(() => {
    pagination.page = 1
    fetchData()
  }, 300)
}

// 监听搜索条件变化
watch(() => searchForm.keyword, () => {
  debouncedSearch()
})

watch(() => searchForm.project_id, () => {
  pagination.page = 1
  fetchData()
})

watch(() => searchForm.date_range, () => {
  pagination.page = 1
  fetchData()
})

// 监听项目列表加载完成
watch(() => projectList.value, (newList) => {
  // 如果当前有选中的项目，重新加载数据
  if (searchForm.project_id && newList.length > 0) {
    fetchData()
  }
})

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.project_id = null
  searchForm.date_range = []
  pagination.page = 1
  fetchData()
}

// 获取天气信息
const fetchWeather = async () => {
  try {
    // 获取地理位置
    const position = await new Promise((resolve, reject) => {
      if (!navigator.geolocation) {
        reject(new Error('浏览器不支持地理定位'))
        return
      }
      navigator.geolocation.getCurrentPosition(resolve, reject, {
        enableHighAccuracy: false,
        timeout: 10000,
        maximumAge: 300000 // 5分钟缓存
      })
    })

    const { latitude, longitude } = position.coords

    // 调用Open-Meteo免费天气API
    const response = await fetch(
      `https://api.open-meteo.com/v1/forecast?latitude=${latitude}&longitude=${longitude}&current=temperature_2m,weather_code&timezone=auto`
    )
    const data = await response.json()

    if (data.current) {
      // 更新温度
      formData.temperature = Math.round(data.current.temperature_2m)

      // 根据天气代码设置天气类型
      const weatherCode = data.current.weather_code
      formData.weather = mapWeatherCode(weatherCode)
    }
  } catch (error) {
    console.warn('获取天气信息失败:', error.message)
    // 保持默认值（晴天，25度）
    formData.weather = 'sunny'
    formData.temperature = 25
  }
}

// 生成标题
const generateTitle = () => {
  const date = new Date(formData.log_date)
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${year}年${month}月${day}日 施工日志`
}

// 映射天气代码到我们的天气类型
const mapWeatherCode = (code) => {
  // WMO天气代码映射
  // 0: 晴天
  // 1-3: 多云/阴
  // 45, 48: 雾
  // 51-67: 雨
  // 71-77: 雪
  // 80-82: 阵雨
  // 85-86: 阵雪
  // 95-99: 雷暴
  if (code === 0) return 'sunny'
  if (code >= 1 && code <= 3) return 'cloudy'
  if (code >= 45 && code <= 48) return 'overcast'
  if ((code >= 51 && code <= 67) || (code >= 80 && code <= 82)) return 'rainy'
  if ((code >= 71 && code <= 77) || (code >= 85 && code <= 86)) return 'snowy'
  if (code >= 95) return 'windy'
  return 'sunny' // 默认晴天
}

// 新增
const handleAdd = async () => {
  resetForm()
  isViewMode.value = false
  dialogVisible.value = true
  fetchProjects()

  // 自动生成标题
  formData.title = generateTitle()

  // 自动获取天气信息
  await fetchWeather()
}

// 查看
const handleView = (row) => {
  Object.assign(formData, {
    id: row.id,
    project_id: row.project_id,
    title: row.title,
    log_date: row.log_date,
    weather: row.weather,
    temperature: row.temperature,
    content: row.content,
    progress: row.progress || '',
    issues: row.issues || '',
    remark: row.remark || ''
  })
  isViewMode.value = true
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  Object.assign(formData, {
    id: row.id,
    project_id: row.project_id,
    title: row.title,
    log_date: row.log_date,
    weather: row.weather,
    temperature: row.temperature,
    content: row.content,
    progress: row.progress || '',
    issues: row.issues || '',
    remark: row.remark || ''
  })
  isViewMode.value = false
  dialogVisible.value = true
  fetchProjects()
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除施工日志"${row.title}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await constructionLogApi.delete(row.id)
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
    dialogLoading.value = true

    const data = {
      project_id: formData.project_id,
      title: formData.title,
      log_date: formData.log_date,
      weather: formData.weather,
      temperature: formData.temperature,
      content: formData.content,
      progress: formData.progress,
      issues: formData.issues,
      remark: formData.remark
    }

    if (formData.id) {
      await constructionLogApi.update(formData.id, data)
      ElMessage.success('更新成功')
    } else {
      await constructionLogApi.create(data)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error(error?.message || '操作失败，请重试')
  } finally {
    dialogLoading.value = false
  }
}

// 导出
const handleExport = async () => {
  try {
    const response = await constructionLogApi.export(searchForm)
    const blob = new Blob([response], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `施工日志_${new Date().getTime()}.xlsx`
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
    project_id: null,
    title: '',
    log_date: new Date().toISOString().split('T')[0],
    weather: 'sunny',
    temperature: 25,
    content: '',
    progress: '',
    issues: '',
    remark: ''
  })
  editorKey.value++ // 强制重新创建富文本编辑器
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 判断是否可编辑
const canEdit = (row) => {
  return authStore.hasPermission('constructionlog_edit')
}

// 判断是否可删除
const canDelete = (row) => {
  if (!authStore.hasPermission('constructionlog_delete')) return false
  return authStore.userId === row.author_id
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

// 监听日期变化，自动更新标题
watch(() => formData.log_date, (newDate) => {
  // 只有在创建新日志且标题为默认格式时才自动更新
  if (!formData.id && newDate) {
    const newTitle = generateTitle()
    // 检查当前标题是否是自动生成的格式，如果是则更新
    const currentTitle = formData.title
    if (currentTitle.match(/^\d{4}年\d{1,2}月\d{1,2}日 施工日志$/)) {
      formData.title = newTitle
    }
  }
})

onMounted(() => {
  fetchProjects()
  fetchData()
})
</script>

<style scoped>
.construction-log-container {
  padding: 0;
}

.mt-20 {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.form-tip {
  font-size: 12px;
  color: #409eff;
  margin-top: 4px;
  line-height: 1.4;
}
</style>
