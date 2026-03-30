# 移动端前端文档

本文档详细说明了移动端前端项目的结构、组件和开发指南。

## 项目概述

- **框架**: Vue 3 + Composition API
- **UI组件库**: Vant 4
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **构建工具**: Vite
- **原生桥接**: Capacitor

## 目录结构

```
mobile-app/
├── public/                 # 静态资源
│   └── version.json        # 版本信息
├── src/
│   ├── api/               # API接口定义
│   │   ├── index.js       # 统一API导出
│   │   ├── appointment.js # 预约API
│   │   ├── attendance.js  # 考勤API
│   │   ├── inbound.js     # 入库API
│   │   ├── requisition.js # 出库API
│   │   ├── material_plan.js # 物资计划API
│   │   └── agent.js       # AI助手API
│   ├── assets/            # 静态资源
│   ├── components/        # 公共组件
│   │   ├── common/        # 通用组件
│   │   │   ├── ListItemCard.vue
│   │   │   ├── ProjectSelector.vue
│   │   │   └── WorkerPicker.vue
│   │   └── AiChatPopup.vue # AI对话弹窗
│   ├── layouts/           # 布局组件
│   │   └── TabbarLayout.vue
│   ├── router/            # 路由配置
│   │   └── index.js
│   ├── stores/            # Pinia状态管理
│   │   ├── auth.js        # 认证状态
│   │   └── notification.js # 通知状态
│   ├── utils/             # 工具函数
│   │   ├── websocket.js   # WebSocket连接
│   │   └── index.js
│   ├── views/             # 页面组件
│   │   ├── Login/         # 登录
│   │   ├── Dashboard/     # 仪表板
│   │   ├── Plan/          # 物资计划
│   │   │   ├── List.vue
│   │   │   ├── Detail.vue
│   │   │   ├── Approve.vue
│   │   │   └── Create.vue
│   │   ├── Inbound/       # 入库管理
│   │   │   ├── List.vue
│   │   │   ├── Detail.vue
│   │   │   ├── Approve.vue
│   │   │   └── Create.vue
│   │   ├── Requisition/   # 出库管理
│   │   │   ├── List.vue
│   │   │   ├── Detail.vue
│   │   │   ├── Approve.vue
│   │   │   └── Create.vue
│   │   ├── Appointment/   # 施工预约
│   │   │   ├── List.vue
│   │   │   ├── Detail.vue
│   │   │   ├── Approve.vue
│   │   │   ├── Create.vue
│   │   │   └── Calendar.vue
│   │   ├── Attendance/    # 考勤打卡
│   │   │   ├── ClockIn.vue
│   │   │   └── RecordList.vue
│   │   ├── Notification/  # 通知
│   │   │   └── List.vue
│   │   └── Profile/       # 个人中心
│   │       └── index.vue
│   ├── App.vue
│   └── main.js
├── android/               # Android原生项目
│   └── app/
│       ├── build.gradle   # 构建配置
│       └── src/main/
│           └── AndroidManifest.xml
├── dist-capacitor/        # Capacitor构建输出
├── capacitor.config.json  # Capacitor配置
├── index.html
├── vite.config.js
├── package.json
└── build.sh               # 构建脚本
```

## 页面清单

| 页面 | 路由 | 功能 |
|------|------|------|
| Login | /login | 用户登录 |
| Dashboard | / | 仪表板/首页 |
| PlanList | /plans | 物资计划列表 |
| PlanDetail | /plans/:id | 计划详情 |
| PlanApprove | /plans/:id/approve | 计划审批 |
| PlanCreate | /plans/create | 创建计划 |
| InboundList | /inbound | 入库单列表 |
| InboundDetail | /inbound/:id | 入库单详情 |
| InboundApprove | /inbound/:id/approve | 入库审批 |
| InboundCreate | /inbound/create | 创建入库 |
| RequisitionList | /requisition | 出库单列表 |
| RequisitionDetail | /requisition/:id | 出库单详情 |
| RequisitionApprove | /requisition/:id/approve | 出库审批 |
| RequisitionCreate | /requisition/create | 创建出库 |
| AppointmentList | /appointments | 预约列表 |
| AppointmentCalendar | /appointments/calendar | 预约日历 |
| AppointmentDetail | /appointment/:id | 预约详情 |
| AppointmentCreate | /appointment/create | 创建预约 |
| AppointmentEdit | /appointment/:id/edit | 编辑预约 |
| AppointmentApprove | /appointment/:id/approve | 预约审批 |
| AttendanceClockIn | /attendance/clock-in | 考勤打卡 |
| AttendanceRecords | /attendance/records | 打卡记录 |
| NotificationList | /notifications | 通知列表 |
| Profile | /profile | 个人中心 |

