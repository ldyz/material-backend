<template>
  <div class="ai-page">
    <!-- 智能问答 -->
    <div class="chat-section">
      <div class="chat-header">
        <h2>AI 智能助手</h2>
        <p class="subtitle">用自然语言查询数据，获取智能分析</p>
      </div>

      <!-- 对话历史 -->
      <div class="chat-messages" ref="messagesContainer">
        <div
          v-for="(message, index) in chatHistory"
          :key="index"
          class="message"
          :class="message.role"
        >
          <div class="message-avatar">
            <van-icon
              :name="message.role === 'user' ? 'user' : 'robot-o'"
              size="20"
            />
          </div>
          <div class="message-content">
            <div class="message-text" v-html="formatMessage(message.content)"></div>
            <div v-if="message.data" class="message-data">
              <!-- 数据可视化 -->
              <div v-if="message.data.chart" class="data-chart">
                <div class="chart-title">{{ message.data.chart.title }}</div>
                <div class="chart-bars">
                  <div
                    v-for="(item, i) in message.data.chart.data"
                    :key="i"
                    class="chart-bar-wrapper"
                  >
                    <div
                      class="chart-bar"
                      :style="{ height: `${item.percentage}%` }"
                    ></div>
                    <span class="chart-label">{{ item.label }}</span>
                  </div>
                </div>
              </div>

              <!-- 表格数据 -->
              <div v-if="message.data.table" class="data-table">
                <div
                  v-for="(row, i) in message.data.table.data"
                  :key="i"
                  class="table-row"
                >
                  <div
                    v-for="col in message.data.table.columns"
                    :key="col.key"
                    class="table-cell"
                  >
                    <span class="cell-label">{{ col.label }}:</span>
                    <span class="cell-value">{{ row[col.key] }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 加载动画 -->
        <div v-if="analyzing" class="message assistant">
          <div class="message-avatar">
            <van-icon name="robot-o" size="20" />
          </div>
          <div class="message-content">
            <van-loading size="16" color="#1989fa" />
          </div>
        </div>
      </div>

      <!-- 输入框 -->
      <div class="chat-input">
        <van-field
          v-model="userQuestion"
          placeholder="问：本月入库金额是多少？"
          :border="false"
          @keyup.enter="sendQuestion"
        >
          <template #button>
            <van-button
              type="primary"
              size="small"
              :loading="analyzing"
              @click="sendQuestion"
            >
              发送
            </van-button>
          </template>
        </van-field>
      </div>
    </div>

    <!-- 快捷问题 -->
    <div class="quick-questions">
      <div class="section-title">快捷问题</div>
      <div class="questions-list">
        <van-tag
          v-for="(question, index) in quickQuestions"
          :key="index"
          plain
          type="primary"
          size="medium"
          @click="askQuestion(question)"
        >
          {{ question }}
        </van-tag>
      </div>
    </div>

    <!-- 数据洞察 -->
    <div class="insights-section">
      <div class="section-title">数据洞察</div>
      <div class="insights-cards">
        <div
          v-for="(insight, index) in insights"
          :key="index"
          class="insight-card"
        >
          <div class="insight-icon">
            <van-icon :name="insight.icon" />
          </div>
          <div class="insight-info">
            <p class="insight-title">{{ insight.title }}</p>
            <p class="insight-value">{{ insight.value }}</p>
            <p class="insight-trend" :class="insight.trendClass">
              {{ insight.trend }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- 推荐建议 -->
    <div class="recommendations-section">
      <div class="section-title">智能推荐</div>
      <van-cell-group inset>
        <van-cell
          v-for="(rec, index) in recommendations"
          :key="index"
          :title="rec.title"
          :label="rec.description"
          is-link
          @click="handleRecommendation(rec)"
        />
      </van-cell-group>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { analyzeQuestion, getInsights, getRecommendations } from '@/api/ai'

const router = useRouter()

const chatHistory = ref([])
const userQuestion = ref('')
const analyzing = ref(false)
const messagesContainer = ref(null)
const insights = ref([])
const recommendations = ref([])

// 快捷问题（从后端获取）
const quickQuestions = ref([])

// 发送问题
async function sendQuestion() {
  const question = userQuestion.value.trim()
  if (!question) {
    showToast('请输入问题')
    return
  }

  // 添加用户消息
  chatHistory.value.push({
    role: 'user',
    content: question,
  })

  userQuestion.value = ''
  analyzing.value = true

  try {
    const response = await analyzeQuestion(question)

    // 注意：AI分析接口直接返回结果对象，不是标准格式
    // 添加AI回复
    chatHistory.value.push({
      role: 'assistant',
      content: response.answer,
      data: response.visualization, // 图表或表格数据
    })

    // 滚动到底部
    scrollToBottom()
  } catch (error) {
    // 检测配额用完错误
    const errorMsg = error.message || error.toString() || ''
    if (errorMsg.includes('quota') || errorMsg.includes('exceeded')) {
      chatHistory.value.push({
        role: 'assistant',
        content: '抱歉，AI 服务配额已用完，请联系管理员充值或稍后再试。',
      })
    } else {
      chatHistory.value.push({
        role: 'assistant',
        content: `抱歉，处理请求时出错：${errorMsg}`,
      })
    }
  } finally {
    analyzing.value = false
  }
}

// 快捷提问
function askQuestion(question) {
  userQuestion.value = question
  sendQuestion()
}

// 格式化消息内容（支持Markdown）
function formatMessage(content) {
  if (!content) return ''

  // 简单的 Markdown 转换
  let html = content

  // 粗体
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
  // 代码
  html = html.replace(/`(.*?)`/g, '<code>$1</code>')
  // 换行
  html = html.replace(/\n/g, '<br>')

  return html
}

// 滚动到底部
function scrollToBottom() {
  setTimeout(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  }, 100)
}

// 加载数据洞察（真实数据）
async function loadInsights() {
  try {
    // 适配统一响应格式
    const { data } = await getInsights({ type: 'dashboard' })

    insights.value = [
      {
        icon: 'gold-coin-o',
        title: '本月入库金额',
        value: `¥${(data.total_inbound || 0).toFixed(2)}`,
        trend: `环比 ${data.inbound_growth >= 0 ? '+' : ''}${(data.inbound_growth || 0).toFixed(1)}%`,
        trendClass: data.inbound_growth >= 0 ? 'up' : 'down',
      },
      {
        icon: 'send-gift-o',
        title: '本月出库金额',
        value: `¥${(data.total_outbound || 0).toFixed(2)}`,
        trend: `环比 ${data.outbound_growth >= 0 ? '+' : ''}${(data.outbound_growth || 0).toFixed(1)}%`,
        trendClass: data.outbound_growth >= 0 ? 'up' : 'down',
      },
      {
        icon: 'warning-o',
        title: '库存预警',
        value: `${data.alert_count || 0} 种`,
        trend: data.alert_count > 0 ? '需要关注' : '库存正常',
        trendClass: data.alert_count > 0 ? 'warning' : 'up',
      },
      {
        icon: 'clock-o',
        title: '待审批单据',
        value: `${data.pending_count || 0} 条`,
        trend: data.pending_count > 0 ? '待处理' : '无待办',
        trendClass: data.pending_count > 0 ? 'info' : 'up',
      },
    ]
  } catch (error) {
    console.error('加载洞察失败:', error.message)

    // 显示错误提示给用户
    if (error.message && (error.message.includes('quota') || error.message.includes('exceeded'))) {
      showToast('AI 配额已用完，暂无法加载数据洞察')
    } else {
      showToast('加载数据洞察失败，请稍后重试')
    }

    // 保留空状态，不显示假数据
    insights.value = []
  }
}

// 加载推荐建议（真实数据）
async function loadRecommendations() {
  try {
    // 适配统一响应格式
    const { recommendations } = await getRecommendations('recommendations')
    recommendations.value = recommendations || []
  } catch (error) {
    console.error('加载推荐失败:', error.message)

    if (error.message && (error.message.includes('quota') || error.message.includes('exceeded'))) {
      showToast('AI 配额已用完，暂无法加载智能推荐')
    } else {
      showToast('加载智能推荐失败，请稍后重试')
    }

    // 空状态
    recommendations.value = []
  }
}

// 处理推荐建议点击
function handleRecommendation(rec) {
  if (rec.action) {
    router.push(rec.action)
  } else {
    showToast('功能开发中')
  }
}

// 加载快捷问题
async function loadQuickQuestions() {
  try {
    // 适配统一响应格式
    const { suggestions } = await getRecommendations('questions')
    quickQuestions.value = suggestions || []
  } catch (error) {
    console.error('加载快捷问题失败:', error.message)
    // 使用默认问题作为后备
    quickQuestions.value = [
      '本月入库金额是多少？',
      '库存预警有哪些材料？',
      '待审批的入库单有多少？',
      '本周出库最多的材料？',
      '哪个项目的采购金额最高？',
    ]
  }
}

onMounted(() => {
  // 添加欢迎消息
  chatHistory.value.push({
    role: 'assistant',
    content: '您好！我是您的AI智能助手，可以帮您查询数据、分析趋势。请问有什么可以帮助您的？',
  })

  // 并行加载数据
  loadInsights()
  loadRecommendations()
  loadQuickQuestions()
})
</script>

<style scoped>
.ai-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 50px;
}

