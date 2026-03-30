# 移动端项目优化建议

## 生成时间
2026-03-30

---

## 一、代码架构优化

### 1. API 请求方式不统一
**问题描述**：
- `AiChatPopup.vue` 中使用原生 `fetch` 而非封装好的 `request` 实例
- 需要手动处理 Capacitor 环境的 URL 判断

**建议**：
- 创建统一的 API 调用方式，所有请求都通过 `src/utils/request.js`
- 或创建一个专门的 `agent.js` API 文件，复用 request 配置

**涉及文件**：
- `src/components/AiChatPopup.vue` (已部分修复)
- `src/api/agent.js`

---

### 2. 代码重复问题

**问题描述**：
- 多个视图组件（Plan、Inbound、Requisition、Appointment）的列表、详情、审批页面结构高度相似
- 状态标签、时间格式化等工具函数在各组件中重复定义

**建议**：
- 抽取通用的 `BaseList.vue`、`BaseDetail.vue`、`BaseApprove.vue` 组件
- 将通用工具函数移至 `src/utils/` 目录
- 使用 `src/composables/` 目录下的组合式函数复用逻辑

**示例**：
```javascript
// 已存在但未充分利用的 composables
src/composables/useStatus.js      // 状态相关
src/composables/useDateTime.js    // 时间相关
src/composables/useApprovalWorkflow.js // 审批流程
```

---

### 3. 控制台日志过多
**问题描述**：
- 发现 146 处 `console.log` 调用
- 生产环境不应包含大量调试日志

**建议**：
- 创建统一的日志工具 `src/utils/logger.js`
- 根据环境变量控制日志输出
- 示例实现：
```javascript
// src/utils/logger.js
const isDev = import.meta.env.DEV
export const logger = {
  log: (...args) => isDev && console.log(...args),
  error: (...args) => console.error(...args),
  warn: (...args) => isDev && console.warn(...args)
}
```

---

## 二、性能优化

### 1. 数据缓存机制缺失
**问题描述**：
- 每次进入页面都重新请求数据
- 没有数据缓存策略，用户体验不佳

**建议**：
- 添加 API 响应缓存机制
- 对不常变化的数据（如项目列表、作业人员列表）进行本地缓存
- 可使用 `localStorage` 或内存缓存

**实现方案**：
```javascript
// src/utils/cache.js
const cache = new Map()
const CACHE_TTL = 5 * 60 * 1000 // 5分钟

export function getCachedData(key, fetcher) {
  const cached = cache.get(key)
  if (cached && Date.now() - cached.time < CACHE_TTL) {
    return Promise.resolve(cached.data)
  }
  return fetcher().then(data => {
    cache.set(key, { data, time: Date.now() })
    return data
  })
}
```

---

### 2. 图片懒加载
**问题描述**：
- 用户头像等图片没有懒加载机制
- 可能影响列表滚动性能

**建议**：
- 使用 Vant 的 `van-image` 组件的 `lazy-load` 属性
- 或实现自定义图片懒加载指令

---

### 3. WebSocket 重连策略优化
**问题描述**：
- 当前最多重试 5 次后放弃
- 没有指数退避策略

**建议**：
```javascript
// 当前实现
const delay = this.reconnectDelay * this.reconnectAttempts

// 建议改为指数退避
const delay = Math.min(
  this.reconnectDelay * Math.pow(2, this.reconnectAttempts),
  30000 // 最大 30 秒
)
```

---

## 三、用户体验优化

### 1. 离线功能支持
**问题描述**：
- 当前应用完全依赖网络
- 网络不佳时用户体验差

**建议**：
- 添加离线状态检测
- 对关键操作提供离线队列
- 使用 Service Worker 缓存静态资源

---

### 2. 骨架屏缺失
**问题描述**：
- 页面加载时只有 loading 状态
- 用户感知等待时间较长

**建议**：
- 为列表页和详情页添加骨架屏组件
- 使用 Vant 的 `van-skeleton` 组件

---

### 3. 错误处理不统一
**问题描述**：
- 各组件错误处理方式不一致
- 部分错误没有用户友好提示

**建议**：
- 创建全局错误处理器
- 统一错误提示样式和内容

---

## 四、代码质量优化

### 1. localStorage 使用分散
**问题描述**：
- 发现 37 处直接使用 localStorage
- 已有 `src/utils/storage.js` 封装但未完全使用

