# 贡献指南

感谢您对甘特图组件的关注！我们欢迎各种形式的贡献。

## 目录

- [行为准则](#行为准则)
- [开发环境设置](#开发环境设置)
- [代码规范](#代码规范)
- [提交代码](#提交代码)
- [问题报告](#问题报告)
- [功能建议](#功能建议)

---

## 行为准则

### 我们的承诺

为了营造开放和友好的环境，我们承诺：

- 🤝 尊重不同的观点和经验
- 🎉 使用欢迎和包容的语言
- 🙏 感恩他人的贡献
- 👀 专注于对社区最有利的事情
- 😊 保持友善和耐心

### 不可接受的行为

- 使用性化或物化的语言或图像
- 人身攻击或侮辱性言论
- 公开或私下骚扰
- 未经许可发布他人私人信息
- 其他不道德或不专业的行为

---

## 开发环境设置

### 1. Fork 并克隆仓库

```bash
# Fork 仓库
git clone https://github.com/your-username/gantt-chart.git
cd gantt-chart

# 添加上游仓库
git remote add upstream https://github.com/original-org/gantt-chart.git
```

### 2. 安装依赖

```bash
# 安装 Node.js 依赖
npm install

# 安装开发工具
npm install -g commitizen
npm install -g @commitlint/cli
```

### 3. 启动开发服务器

```bash
# 启动开发环境
npm run dev

# 运行类型检查
npm run type-check

# 运行代码检查
npm run lint
```

---

## 代码规范

### TypeScript 规范

```typescript
// ✅ 好的实践
interface TaskProps {
  id: string
  name: string
  onComplete?: () => void
}

// ❌ 避免
const data: any = {}  // 不要使用 any
```

### Vue 3 规范

```vue
<!-- ✅ 使用 Script Setup -->
<script setup lang="ts">
import { ref } from 'vue'

const count = ref(0)
</script>

<!-- ❌ 避免 Options API（新代码） -->
<script>
export default {
  data() {
    return { count: 0 }
  }
}
</script>
```

### 命名规范

```typescript
// 组件：PascalCase
TaskBar.vue
GanttChart.vue

// Composable：camelCase，use 前缀
useGanttDrag.ts
useTheme.ts

// 工具函数：camelCase
formatDate()
calculateDuration()

// 常量：UPPER_SNAKE_CASE
const MAX_TASKS = 1000
const DEFAULT_ZOOM = 40

// 类型：PascalCase，I 后缀
interface GanttTask
enum TaskStatus
```

### 注释规范

```typescript
/**
 * 计算关键路径
 * @param tasks - 任务列表
 * @returns 关键路径上的任务ID数组
 * @throws {GanttError} 当存在循环依赖时
 */
function calculateCriticalPath(tasks: GanttTask[]): string[] {
  // 实现...
}
```

---

## 提交代码

### 1. 创建分支

```bash
# 从 main 分支创建特性分支
git checkout -b feature/my-awesome-feature

# 或从 issue 创建分支
git checkout -b issue-123-fix-bug
```

### 2. 编写代码

- 遵循代码规范
- 添加必要的注释
- 更新相关文档
- 确保测试通过

### 3. 提交更改

```bash
# 添加更改
git add .

# 使用 Commitizen 提交
npx git-cz

# 或手动提交（遵循约定式提交）
git commit -m "feat: add virtual scrolling support"
```

### 提交信息格式

遵循约定式提交规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型 (type):**

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建/工具更新

**示例：**

```
feat(timeline): add virtual scrolling for large datasets

- Implement useVirtualScroll composable
- Add VirtualList component
- Optimize rendering for 10000+ tasks

Closes #123
```

### 4. 推送更改

```bash
# 推送到你的 fork
git push origin feature/my-awesome-feature

# 创建 Pull Request
```

### 5. 创建 Pull Request

1. 在 GitHub 上打开原仓库
2. 点击 "New Pull Request"
3. 选择你的分支
4. 填写 PR 模板

**PR 模板：**

```markdown
## 描述
简要描述此 PR 的目的和改动。

## 改动类型
- [ ] Bug 修复
- [ ] 新功能
- [ ] 重构
- [ ] 文档
- [ ] 性能优化
- [ ] 其他

## 测试
- [ ] 单元测试通过
- [ ] E2E 测试通过
- [ ] 手动测试通过

## 截图
（如有 UI 改动，请提供截图）

## Checklist
- [ ] 遵循代码规范
- [ ] 更新了相关文档
- [ ] 添加了必要的测试
- [ ] 所有测试通过
```

---

## 问题报告

### 报告 Bug

请在 Issue 中包含以下信息：

**Bug 描述**
简要描述遇到的问题。

**复现步骤**
1. 访问页面 '...'
2. 点击按钮 '....'
3. 滚动到 '....'
4. 看到错误

**期望行为**
描述你期望发生什么。

**实际行为**
描述实际发生了什么。

**环境信息**
- OS: [e.g. Windows 11, macOS 14]
- Browser: [e.g. Chrome 120, Firefox 121]
- Vue Version: [e.g. 3.4.0]
- Component Version: [e.g. 2.0.0]

**截图**

如果可能，添加截图来展示问题。

---

## 功能建议

### 建议新功能

请先在 Issues 中讨论功能提案，我们会在实现前确认：

1. 功能是否符合项目目标
2. 实现难度和优先级
3. 对现有功能的影响

**建议模板：**

```markdown
## 功能描述
简要描述建议的功能。

## 使用场景
描述这个功能解决什么问题，谁会受益。

## 实现建议（可选）
如果你有技术实现的想法，请分享。

## 替代方案（可选）
是否有其他方式达到类似效果？
```

---

## 代码审查流程

### 审查重点

1. **代码质量**
   - 遵循代码规范
   - 无 lint 错误
   - 合理的复杂度

2. **功能正确性**
   - 实现符合需求
   - 边界情况处理
   - 错误处理

3. **性能**
   - 无明显性能问题
   - 大数据量测试通过
   - 无内存泄漏

4. **测试**
   - 测试覆盖率充足
   - 测试用例合理

5. **文档**
   - API 文档更新
   - 变更日志更新
   - 必要的注释添加

### 审查反馈

审查者会在 PR 中留下评论，请：

- ✅ 积极响应反馈
- ✅ 及时修改问题
- ✅ 保持友好沟通

---

## 发布流程

### 版本号规则

遵循语义化版本 `MAJOR.MINOR.PATCH`：

- **MAJOR**: 破坏性变更
- **MINOR**: 新功能（向后兼容）
- **PATCH**: Bug 修复

### 发布检查清单

- [ ] 所有测试通过
- [ ] 文档已更新
- [ ] CHANGELOG 已更新
- [ ] 版本号已更新
- [ ] 标签已创建

---

## 获取帮助

### 沟道

- **Issues**: 报告 Bug 和功能建议
- **Discussions**: 技术讨论和疑问
- **Email**: support@example.com

### 社区资源

- [API 文档](./API.md)
- [迁移指南](./MIGRATION.md)
- [架构文档](./ARCHITECTURE.md)

---

## 认可贡献者

您的名字将会出现在项目的贡献者列表中！

---

**再次感谢您的贡献！**

---

**更新日期**: 2026-01-31
**版本**: 2.0.0
