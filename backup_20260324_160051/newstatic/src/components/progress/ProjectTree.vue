<template>
  <div class="project-tree">
    <el-empty
      v-if="loading"
      description="加载中..."
      :image-size="100"
    />
    <el-empty
      v-else-if="!treeData || treeData.length === 0"
      description="暂无项目层级数据"
      :image-size="100"
    />
    <el-tree
      v-else
      :data="treeData"
      :props="treeProps"
      node-key="id"
      :expand-on-click-node="false"
      :default-expand-all="false"
      @node-click="handleNodeClick"
    >
      <template #default="{ node, data }">
        <div class="tree-node">
          <span class="node-label">
            <el-icon
              v-if="data.level > 0"
              :size="16"
              style="margin-right: 4px; color: #909399"
            >
              <component :is="getLevelIcon(data.level)" />
            </el-icon>
            {{ node.label }}
          </span>
          <div class="node-info">
            <el-tag
              :type="getLevelTagType(data.level)"
              size="small"
              style="margin-right: 8px"
            >
              {{ getLevelLabel(data.level) }}
            </el-tag>
            <el-progress
              :percentage="data.progress_percentage || 0"
              :stroke-width="6"
              :show-text="false"
              style="width: 80px; margin-right: 8px"
            />
            <span class="progress-text">{{
              data.progress_percentage || 0
            }}%</span>
            <el-button
              link
              type="primary"
              size="small"
              @click.stop="handleAddChild(data)"
              v-if="data.level < 3"
            >
              <el-icon><Plus /></el-icon>
              添加子项目
            </el-button>
            <el-button
              link
              type="primary"
              size="small"
              @click.stop="handleViewSchedule(data)"
            >
              <el-icon><Calendar /></el-icon>
              进度计划
            </el-button>
          </div>
        </div>
      </template>
    </el-tree>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Calendar, Folder, FolderOpened, Document } from '@element-plus/icons-vue'
import { progressApi, projectApi } from '@/api'
import { useRouter } from 'vue-router'

const props = defineProps({
  projectId: {
    type: [Number, String],
    required: true
  }
})

const emit = defineEmits(['project-selected', 'add-child', 'view-schedule'])

const router = useRouter()
const treeData = ref([])
const loading = ref(false)
const treeProps = {
  children: 'children',
  label: 'name'
}

// 获取项目树数据
const fetchProjectTree = async () => {
  loading.value = true
  try {
    console.log('获取项目树，projectId:', props.projectId)

    // 尝试使用 projectApi
    const response = await projectApi.getProjectTree(props.projectId)

    console.log('项目树响应:', response)

    if (response && response.data) {
      // 确保数据是数组格式
      if (Array.isArray(response.data)) {
        treeData.value = response.data
      } else {
        treeData.value = [response.data]
      }
      console.log('树数据已设置:', treeData.value)
    } else {
      treeData.value = []
      console.warn('项目树响应数据为空')
    }
  } catch (error) {
    console.error('获取项目树失败:', error)
    ElMessage.error('获取项目树失败: ' + (error.message || '未知错误'))
    treeData.value = []
  } finally {
    loading.value = false
  }
}

// 节点点击事件
const handleNodeClick = (data) => {
  emit('project-selected', data)
}

// 添加子项目
const handleAddChild = (data) => {
  emit('add-child', data)
}

// 查看进度计划
const handleViewSchedule = (data) => {
  emit('view-schedule', data)
  // 跳转到进度页面
  router.push({
    name: 'Progress',
    query: { projectId: data.id }
  })
}

// 获取层级图标
const getLevelIcon = (level) => {
  switch (level) {
    case 0:
      return Folder
    case 1:
      return FolderOpened
    case 2:
      return Document
    default:
      return Document
  }
}

// 获取层级标签类型
const getLevelTagType = (level) => {
  switch (level) {
    case 0:
      return ''
    case 1:
      return 'success'
    case 2:
      return 'warning'
    case 3:
      return 'info'
    default:
      return 'info'
  }
}

// 获取层级标签文字
const getLevelLabel = (level) => {
  switch (level) {
    case 0:
      return '主项目'
    case 1:
      return '分部分项'
    case 2:
      return '工作包'
    case 3:
      return '活动'
    default:
      return '未知'
  }
}

onMounted(() => {
  fetchProjectTree()
})

// 暴露刷新方法
defineExpose({
  refresh: fetchProjectTree
})
</script>

<style scoped>
.project-tree {
  background: #fff;
  padding: 16px;
  border-radius: 4px;
}

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding-right: 16px;
}

.node-label {
  font-size: 14px;
  color: #303133;
  flex-shrink: 0;
}

.node-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.progress-text {
  font-size: 12px;
  color: #909399;
  min-width: 40px;
}

:deep(.el-tree-node__content) {
  height: 48px;
}

:deep(.el-tree-node__content:hover) {
  background-color: #f5f7fa;
}
</style>