## 核心组件

### TabbarLayout 底部导航布局

**路径**: `src/layouts/TabbarLayout.vue`

**功能**:
- 底部Tab导航
- 页面缓存支持
- 通知未读数显示

**Tab配置**:
| 图标 | 标题 | 路由 |
|------|------|------|
| home-o | 首页 | / |
| orders-o | 计划 | /plans |
| sign | 打卡 | /attendance/clock-in |
| bell | 消息 | /notifications |
| user-o | 我的 | /profile |

### ListItemCard 列表项卡片

**路径**: `src/components/common/ListItemCard.vue`

**功能**:
- 统一的列表项展示
- 支持标题、状态、时间等信息
- 点击跳转详情

**Props**:
| 名称 | 类型 | 默认值 | 说明 |
|------|------|-------|------|
| title | String | '' | 标题 |
| status | String | '' | 状态 |
| statusType | String | 'primary' | 状态类型 |
| subtitle | String | '' | 副标题 |
| rightText | String | '' | 右侧文本 |
| time | String | '' | 时间 |

**使用示例**:
```vue
<ListItemCard
  :title="item.name"
  :status="item.status"
  status-type="success"
  :time="item.created_at"
  @click="goDetail(item.id)"
/>
```

### WorkerPicker 作业人员选择器

**路径**: `src/components/common/WorkerPicker.vue`

**功能**:
- 多选作业人员
- 显示可用状态
- 支持搜索

**Props**:
| 名称 | 类型 | 默认值 | 说明 |
|------|------|-------|------|
| modelValue | Array | [] | 已选人员ID列表 |
| workDate | String | '' | 作业日期 |
| timeSlot | String | '' | 时间段 |

**Events**:
| 名称 | 参数 | 说明 |
|------|------|------|
| update:modelValue | Array | 选择变化 |

### AiChatPopup AI对话弹窗

**路径**: `src/components/AiChatPopup.vue`

**功能**:
- 文本对话
- 语音输入
- 快捷操作

**Props**:
| 名称 | 类型 | 默认值 | 说明 |
|------|------|-------|------|
| show | Boolean | false | 是否显示 |

## 状态管理

### auth.js 认证状态

```javascript
{
  user: null,
  token: null,
  isAuthenticated: false,
  permissions: [],
  roles: []
}
```

### notification.js 通知状态

```javascript
{
  unreadCount: 0,
  notifications: [],
  wsConnected: false
}
```

**Actions**:
- `connectWebSocket()` - 连接WebSocket
- `disconnectWebSocket()` - 断开连接
- `fetchNotifications()` - 获取通知列表
- `markAsRead(id)` - 标记已读

## API调用层

### appointment.js 预约API

```javascript
// 获取预约列表
export function getAppointments(params)

// 获取我的预约
export function getMyAppointments(params)

// 获取预约详情
export function getAppointment(id)

// 创建预约
export function createAppointment(data)

// 更新预约
export function updateAppointment(id, data)

// 提交审批
export function submitAppointment(id)

// 审批预约
export function approveAppointment(id, data)

// 分配作业人员
export function assignWorkers(id, data)

// 获取作业人员列表
export function getWorkers(params)

// 获取每日统计
export function getDailyStatistics(params)
```

### attendance.js 考勤API

```javascript
// 打卡
export function clockIn(data)

// 获取打卡记录
export function getAttendanceRecords(params)

// 获取今日待打卡任务
export function getTodayAppointments()

// 确认记录
export function confirmRecord(id, data)

// 驳回记录
export function rejectRecord(id, data)
```

