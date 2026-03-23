# 配置文件说明

## 配置文件位置

项目使用 YAML 格式的配置文件：`config.yaml`

## 启动参数

```bash
# 使用默认配置文件 (config.yaml)
go run cmd/server/main.go

# 指定配置文件
go run cmd/server/main.go -c /path/to/config.yaml

# 编译后运行
./server -c config.yaml
```

## 配置项说明

### 服务器配置 (server)

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| port | int | 8088 | 服务监听端口 |
| mode | string | debug | 运行模式: debug/release/test |
| read_timeout | duration | 60s | 读取超时时间 |
| write_timeout | duration | 60s | 写入超时时间 |
| shutdown_timeout | duration | 10s | 优雅关闭超时时间 |

### 数据库配置 (database)

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| type | string | postgresql | 数据库类型: postgresql/mysql/sqlite |
| host | string | 127.0.0.1 | 数据库主机地址 |
| port | int | 5432 | 数据库端口 |
| user | string | - | 数据库用户名 |
| password | string | - | 数据库密码 |
| database | string | - | 数据库名称 |
| ssl_mode | string | disable | SSL模式 |
| max_idle_conns | int | 10 | 最大空闲连接数 |
| max_open_conns | int | 100 | 最大打开连接数 |
| max_lifetime | int | 3600 | 连接最大生存时间(秒) |

**数据库类型说明：**

- **PostgreSQL**:
  ```yaml
  database:
    type: postgresql
    host: 127.0.0.1
    port: 5432
    user: materials
    password: your-password
    database: materials
  ```

- **MySQL**:
  ```yaml
  database:
    type: mysql
    host: 127.0.0.1
    port: 3306
    user: root
    password: your-password
    database: materials
  ```

- **SQLite**:
  ```yaml
  database:
    type: sqlite
    database: ./data/materials.db
  ```

### JWT配置 (jwt)

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| secret | string | - | JWT签名密钥 (生产环境必须修改) |
| expire_time | duration | 24h | Token过期时间 |
| issuer | string | material-backend | JWT签发者 |

### 上传配置 (upload)

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| max_file_size | int64 | 5242880 | 最大文件大小(字节) |
| max_upload_count | int | 10 | 单次最大上传文件数 |
| allowed_types | string | - | 允许的文件类型(逗号分隔) |
| upload_dir | string | static/uploads | 上传目录 |
| generate_file_name | bool | true | 是否生成随机文件名 |

### 日志配置 (log)

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| level | string | info | 日志级别: debug/info/warn/error |
| file_name | string | logs/app.log | 日志文件路径 |
| max_size | int | 100 | 单个日志文件最大大小(MB) |
| max_backups | int | 3 | 保留的旧日志文件数量 |
| max_age | int | 7 | 保留旧日志文件的最大天数 |
| compress | bool | true | 是否压缩旧日志文件 |

## 环境迁移

从 .env 文件迁移到 config.yaml:

旧方式 (.env):
```bash
POSTGRES_DSN=host=127.0.0.1 port=5432 user=materials password=julei1984 dbname=materials sslmode=disable
```

新方式 (config.yaml):
```yaml
database:
  type: postgresql
  host: 127.0.0.1
  port: 5432
  user: materials
  password: julei1984
  database: materials
  ssl_mode: disable
```

## 配置文件热更新

配置系统支持监听文件变化，自动重载配置（开发中功能）。

## 安全建议

1. **生产环境必须修改默认配置**:
   - JWT secret
   - 数据库密码
   - 服务器模式改为 `release`

2. **敏感信息保护**:
   - 不要将包含真实密码的 config.yaml 提交到版本控制
   - 使用 config.example.yaml 作为模板
   - 生产环境使用环境变量或密钥管理系统

3. **文件权限**:
   ```bash
   chmod 600 config.yaml  # 只有所有者可读写
   ```
