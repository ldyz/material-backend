# 材料管理系统 - Vue 3 重构版

## 项目简介

这是使用 Vue 3 + Vite + Element Plus 重构的材料管理系统前端项目。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **Vite** - 新一代前端构建工具
- **Vue Router 4** - Vue.js 官方路由
- **Pinia** - Vue 的状态管理库
- **Element Plus** - 基于 Vue 3 的组件库
- **Axios** - HTTP 客户端

## 项目结构

```
newstatic/
├── src/
│   ├── api/              # API 接口
│   │   ├── index.js      # API 模块
│   │   └── request.js    # Axios 封装
│   ├── assets/           # 静态资源
│   │   ├── css/          # 样式文件
│   │   └── images/       # 图片资源
│   ├── components/       # 组件
│   │   ├── common/       # 通用组件
│   │   └── layout/       # 布局组件
│   ├── composables/      # 组合式函数
│   ├── router/           # 路由配置
│   ├── stores/           # Pinia 状态管理
│   ├── utils/            # 工具函数
│   ├── views/            # 页面视图
│   ├── App.vue           # 根组件
│   └── main.js           # 入口文件
├── index.html            # HTML 模板
├── vite.config.js        # Vite 配置
└── package.json          # 项目依赖

```

## 快速开始

### 安装依赖

```bash
cd newstatic
npm install
```

### 开发模式

```bash
npm run dev
```

开发服务器将在 `http://localhost:3000` 启动

### 构建生产版本

```bash
npm run build
```

构建后的文件将输出到 `../static/` 目录

## 功能模块

- [x] 用户认证（登录/登出）
- [x] 仪表板
- [ ] 项目管理
- [ ] 施工日志
- [ ] 进度管理
- [ ] 物资管理
- [ ] 库存管理
- [ ] 出库单管理
- [ ] 入库单管理
- [ ] AI 数据分析
- [ ] 系统管理

## 开发说明

### API 代理配置

在 `vite.config.js` 中已配置 API 代理：

```javascript
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

### 权限控制

使用 Pinia store (`useAuthStore`) 管理用户认证和权限：

```javascript
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// 检查权限
if (authStore.hasPermission('project_view')) {
  // 有权限
}
```

### 路由守卫

在 `src/router/index.js` 中配置了路由守卫，自动处理：
- 认证检查
- 权限验证
- 页面标题设置

## 部署

生产环境构建会自动输出到 `backend/static/` 目录，Go 后端可以直接使用。

```bash
npm run build
```

## 注意事项

1. 确保后端 API 服务运行在 `http://localhost:8080`
2. 登录后 token 会保存在 localStorage 中
3. 登录状态有效期 72 小时
4. 所有 API 请求都会自动携带 Authorization 头

## 下一步计划

- [ ] 完善各个功能模块
- [ ] 集成 Quill 富文本编辑器
- [ ] 集成 Chart.js 图表
- [ ] 添加单元测试
- [ ] 优化性能
- [ ] 完善错误处理
