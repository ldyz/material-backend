<template>
  <div class="materials-container">
    <el-card shadow="never">
      <!-- 工具栏 -->
      <TableToolbar>
        <template #left>
          <ProjectSelector
            v-model="searchForm.project_id"
            :projects="projectList"
            placeholder="选择项目"
            width="200px"
            @change="handleProjectChange"
          />
          <el-select
            v-model="searchForm.plan_id"
            placeholder="选择计划"
            clearable
            style="width: 200px"
            @change="handleSearch"
          >
            <el-option
              v-for="plan in planList"
              :key="plan.id"
              :label="plan.plan_name || plan.plan_no"
              :value="plan.id"
            />
          </el-select>
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索物资名称、编码、规格"
            clearable
            style="width: 250px"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
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
        <el-table-column prop="material_code" label="物资编码" width="120" />
        <el-table-column prop="material_name" label="物资名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="specification" label="规格型号" width="120" show-overflow-tooltip />
        <el-table-column prop="category" label="分类" width="100">
          <template #default="scope">
            <el-tag size="small">{{ scope.row.category || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="unit" label="单位" width="70" />
        <el-table-column prop="planned_quantity" label="计划数量" width="90" align="right">
          <template #default="scope">
            {{ scope.row.planned_quantity || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="arrived_quantity" label="已到货" width="90" align="right">
          <template #default="scope">
            {{ scope.row.arrived_quantity || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="remaining_quantity" label="未到" width="90" align="right">
          <template #default="scope">
            {{ scope.row.remaining_quantity || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="到货进度" width="120" align="center">
          <template #default="scope">
            <el-progress
              :percentage="Math.round(scope.row.arrival_percent || 0)"
              :stroke-width="8"
            />
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="所属项目" width="180" show-overflow-tooltip />
        <el-table-column prop="plan_name" label="所属计划" width="150" show-overflow-tooltip />
        <el-table-column prop="plan_no" label="计划编号" width="130">
          <template #default="scope">
            <el-link type="primary" @click="handleViewPlan(scope.row.plan_id)">
              {{ scope.row.plan_no }}
            </el-link>
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

      <!-- 计划详情对话框 -->
      <PlanDetailDialog
        v-model="detailDialogVisible"
        :plan-id="viewingPlanId"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { materialApi, projectApi, materialPlanApi } from '@/api'
import { Search, Refresh } from '@element-plus/icons-vue'
import TableToolbar from '@/components/common/TableToolbar.vue'
import ProjectSelector from '@/components/common/ProjectSelector.vue'
import PlanDetailDialog from './PlanDetailDialog.vue'

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
  project_id: '',
  plan_id: ''
})

// 项目列表
const projectList = ref([])

// 计划列表
const planList = ref([])

// 计划详情对话框
const detailDialogVisible = ref(false)
const viewingPlanId = ref(0)

// 查看计划详情
const handleViewPlan = (planId) => {
  viewingPlanId.value = planId
  detailDialogVisible.value = true
}

// 获取列表数据
const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      search: searchForm.keyword || undefined,
      project_id: searchForm.project_id || undefined,
      plan_id: searchForm.plan_id || undefined
    }
    const { data, pagination: pag } = await materialApi.getList(params)
    tableData.value = data || []
    pagination.total = pag?.total || 0
  } catch (error) {
    console.error('获取物资列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载项目列表
const fetchProjects = async () => {
  try {
    const { data } = await projectApi.getList({ pageSize: 1000 })
    // 构建树形结构
    projectList.value = buildProjectTree(data || [])
  } catch (error) {
    console.error('获取项目列表失败:', error)
  }
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

// 加载计划列表
const fetchPlans = async () => {
  try {
    const params = { page_size: 1000 }
    if (searchForm.project_id) {
      params.project_id = searchForm.project_id
    }
    const { data } = await materialPlanApi.getPlans(params)
    planList.value = data || []
  } catch (error) {
    console.error('获取计划列表失败:', error)
  }
}

// 项目变化时加载计划列表
const handleProjectChange = () => {
  searchForm.plan_id = ''
  fetchPlans()
  handleSearch()
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.project_id = ''
  searchForm.plan_id = ''
  planList.value = []
  pagination.page = 1
  fetchData()
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

// 防抖定时器
let searchTimer = null

// 即时搜索函数（带防抖）
const debouncedSearch = () => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  searchTimer = setTimeout(() => {
    pagination.page = 1
    fetchData()
  }, 500)
}

// 监听搜索字段变化，实现即时搜索
watch(() => searchForm.keyword, debouncedSearch)

onMounted(() => {
  fetchProjects()
  fetchPlans()
  fetchData()
})
</script>

<style scoped>
.materials-container {
  padding: 0;
}

.mt-20 {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