### agent.js AI助手API

```javascript
// 文本对话
export function chat(data)

// 语音对话
export function voiceChat(formData)

// 获取对话历史
export function getConversationHistory()
```

## Capacitor配置

### capacitor.config.json

```json
{
  "appId": "com.material.management",
  "appName": "材料管理",
  "webDir": "dist-capacitor",
  "server": {
    "androidScheme": "https"
  },
  "plugins": {
    "Camera": {
      "permissions": ["camera", "photos"]
    },
    "Geolocation": {
      "permissions": ["location"]
    }
  }
}
```

### 原生权限配置

**AndroidManifest.xml**:
```xml
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.CAMERA" />
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />
<uses-permission android:name="android.permission.RECORD_AUDIO" />
```

## 原生功能使用

### 相机拍照

```javascript
import { Camera, CameraResultType } from '@capacitor/camera'

async function takePhoto() {
  const photo = await Camera.getPhoto({
    quality: 90,
    allowEditing: false,
    resultType: CameraResultType.Uri
  })
  return photo.webPath
}
```

### 获取位置

```javascript
import { Geolocation } from '@capacitor/geolocation'

async function getLocation() {
  const position = await Geolocation.getCurrentPosition()
  return {
    latitude: position.coords.latitude,
    longitude: position.coords.longitude
  }
}
```

### 语音录制

```javascript
import { VoiceRecorder } from 'capacitor-voice-recorder'

async function recordAudio() {
  const { value: record } = await VoiceRecorder.startRecording()
  // ...
  const { value: audio } = await VoiceRecorder.stopRecording()
  return audio
}
```

## 开发命令

```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 构建Capacitor版本（必须设置环境变量）
CAPACITOR_BUILD=true npm run build

# 同步到Android
npx cap sync android

# 打开Android Studio
npx cap open android

# 构建APK（使用脚本）
./build.sh
```

## 构建脚本

`build.sh` 脚本内容：

```bash
#!/bin/bash
echo "Building mobile app..."

# 设置环境变量
export CAPACITOR_BUILD=true

# 构建前端
npm run build

# 同步到Android
npx cap sync android

# 构建APK
cd android
./gradlew assembleDebug

# 复制APK
cp app/build/outputs/apk/debug/app-debug.apk ../material-management.apk

echo "Build completed!"
echo "APK: material-management.apk"
```

## 版本更新

### 更新流程

1. 更新版本号
```json
// public/version.json
{
  "version": "1.0.101",
  "buildNumber": 101
}
```

2. 更新Android配置
```gradle
// android/app/build.gradle
android {
  defaultConfig {
    versionCode 101
    versionName "1.0.101"
  }
}
```

3. 构建并上传
```bash
./build.sh
# 将APK上传到 mobile-app-updates/android/
```

4. 创建版本信息文件
```json
// mobile-app-updates/android/version-1.0.101.json
{
  "platform": "android",
  "version": "1.0.101",
  "build_number": 101,
  "download_url": "https://xxx/mobile-updates/android/material-management-1.0.101.apk",
  "force_update": false,
  "update_message": "修复若干问题",
  "release_notes": "..."
}
```

### 热更新检查

应用启动时会检查版本更新：

```javascript
async function checkUpdate() {
  const response = await fetch('/api/app/version')
  const latest = await response.json()

  if (latest.build_number > currentBuildNumber) {
    // 显示更新提示
    showUpdateDialog(latest)
  }
}
```

## 样式规范

- 使用Vant默认主题
- 自定义样式在 `src/assets/styles/`
- 适配方案：使用 `postcss-px-to-viewport` 自动转换

## 常见问题

### 1. 白屏问题

确保构建时设置了 `CAPACITOR_BUILD=true` 环境变量。

### 2. 权限问题

在AndroidManifest.xml中添加所需权限。

### 3. 网络请求失败

检查API地址配置，确保使用正确的服务器地址。

### 4. 推送通知

需要配置Firebase Cloud Messaging (FCM)。
