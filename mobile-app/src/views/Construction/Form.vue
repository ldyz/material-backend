<template>
  <div class="construction-form-page">
    <van-nav-bar
      :title="isEdit ? '编辑日志' : '新建日志'"
      left-arrow
      @click-left="onClickLeft"
    />

    <van-form @submit="onSubmit" class="form-content">
      <!-- 项目选择 -->
      <van-cell-group inset title="基本信息">
        <van-field
          v-model="formData.project_id"
          name="project_id"
          label="项目"
          placeholder="请选择项目"
          readonly
          is-link
          @click="showProjectPicker = true"
          :rules="[{ required: true, message: '请选择项目' }]"
        />
        <van-field
          v-model="formData.log_date"
          name="log_date"
          label="日期"
          placeholder="请选择日期"
          readonly
          is-link
          @click="showDatePicker = true"
          :rules="[{ required: true, message: '请选择日期' }]"
        />
      </van-cell-group>

      <!-- 天气信息 -->
      <van-cell-group inset title="天气信息">
        <van-field name="weather" label="天气">
          <template #input>
            <van-radio-group v-model="formData.weather" direction="horizontal">
              <van-radio name="sunny">晴</van-radio>
              <van-radio name="cloudy">多云</van-radio>
              <van-radio name="rainy">雨</van-radio>
              <van-radio name="snowy">雪</van-radio>
            </van-radio-group>
          </template>
        </van-field>
        <van-field
          v-model.number="formData.temperature"
          name="temperature"
          type="number"
          label="温度"
          placeholder="请输入温度（°C）"
        />
      </van-cell-group>

      <!-- 日志内容（富文本） -->
      <van-cell-group inset title="日志内容">
        <van-field name="content" label="内容" :rules="[{ required: true, message: '请输入内容' }]">
          <template #input>
            <div class="rich-editor">
              <!-- 工具栏 -->
              <div class="editor-toolbar">
                <van-button
                  size="small"
                  icon="bold"
                  @click="insertFormat('**', '**')"
                />
                <van-button
                  size="small"
                  icon="italic"
                  @click="insertFormat('*', '*')"
                />
                <van-button
                  size="small"
                  icon="heading"
                  @click="insertLine('## ')"
                />
                <van-button
                  size="small"
                  icon="list-switch"
                  @click="insertLine('- ')"
                />
              </div>

              <!-- 文本域 -->
              <van-field
                v-model="formData.content"
                type="textarea"
                placeholder="请输入日志内容（支持Markdown格式）"
                rows="10"
                :border="false"
              />

              <!-- 预览 -->
              <div v-if="formData.content" class="content-preview">
                <div class="preview-label">预览：</div>
                <div v-html="renderedContent" class="preview-content"></div>
              </div>
            </div>
          </template>
        </van-field>
      </van-cell-group>

      <!-- 图片上传 -->
      <van-cell-group inset title="现场图片">
        <van-field name="images" label="图片">
          <template #input>
            <div class="image-uploader">
              <van-uploader
                v-model="imageList"
                multiple
                :max-count="9"
                :after-read="afterRead"
                :before-delete="beforeDelete"
              />
            </div>
          </template>
        </van-field>
      </van-cell-group>

      <!-- 提交按钮 -->
      <div class="submit-button">
        <van-button
          round
          block
          type="primary"
          native-type="submit"
          :loading="submitting"
        >
          保存日志
        </van-button>
      </div>
    </van-form>

    <!-- 项目选择弹窗 -->
    <van-popup v-model:show="showProjectPicker" position="bottom">
      <van-picker
        :columns="projectColumns"
        @confirm="onProjectConfirm"
        @cancel="showProjectPicker = false"
      />
    </van-popup>

    <!-- 日期选择弹窗 -->
    <van-popup v-model:show="showDatePicker" position="bottom">
      <van-date-picker
        v-model="currentDate"
        @confirm="onDateConfirm"
        @cancel="showDatePicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { createConstructionLog, updateConstructionLog, getConstructionLogDetail, uploadImage } from '@/api/construction'
import { formatDate } from '@/utils/date'

const router = useRouter()
const route = useRoute()

const isEdit = computed(() => !!route.params.id && route.params.id !== 'create')
const submitting = ref(false)

const formData = ref({
  project_id: '',
  log_date: formatDate(new Date(), 'YYYY-MM-DD'),
  weather: 'sunny',
  temperature: null,
  content: '',
  images: [],
})

const imageList = ref([])
const showProjectPicker = ref(false)
const showDatePicker = ref(false)
const currentDate = ref([2025, 1, 24]) // [year, month, day]

// 项目选项（模拟数据）
const projectColumns = [
  { text: '项目A', value: '1' },
  { text: '项目B', value: '2' },
]

