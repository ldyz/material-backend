# 开发指南

本文档说明了开发环境的搭建、开发流程和代码规范。

## 环境要求

### 后端环境

| 工具 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.21+ | 后端语言 |
| PostgreSQL | 14+ | 数据库 |
| Git | 2.x | 版本控制 |

### 前端环境

| 工具 | 版本要求 | 说明 |
|------|---------|------|
| Node.js | 18+ | JavaScript运行时 |
| npm/pnpm | 9+/8+ | 包管理器 |

### 移动端环境

| 工具 | 版本要求 | 说明 |
|------|---------|------|
| Android Studio | 最新版 | Android开发 |
| JDK | 17+ | Java开发工具包 |
| Android SDK | API 33+ | Android SDK |

## 环境搭建

### 1. 克隆代码

```bash
git clone <repository-url>
cd backend
```

### 2. 后端环境

```bash
# 安装Go依赖
go mod download

# 创建配置文件
cp config.example.yaml config.yaml

# 编辑配置文件
vim config.yaml
```

### 3. 数据库初始化

```bash
# 创建数据库
createdb -U postgres materials

# 执行迁移脚本
for f in migrations/*.sql; do
  psql -U materials -d materials -f "$f"
done
```

### 4. 前端环境

```bash
# Web前端
cd newstatic
npm install

# 移动端
cd ../mobile-app
npm install
```

### 5. 启动开发服务

```bash
# 后端
./rebuild.sh

# Web前端
cd newstatic
npm run dev

# 移动端
cd mobile-app
npm run dev
```

## 项目配置

### 后端配置 (config.yaml)

```yaml
server:
  port: 8088           # 服务端口
  mode: debug          # 运行模式: debug/release

database:
  type: postgresql
  host: 127.0.0.1
  port: 5432
  user: materials
  password: your_password
  database: materials

jwt:
  secret: "your-secret-key-change-in-production"
  expire_time: 24h

upload:
  max_file_size: 5242880  # 5MB
  upload_dir: "static/uploads"

ai:
  provider: baidu        # baidu 或 deepseek
  baidu_api_key: ""
  baidu_secret_key: ""
  deepseek_api_key: ""
```

### 前端配置

**Web前端 (.env.development)**:
```
VITE_API_BASE_URL=/api
```

**移动端 (.env.development)**:
```
VITE_API_BASE_URL=http://192.168.1.100:8088/api
```

## 开发流程

### 1. 创建功能分支

```bash
git checkout -b feature/your-feature
```

### 2. 开发代码

遵循代码规范进行开发。

### 3. 本地测试

```bash
# 后端测试
go test ./...

# 前端测试
npm run test
```

### 4. 提交代码

```bash
git add .
git commit -m "feat: 添加新功能描述"
```

### 5. 合并代码

```bash
git checkout main
git merge feature/your-feature
```

## 代码规范

### Go代码规范

#### 命名规范

```go
// 包名：小写单词
package auth

// 导出函数：驼峰命名
func GetUserByID(id uint) (*User, error) { }

// 私有函数：小写开头
func validateUser(u *User) error { }

// 常量：驼峰命名
const (
    StatusActive = "active"
    StatusClosed = "closed"
)

// 接口：动词+er
type UserRepositorier interface {
    FindByID(id uint) (*User, error)
}
```

#### 错误处理

```go
// 使用 errors.New 或 fmt.Errorf
if err != nil {
    return fmt.Errorf("查询用户失败: %w", err)
}

// 统一响应格式
if err := db.First(&user, id).Error; err != nil {
    response.NotFound(c, "用户不存在")
    return
}
```

#### 注释规范

```go
// GetUserByID 根据ID获取用户
// 参数:
//   - id: 用户ID
// 返回:
//   - *User: 用户对象
//   - error: 错误信息
func GetUserByID(id uint) (*User, error) {
    // 实现...
}
```

### Vue代码规范

