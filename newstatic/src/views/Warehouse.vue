<template>
  <div class="warehouse-container">
    <!-- TAB切换 -->
    <el-card shadow="never" class="tab-card">
      <el-radio-group v-model="activeTab" @change="handleTabChange" size="large">
        <el-radio-button label="inbound">
          <el-icon style="margin-right: 4px;"><Download /></el-icon>
          入库单
        </el-radio-button>
        <el-radio-button label="requisition">
          <el-icon style="margin-right: 4px;"><Upload /></el-icon>
          出库单
        </el-radio-button>
      </el-radio-group>
    </el-card>

    <!-- 内容区域 -->
    <div class="content-area">
      <!-- 入库单组件 -->
      <Inbound v-if="activeTab === 'inbound'" ref="inboundRef" />

      <!-- 出库单组件 -->
      <Requisitions v-if="activeTab === 'requisition'" ref="requisitionRef" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Download, Upload } from '@element-plus/icons-vue'
import Inbound from './Inbound.vue'
import Requisitions from './Requisitions.vue'

const activeTab = ref('inbound')
const inboundRef = ref(null)
const requisitionRef = ref(null)

// TAB切换
const handleTabChange = (tab) => {
  console.log('切换到:', tab)
}

onMounted(() => {
  // 初始化时加载入库单数据
  if (inboundRef.value) {
    inboundRef.value.fetchData && inboundRef.value.fetchData()
  }
})
</script>

<style scoped>
.warehouse-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.tab-card {
  margin-bottom: 16px;
}

.tab-card :deep(.el-card__body) {
  padding: 16px 20px;
}

.content-area {
  flex: 1;
  overflow: hidden;
}

:deep(.el-radio-group) {
  display: flex;
  gap: 16px;
}

:deep(.el-radio-button__inner) {
  padding: 12px 24px;
  font-size: 14px;
  display: flex;
  align-items: center;
}
</style>
