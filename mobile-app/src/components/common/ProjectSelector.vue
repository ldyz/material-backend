<template>
  <!-- 单元格模式（用于表单） -->
  <van-field
    v-if="mode === 'cell'"
    :name="name"
    :model-value="internalValue"
    :label="label"
    :placeholder="placeholder"
    :readonly="readonly"
    :required="required"
    :rules="computedRules"
    is-link
    clickable
    @click="handleClick"
  >
    <template #input>
      <div :class="['project-selector-input', { 'has-value': selectedProject }]">
        {{ displayText }}
      </div>
    </template>
  </van-field>

  <!-- 下拉菜单模式（用于筛选器） -->
  <van-dropdown-item
    v-else-if="mode === 'dropdown'"
    v-model="internalValue"
    :options="dropdownOptions"
    :title="label"
    @change="handleDropdownChange"
  />

  <!-- 弹窗选择器 -->
  <van-popup
    v-model:show="showPicker"
    position="bottom"
    :style="{ height: '50%' }"
  >
    <van-picker
      :model-value="pickerValue"
      :columns="pickerColumns"
      :loading="loading"
      @confirm="handleConfirm"
      @cancel="showPicker = false"
    >
      <template #title>
        <div class="picker-title">{{ label }}</div>
      </template>
    </van-picker>
  </van-popup>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { getProjects } from '@/api/project'
import { logger } from '@/utils/logger'

const props = defineProps({
  /**
   * 组件模式
   * - cell: 单元格模式（用于表单）
   * - dropdown: 下拉菜单模式（用于筛选器）
   */
  mode: {
    type: String,
    default: 'cell'
  },
  /**
   * 字段名称（用于表单验证）
   */
  name: {
    type: String,
    default: 'project_id'
  },
  /**
   * 标签文本
   */
  label: {
    type: String,
    default: '项目'
  },
  /**
   * 占位符文本
   */
  placeholder: {
    type: String,
    default: '请选择项目'
  },
  /**
   * 是否必填
   */
  required: {
    type: Boolean,
    default: false
  },
  /**
   * 是否只读
   */
  readonly: {
    type: Boolean,
    default: false
  },
  /**
   * 选中的项目ID（v-model）
   */
  modelValue: {
    type: [Number, String],
    default: null
  },
  /**
   * 验证规则（仅用于cell模式）
   */
  rules: {
    type: Array,
    default: () => []
  },
  /**
   * 项目状态筛选
   */
  statusFilter: {
    type: String,
    default: '' // 空=全部, active=进行中, closed=已关闭
  }
})

const emit = defineEmits(['update:modelValue', 'change', 'select'])

// 计算验证规则
const computedRules = computed(() => {
  const rules = [...props.rules]
  if (props.required) {
    rules.unshift({
      required: true,
      message: `请选择${props.label}`
    })
  }
  return rules
})

// 内部值
const internalValue = ref(props.modelValue)
const showPicker = ref(false)
const loading = ref(false)
const projects = ref([])

// 获取项目列表
async function fetchProjects() {
  loading.value = true
  try {
    const params = {
      page: 1,
      page_size: 1000, // 获取所有项目
      show_all: 'true' // 显示所有项目，不受用户关联限制
    }
    if (props.statusFilter) {
      params.status = props.statusFilter
    }
    const response = await getProjects(params)
    // 后端返回格式：{ success: true, data: [...], pagination: {...} }
    // axios 拦截器返回 response.data，所以这里的 response 就是上面的对象
    // response.data 直接就是项目数组
    projects.value = response.data || []
  } catch (error) {
    logger.error('获取项目列表失败:', error)
    projects.value = []
  } finally {
    loading.value = false
  }
}

// 选择器列配置
const pickerColumns = computed(() => {
  if (!projects.value || projects.value.length === 0) {
    return []
  }
  return projects.value.map(p => ({
    text: p.name,
    value: p.id
  }))
})

// Picker 选中的值（数组格式）
const pickerValue = computed(() => {
  if (internalValue.value === null || internalValue.value === undefined || internalValue.value === '') {
    return []
  }
  return [internalValue.value]
})

// 下拉菜单选项
const dropdownOptions = computed(() => {
  return [
    { text: '全部', value: '' },
    ...projects.value.map(p => ({
      text: p.name,
      value: p.id
    }))
  ]
})

// 当前选中的项目
const selectedProject = computed(() => {
  if (!internalValue.value) return null
  return projects.value.find(p => p.id === Number(internalValue.value))
})

// 显示文本
const displayText = computed(() => {
  if (selectedProject.value) {
    return selectedProject.value.name
  }
  return props.placeholder
})

// 处理点击（cell模式）
function handleClick() {
  if (!props.readonly) {
    showPicker.value = true
  }
}

// 处理确认
function handleConfirm({ selectedOptions, selectedValues }) {
  // 兼容两种参数格式
  const value = selectedValues?.[0] ?? selectedOptions?.[0]?.value
  if (value !== undefined) {
    internalValue.value = value
    emit('update:modelValue', value)
    emit('change', value)
    emit('select', selectedProject.value)
  }
  showPicker.value = false
}

// 处理下拉菜单变化
function handleDropdownChange(value) {
  emit('update:modelValue', value || null)
  emit('change', value || null)
  const project = value ? projects.value.find(p => p.id === Number(value)) : null
  emit('select', project)
}

// 监听外部值变化
watch(() => props.modelValue, (newVal) => {
  internalValue.value = newVal
})

// 初始化
onMounted(() => {
  fetchProjects()
})

// 暴露刷新方法
defineExpose({
  refresh: fetchProjects
})
</script>

<style scoped>
.project-selector-input {
  color: #969799;
  min-height: 24px;
  display: flex;
  align-items: center;
}

.project-selector-input.has-value {
  color: #323233;
}

.picker-title {
  font-size: 16px;
  font-weight: 500;
  color: #323233;
  padding: 12px 0;
}
</style>