#### 组件命名

```vue
<!-- 使用PascalCase -->
<template>
  <UserSelector v-model="userId" />
</template>

<script setup>
import UserSelector from '@/components/UserSelector.vue'
</script>
```

#### Props定义

```javascript
const props = defineProps({
  modelValue: {
    type: Number,
    required: true
  },
  disabled: {
    type: Boolean,
    default: false
  }
})
```

#### 事件定义

```javascript
const emit = defineEmits(['update:modelValue', 'change'])

function handleChange(value) {
  emit('update:modelValue', value)
  emit('change', value)
}
```

#### API调用

```javascript
// 使用 async/await
async function fetchData() {
  try {
    loading.value = true
    const res = await api.getData(params)
    data.value = res.data
  } catch (error) {
    ElMessage.error(error.message)
  } finally {
    loading.value = false
  }
}
```

### Git提交规范

使用约定式提交：

```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式（不影响功能）
refactor: 重构
perf: 性能优化
test: 测试
chore: 构建/工具变更
```

示例：
```
feat: 添加用户头像上传功能
fix: 修复登录Token过期问题
docs: 更新API文档
refactor: 重构权限检查逻辑
```

## 调试技巧

### 后端调试

#### 日志输出

```go
import "log"

// 使用标准日志
log.Printf("[DEBUG] 用户ID: %d", userID)

// 查看日志
tail -f /tmp/backend.log
```

#### 数据库调试

```go
// 启用GORM日志
db.Debug().Find(&users)
```

#### 性能分析

```go
import "net/http/pprof"

// 在开发模式启用pprof
go func() {
    http.ListenAndServe(":6060", nil)
}()

// 访问分析
// http://localhost:6060/debug/pprof/
```

### 前端调试

#### Vue Devtools

安装Vue Devtools浏览器扩展，可以查看组件树、状态、事件等。

#### 网络请求

使用浏览器开发者工具的Network面板查看API请求。

#### Console调试

```javascript
// 打印响应式数据
console.log('data:', JSON.stringify(data.value))

// 使用debugger断点
debugger
```

## 常用命令速查

### 后端

```bash
# 编译
go build -o server ./cmd/server

# 编译并重启
./rebuild.sh

# 运行测试
go test ./...

# 查看日志
tail -f /tmp/backend.log

# 数据库迁移
psql -U materials -d materials -f migrations/xxx.sql
```

### Web前端

```bash
# 开发模式
npm run dev

# 构建
npm run build

# 代码检查
npm run lint

# 格式化
npm run format
```

### 移动端

```bash
# 开发模式
npm run dev

# 构建
CAPACITOR_BUILD=true npm run build

# 同步Capacitor
npx cap sync

# 打开Android Studio
npx cap open android

# 构建APK
./build.sh
```

## 故障排查

### 后端无法启动

1. 检查配置文件是否正确
2. 检查数据库连接
3. 检查端口是否被占用

```bash
# 检查端口
lsof -i :8088

# 杀死占用进程
kill -9 <PID>
```

### 前端无法访问API

1. 检查API地址配置
2. 检查后端是否启动
3. 检查CORS配置

### 移动端白屏

1. 确保构建时设置了 `CAPACITOR_BUILD=true`
2. 检查API地址是否正确
3. 查看WebView控制台日志

### 数据库连接失败

1. 检查PostgreSQL是否启动
2. 检查用户名密码
3. 检查数据库是否存在

```bash
# 测试连接
psql -U materials -d materials -c "SELECT 1"
```

## 开发工具推荐

### IDE

- **Go**: GoLand / VS Code + Go扩展
- **Vue**: VS Code + Volar扩展
- **Android**: Android Studio

### VS Code扩展

- Go
- Vue - Official
- ESLint
- Prettier
- GitLens

### 其他工具

- Postman: API测试
- DBeaver: 数据库管理
- Sourcetree: Git图形界面