// Markdown渲染（简单实现）
const renderedContent = computed(() => {
  if (!formData.value.content) return ''

  let html = formData.value.content

  // 粗体
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
  // 斜体
  html = html.replace(/\*(.*?)\*/g, '<em>$1</em>')
  // 标题
  html = html.replace(/^## (.*$)/gm, '<h3>$1</h3>')
  html = html.replace(/^# (.*$)/gm, '<h2>$1</h2>')
  // 列表
  html = html.replace(/^- (.*$)/gm, '<li>$1</li>')
  // 换行
  html = html.replace(/\n/g, '<br>')

  return html
})

// 项目选择确认
function onProjectConfirm({ selectedOptions }) {
  formData.value.project_id = selectedOptions[0].value
  showProjectPicker.value = false
}

// 日期选择确认
function onDateConfirm(value) {
  // van-date-picker 返回数组 [year, month, day]
  const [year, month, day] = value
  formData.value.log_date = `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`
  showDatePicker.value = false
}

// 插入格式
function insertFormat(before, after) {
  const textarea = document.querySelector('.rich-editor textarea')
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = formData.value.content

  formData.value.content = text.substring(0, start) + before + text.substring(start, end) + after + text.substring(end)

  // 重新聚焦并设置光标位置
  setTimeout(() => {
    textarea.focus()
    textarea.setSelectionRange(start + before.length, end + before.length)
  }, 0)
}

// 插入行
function insertLine(prefix) {
  const textarea = document.querySelector('.rich-editor textarea')
  const start = textarea.selectionStart
  const text = formData.value.content

  // 在当前行前插入
  const lines = text.split('\n')
  let currentLine = 0
  let charCount = 0

  for (let i = 0; i < lines.length; i++) {
    charCount += lines[i].length + 1
    if (charCount > start) {
      currentLine = i
      break
    }
  }

  lines[currentLine] = prefix + lines[currentLine]
  formData.value.content = lines.join('\n')

  setTimeout(() => {
    textarea.focus()
  }, 0)
}

// 图片上传后处理
async function afterRead(file, detail) {
  if (Array.isArray(file)) {
    file.forEach(f => uploadSingleImage(f))
  } else {
    await uploadSingleImage(file)
  }
}

// 上传单张图片
async function uploadSingleImage(file) {
  if (file.status === 'uploading') {
    const formData = new FormData()
    formData.append('file', file.file)

    try {
      const response = await uploadImage(formData)
      // 上传接口返回：{ success: true, url: '...' }
      file.url = response.url
      file.status = 'done'

      // 添加到表单数据
      if (!formData.value.images.includes(response.url)) {
        formData.value.images.push(response.url)
      }
    } catch (error) {
      file.status = 'failed'
      showToast({
        type: 'fail',
        message: '上传失败',
      })
    }
  }
}

// 删除图片前确认
function beforeDelete(file) {
  return new Promise((resolve) => {
    showConfirmDialog({
      title: '删除图片',
      message: '确定要删除这张图片吗？',
    })
      .then(() => {
        const index = formData.value.images.indexOf(file.url)
        if (index > -1) {
          formData.value.images.splice(index, 1)
        }
        resolve(true)
      })
      .catch(() => {
        resolve(false)
      })
  })
}

// 提交表单
async function onSubmit() {
  try {
    submitting.value = true

    const data = {
      ...formData.value,
      images: formData.value.images.join(','),
    }

    if (isEdit.value) {
      await updateConstructionLog(route.params.id, data)
      showToast({
        type: 'success',
        message: '保存成功',
      })
    } else {
      await createConstructionLog(data)
      showToast({
        type: 'success',
        message: '创建成功',
      })
    }

    router.back()
  } catch (error) {
    showToast({
      type: 'fail',
      message: error.message || '保存失败',
    })
  } finally {
    submitting.value = false
  }
}

// 返回
function onClickLeft() {
  router.back()
}

// 加载编辑数据
async function loadEditData() {
  if (isEdit.value) {
    try {
      const response = await getConstructionLogDetail(route.params.id)
      // 注意：construction_log接口直接返回对象，不是标准格式
      formData.value = {
        project_id: response.project_id.toString(),
        log_date: response.log_date,
        weather: response.weather,
        temperature: response.temperature,
        content: response.content || '',
        images: response.images || [],
      }

      // 设置图片列表
      imageList.value = response.images.map(url => ({
        url,
        status: 'done',
      }))
    } catch (error) {
      showToast({
        type: 'fail',
        message: '加载失败',
      })
      router.back()
    }
  }
}

onMounted(() => {
  loadEditData()
})
</script>

<style scoped>
.construction-form-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.form-content {
  padding: 16px 0 80px 0;
}

.rich-editor {
  width: 100%;
}

.editor-toolbar {
  display: flex;
  gap: 8px;
  padding: 8px 0;
  border-bottom: 1px solid #ebedf0;
  margin-bottom: 8px;
}

.content-preview {
  margin-top: 12px;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
}

.preview-label {
  font-size: 12px;
  color: #969799;
  margin-bottom: 8px;
}

.preview-content {
  font-size: 14px;
  line-height: 1.6;
  color: #323233;
}

.preview-content :deep(h2),
.preview-content :deep(h3) {
  margin: 8px 0;
  font-weight: bold;
}

.preview-content :deep(strong) {
  font-weight: bold;
}

.preview-content :deep(em) {
  font-style: italic;
}

.preview-content :deep(li) {
  margin-left: 20px;
}

.image-uploader {
  width: 100%;
}

.submit-button {
  padding: 16px;
  margin-top: 16px;
}
</style>
