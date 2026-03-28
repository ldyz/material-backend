<template>
  <div class="material-categories">
    <el-card>
      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchForm.keyword"
          placeholder="搜索分类名称、编码"
          clearable
          style="width: 300px"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" :icon="Plus" @click="handleAdd">
          添加顶级分类
        </el-button>
      </div>

      <!-- 分类表格 -->
      <el-table
        v-loading="loading"
        :data="categoryList"
        row-key="id"
        border
        stripe
        default-expand-all
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        style="margin-top: 20px"
      >
        <el-table-column prop="name" label="分类名称" min-width="200" />
        <el-table-column prop="code" label="分类编码" width="120" />
        <el-table-column prop="level" label="层级" width="80" align="center">
          <template #default="scope">
            <el-tag size="small">{{ scope.row.level }}级</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="material_count" label="物资数量" width="100" align="center" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="scope">
            <el-button
              v-if="scope.row.level < 4"
              type="primary"
              size="small"
              :icon="Plus"
              link
              @click="handleAddChild(scope.row)"
            >
              添加子分类
            </el-button>
            <el-button
              type="primary"
              size="small"
              :icon="Edit"
              link
              @click="handleEdit(scope.row)"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              link
              @click="handleDelete(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 提示信息 -->
      <el-alert
        title="使用说明"
        type="info"
        :closable="false"
        style="margin-top: 20px"
      >
        <template #default>
          <div>• 系统支持最多4级分类</div>
          <div>• 点击"添加子分类"可以为当前分类添加下级分类</div>
          <div>• 删除分类前请确保没有子分类和物资使用该分类</div>
        </template>
      </el-alert>
    </el-card>

    <!-- 分类对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      :close-on-click-modal="false"
      :loading="dialogLoading"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item label="分类名称" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入分类名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="分类编码" prop="code">
          <el-input
            v-model="form.code"
            placeholder="请输入分类编码（大写字母）"
            maxlength="20"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="父分类" prop="parent_id" v-if="!isAddingChild">
          <el-cascader
            v-model="form.parent_id"
            :options="treeOptions"
            :props="{
              value: 'id',
              label: 'name',
              emitPath: false,
              checkStrictly: true,
              expandTrigger: 'hover'
            }"
            placeholder="请选择父分类（不选则为顶级分类）"
            clearable
            style="width: 100%"
          />
          <div class="text-gray">选择父分类后，将作为该分类的子分类</div>
        </el-form-item>
        <el-form-item label="所属分类" v-if="isAddingChild">
          <el-input :value="parentCategory?.name" disabled />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="dialogLoading" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Edit, Delete } from '@element-plus/icons-vue'
import { materialApi } from '@/api'

// 搜索表单
const searchForm = reactive({
  keyword: ''
})

// 加载状态
const loading = ref(false)

// 分类列表
const categoryList = ref([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = computed(() => {
  if (isAddingChild.value) {
    return `添加子分类 - ${parentCategory.value?.name}`
  }
  return isEdit.value ? '编辑物资分类' : '添加物资分类'
})
const dialogLoading = ref(false)

// 表单
const formRef = ref(null)
const isEdit = ref(false)
const isAddingChild = ref(false)
const parentCategory = ref(null)

const form = reactive({
  id: null,
  name: '',
  code: '',
  parent_id: null,
  sort: 0,
  remark: ''
})

const formRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入分类编码', trigger: 'blur' }
  ]
}

// 分类树选项（用于级联选择器）
const treeOptions = computed(() => {
  return buildCategoryTreeOptions(categoryList.value)
})

// 获取物资分类
const fetchCategories = async () => {
  loading.value = true
  try {
    const { data } = await materialApi.getCategories()
    // 后端已经返回树形结构，直接使用
    categoryList.value = data || []

    // 如果有搜索关键词，进行过滤
    if (searchForm.keyword) {
      const keyword = searchForm.keyword.toLowerCase()
      categoryList.value = filterCategoriesByKeyword(categoryList.value, keyword)
    }
  } catch (error) {
    console.error('获取物资分类失败:', error)
    ElMessage.error(error?.response?.data?.message || error?.message || '获取物资分类失败')
    categoryList.value = []
  } finally {
    loading.value = false
  }
}

// 根据关键词过滤分类
const filterCategoriesByKeyword = (categories, keyword) => {
  return categories.reduce((acc, cat) => {
    const matchName = cat.name?.toLowerCase().includes(keyword)
    const matchCode = cat.code?.toLowerCase().includes(keyword)

    if (matchName || matchCode) {
      acc.push(cat)
    } else if (cat.children && cat.children.length > 0) {
      const filteredChildren = filterCategoriesByKeyword(cat.children, keyword)
      if (filteredChildren.length > 0) {
        acc.push({ ...cat, children: filteredChildren })
      }
    }

    return acc
  }, [])
}

// 构建分类树选项（用于级联选择器）
const buildCategoryTreeOptions = (categories) => {
  return categories.map(cat => {
    const option = {
      id: cat.id,
      name: cat.name,
      value: cat.id,
      label: cat.name
    }

    if (cat.children && cat.children.length > 0) {
      option.children = buildCategoryTreeOptions(cat.children)
    }

    return option
  })
}

// 搜索
const handleSearch = () => {
  fetchCategories()
}

// 添加顶级分类
const handleAdd = () => {
  isEdit.value = false
  isAddingChild.value = false
  parentCategory.value = null
  resetForm()
  dialogVisible.value = true
}

// 添加子分类
const handleAddChild = (row) => {
  isEdit.value = false
  isAddingChild.value = true
  parentCategory.value = row
  resetForm()
  form.parent_id = row.id
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  isEdit.value = false
  isAddingChild.value = false
  parentCategory.value = null
  resetForm()
  Object.assign(form, {
    id: row.id,
    name: row.name,
    code: row.code,
    parent_id: row.parent_id || null,
    sort: row.sort,
    remark: row.remark || ''
  })
  isEdit.value = true
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除分类"${row.name}"吗？删除前请确保没有物资使用该分类。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await materialApi.deleteCategory(row.id)
    ElMessage.success('删除成功')
    fetchCategories()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除分类失败:', error)
      ElMessage.error(error?.response?.data?.message || error?.message || '删除失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    dialogLoading.value = true

    const data = {
      name: form.name,
      code: form.code.toUpperCase(),
      parent_id: isAddingChild.value ? form.parent_id : (form.parent_id || 0),
      sort: form.sort,
      remark: form.remark
    }

    if (form.id) {
      await materialApi.updateCategory(form.id, data)
      ElMessage.success('更新成功')
    } else {
      await materialApi.createCategory(data)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchCategories()
  } catch (error) {
    console.error('保存分类失败:', error)
    ElMessage.error(error?.response?.data?.message || error?.message || '保存失败')
  } finally {
    dialogLoading.value = false
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    id: null,
    name: '',
    code: '',
    parent_id: null,
    sort: 0,
    remark: ''
  })

  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.material-categories {
  padding: 20px;
}

.search-bar {
  display: flex;
  gap: 10px;
}

.text-gray {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}
</style>