**建议**：
- 统一使用 `storage.js` 封装的方法
- 便于后续迁移到其他存储方案

---

### 2. 类型安全
**问题描述**：
- 项目使用纯 JavaScript
- 缺少类型检查，容易出错

**建议**：
- 考虑迁移到 TypeScript
- 或至少添加 JSDoc 注释提供类型提示

---

### 3. 组件命名规范
**问题描述**：
- 部分组件命名不一致
- 如 `Dashboard/index.vue` 和 `Dashboard.vue` 同时存在

**建议**：
- 统一组件命名规范
- 删除冗余文件

---

## 五、语音识别本地化研究

### 当前方案分析
当前项目使用 `@independo/capacitor-voice-recorder` 进行录音，然后将音频发送到后端进行语音识别。这种方式：
- **优点**：识别准确度高，支持复杂场景
- **缺点**：依赖网络，延迟较高，服务器压力大

### 可选的本地语音识别方案

#### 方案一：Android 原生 SpeechRecognizer API
**平台**：仅 Android
**离线支持**：部分支持（需提前下载语言包）
**实现方式**：
```javascript
// 使用 Capacitor 插件调用原生 API
// 可参考：capacitor-voice-recognition
import { VoiceRecognition } from 'capacitor-voice-recognition'

// 开始识别
await VoiceRecognition.start({
  language: 'zh-CN',
  partialResults: true // 实时返回部分结果
})

// 监听结果
VoiceRecognition.addListener('partialResults', (result) => {
  console.log('识别中:', result.matches)
})
```

**优点**：
- 免费使用系统内置能力
- 支持实时语音识别
- 低延迟

**缺点**：
- 需要 Android 5.0+ 且安装 Google 服务
- 离线识别准确度较低
- iOS 需要单独实现

#### 方案二：iOS 原生 SFSpeechRecognizer
**平台**：仅 iOS
**离线支持**：部分支持
**实现方式**：需要通过 Capacitor 插件调用

#### 方案三：Web Speech API（已有部分使用）
**平台**：浏览器环境
**离线支持**：不支持
**当前状态**：已在项目中用于语音播报

#### 方案四：Vosk（完全离线）
**平台**：跨平台
**离线支持**：完全支持
**特点**：
- 开源语音识别引擎
- 可集成到移动应用中
- 中文模型约 50MB

**实现复杂度**：高，需要：
1. 编写原生插件
2. 集成 Vosk 库
3. 下载并打包语音模型

### 推荐方案

**短期方案**：
1. 使用 `capacitor-voice-recognition` 插件
2. Android 使用原生 SpeechRecognizer
3. iOS 使用原生 SFSpeechRecognizer
4. 首选本地识别，失败时回退到后端识别

**实现建议**：
```javascript
// src/utils/speechRecognition.js
import { Capacitor } from '@capacitor/core'

class SpeechRecognitionService {
  async startRecognition(options = {}) {
    if (Capacitor.isNativePlatform()) {
      // 尝试使用原生识别
      try {
        return await this.nativeRecognition(options)
      } catch (e) {
        // 原生识别失败，回退到后端
        console.log('Native recognition failed, fallback to server')
        return this.serverRecognition(options)
      }
    } else {
      // Web 环境使用后端识别
      return this.serverRecognition(options)
    }
  }

  async nativeRecognition(options) {
    // 使用原生 API
  }

  async serverRecognition(options) {
    // 使用当前的后端 WebSocket 方案
  }
}
```

**相关 Capacitor 插件**：
1. `capacitor-voice-recognition` - 调用原生语音识别
2. `@capacitor-community/speech-recognition` - 社区维护版本

---

## 六、优先级排序

### 高优先级（建议立即处理）
1. 统一 API 请求方式
2. 减少生产环境的 console.log
3. 添加数据缓存机制

### 中优先级（可安排迭代）
1. 抽取通用组件
2. 统一错误处理
3. 添加骨架屏

### 低优先级（长期优化）
1. 迁移到 TypeScript
2. 实现本地语音识别
3. 添加离线功能支持

---

## 七、文件清理建议

以下文件可能需要清理或整合：
- `src/views/Dashboard.vue` 和 `src/views/Dashboard/index.vue` 存在重复
- `src/api/` 目录下的部分 API 文件可以合并

---

*以上建议仅供参考，实施前请评估影响范围和测试充分性。*