.chat-section {
  background: white;
  margin-bottom: 16px;
}

.chat-header {
  padding: 16px;
  text-align: center;
  border-bottom: 1px solid #ebedf0;
}

.chat-header h2 {
  font-size: 18px;
  font-weight: bold;
  color: #323233;
  margin: 0 0 4px 0;
}

.subtitle {
  font-size: 13px;
  color: #969799;
  margin: 0;
}

.chat-messages {
  height: 300px;
  overflow-y: auto;
  padding: 16px;
  background: #f7f8fa;
}

.message {
  display: flex;
  margin-bottom: 16px;
}

.message.user {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #e8f3ff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: #1989fa;
}

.message.user .message-avatar {
  background: #1989fa;
  color: white;
}

.message-content {
  max-width: 70%;
  margin: 0 8px;
}

.message.user .message-content {
  display: flex;
  justify-content: flex-end;
}

.message-text {
  padding: 8px 12px;
  background: white;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.5;
  color: #323233;
}

.message.user .message-text {
  background: #1989fa;
  color: white;
}

.message-text :deep(code) {
  padding: 2px 6px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
  font-size: 12px;
}

.message-data {
  margin-top: 8px;
}

.data-chart {
  padding: 12px;
  background: white;
  border-radius: 8px;
}

.chart-title {
  font-size: 14px;
  font-weight: bold;
  color: #323233;
  margin-bottom: 12px;
}

.chart-bars {
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
  height: 120px;
}

.chart-bar-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.chart-bar {
  width: 20px;
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px 4px 0 0;
}

.chart-label {
  font-size: 10px;
  color: #969799;
  margin-top: 4px;
}

.data-table {
  padding: 12px;
  background: white;
  border-radius: 8px;
}

.table-row {
  padding: 8px 0;
  border-bottom: 1px solid #ebedf0;
}

.table-row:last-child {
  border-bottom: none;
}

.table-cell {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
  font-size: 13px;
}

.cell-label {
  color: #969799;
  flex-shrink: 0;
}

.cell-value {
  color: #323233;
  text-align: right;
}

.chat-input {
  padding: 12px 16px;
  background: white;
  border-top: 1px solid #ebedf0;
}

.quick-questions {
  background: white;
  padding: 16px;
  margin-bottom: 16px;
}

.section-title {
  font-size: 16px;
  font-weight: bold;
  color: #323233;
  margin-bottom: 12px;
}

.questions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.insights-section {
  background: white;
  padding: 16px;
  margin-bottom: 16px;
}

.insights-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.insight-card {
  display: flex;
  align-items: center;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
}

.insight-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
}

.insight-info {
  flex: 1;
}

.insight-title {
  font-size: 12px;
  color: #969799;
  margin: 0 0 4px 0;
}

.insight-value {
  font-size: 16px;
  font-weight: bold;
  color: #323233;
  margin: 0 0 4px 0;
}

.insight-trend {
  font-size: 12px;
  margin: 0;
}

.insight-trend.up {
  color: #07c160;
}

.insight-trend.down {
  color: #ee0a24;
}

.insight-trend.warning {
  color: #ff976a;
}

.insight-trend.info {
  color: #1989fa;
}

.recommendations-section {
  background: white;
  padding: 16px;
}
</style>
