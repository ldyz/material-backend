# 部署指南

本文档说明了生产环境的部署步骤和配置说明。

## 服务器要求

### 硬件要求

| 配置项 | 最低要求 | 推荐配置 |
|--------|---------|---------|
| CPU | 2核 | 4核+ |
| 内存 | 4GB | 8GB+ |
| 存储 | 50GB | 100GB+ SSD |

### 软件要求

| 软件 | 版本 |
|------|------|
| 操作系统 | Ubuntu 22.04 / CentOS 8+ |
| Go | 1.21+ |
| PostgreSQL | 14+ |
| Nginx | 1.20+ |

## 部署架构

```
┌─────────────────────────────────────────────────────────────┐
│                         用户请求                              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Nginx (反向代理)                           │
│  - HTTPS终止                                                 │
│  - 静态资源服务                                              │
│  - 负载均衡                                                  │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│   Go后端服务     │ │   Go后端服务     │ │   Go后端服务     │
│   (端口8088)    │ │   (端口8088)    │ │   (端口8088)    │
└─────────────────┘ └─────────────────┘ └─────────────────┘
              │               │               │
              └───────────────┼───────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   PostgreSQL 数据库                          │
└─────────────────────────────────────────────────────────────┘
```

## 部署步骤

### 1. 准备服务器

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装依赖
sudo apt install -y git curl wget build-essential

# 安装Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装PostgreSQL
sudo apt install -y postgresql postgresql-contrib

# 安装Nginx
sudo apt install -y nginx
```

### 2. 配置数据库

```bash
# 切换到postgres用户
sudo -u postgres psql

# 创建数据库和用户
CREATE USER materials WITH PASSWORD 'your_secure_password';
CREATE DATABASE materials OWNER materials;
GRANT ALL PRIVILEGES ON DATABASE materials TO materials;
\q

# 配置PostgreSQL
sudo vim /etc/postgresql/14/main/postgresql.conf
# 修改:
# listen_addresses = '*'
# max_connections = 200

sudo vim /etc/postgresql/14/main/pg_hba.conf
# 添加:
# host materials materials 127.0.0.1/32 md5

# 重启PostgreSQL
sudo systemctl restart postgresql
```

### 3. 部署后端

```bash
# 创建应用目录
sudo mkdir -p /opt/material-backend
sudo chown $USER:$USER /opt/material-backend

# 克隆代码
cd /opt/material-backend
git clone <repository-url> .

# 创建配置文件
vim config.yaml
```

**生产环境配置**:

```yaml
server:
  port: 8088
  mode: release

database:
  type: postgresql
  host: 127.0.0.1
  port: 5432
  user: materials
  password: your_secure_password
  database: materials

jwt:
  secret: "your-very-secure-secret-key-at-least-32-chars"
  expire_time: 24h

upload:
  max_file_size: 10485760  # 10MB
  upload_dir: "/opt/material-backend/static/uploads"

ai:
  provider: baidu
  baidu_api_key: "your-api-key"
  baidu_secret_key: "your-secret-key"
```

```bash
# 执行数据库迁移
for f in migrations/*.sql; do
  psql -U materials -d materials -f "$f"
done

# 创建管理员用户
go run cmd/create-admin/main.go --username admin --password your_password

# 编译后端
go build -o server ./cmd/server

# 设置systemd服务
sudo vim /etc/systemd/system/material-backend.service
```

**systemd服务配置**:

```ini
[Unit]
Description=Material Management Backend
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/material-backend
ExecStart=/opt/material-backend/server -c /opt/material-backend/config.yaml
Restart=always
RestartSec=5
StandardOutput=append:/var/log/material-backend.log
StandardError=append:/var/log/material-backend.log

[Install]
WantedBy=multi-user.target
```

```bash
# 创建日志文件
sudo touch /var/log/material-backend.log
sudo chown www-data:www-data /var/log/material-backend.log

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable material-backend
sudo systemctl start material-backend

# 检查状态
sudo systemctl status material-backend
```

### 4. 部署前端

```bash
# Web前端
cd /opt/material-backend/newstatic
npm install
npm run build

# 移动端
cd /opt/material-backend/mobile-app
CAPACITOR_BUILD=true npm run build
```

### 5. 配置Nginx

```bash
sudo vim /etc/nginx/sites-available/material
```

**Nginx配置**:

```nginx
# HTTP重定向到HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

# HTTPS配置
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL证书
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # 静态资源
    location /static/ {
        alias /opt/material-backend/newstatic/dist/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # 移动端
    location /mobile/ {
        alias /opt/material-backend/mobile-app/dist-capacitor/;
        try_files $uri $uri/ /mobile/index.html;
    }

    # 移动端更新包
    location /mobile-updates/ {
        alias /opt/material-backend/mobile-app-updates/;
    }

    # 上传文件
    location /uploads/ {
        alias /opt/material-backend/static/uploads/;
        expires 30d;
    }

    # API代理
    location /api/ {
        proxy_pass http://127.0.0.1:8088;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 86400;
    }

    # WebSocket
    location /api/notification/ws {
        proxy_pass http://127.0.0.1:8088;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_read_timeout 86400;
    }

    # 主页面
    location / {
        root /opt/material-backend/newstatic/dist;
        try_files $uri $uri/ /index.html;
    }

    # 日志
    access_log /var/log/nginx/material-access.log;
    error_log /var/log/nginx/material-error.log;
}
```

```bash
# 启用站点
sudo ln -s /etc/nginx/sites-available/material /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载Nginx
sudo systemctl reload nginx
```

### 6. 配置SSL证书

```bash
# 安装Certbot
sudo apt install -y certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo systemctl enable certbot.timer
```

## 配置说明

### 后端配置

| 配置项 | 说明 | 生产环境建议 |
|--------|------|-------------|
| server.mode | 运行模式 | release |
| jwt.secret | JWT密钥 | 至少32字符随机字符串 |
| jwt.expire_time | Token过期时间 | 24h |
| upload.max_file_size | 最大上传文件 | 根据需求设置 |

### Nginx配置

| 配置项 | 说明 |
|--------|------|
| worker_processes | 工作进程数，建议设置为CPU核心数 |
| worker_connections | 每个进程的最大连接数 |
| keepalive_timeout | 长连接超时时间 |
| client_max_body_size | 最大请求体大小 |

### PostgreSQL配置

| 配置项 | 说明 | 建议值 |
|--------|------|--------|
| max_connections | 最大连接数 | 200 |
| shared_buffers | 共享缓冲区 | 内存的25% |
| effective_cache_size | 有效缓存大小 | 内存的50% |
| work_mem | 工作内存 | 64MB |

## 安全加固

### 1. 防火墙配置

```bash
# 启用UFW
sudo ufw enable

# 只开放必要端口
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS

# 检查状态
sudo ufw status
```

### 2. SSH安全

```bash
sudo vim /etc/ssh/sshd_config

# 修改:
# Port 22 -> Port 自定义端口
# PermitRootLogin no
# PasswordAuthentication no

sudo systemctl restart sshd
```

### 3. 定期更新

```bash
# 设置自动安全更新
sudo apt install -y unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades
```

## 备份策略

### 数据库备份

```bash
# 创建备份脚本
vim /opt/scripts/backup-db.sh
```

```bash
#!/bin/bash
BACKUP_DIR="/opt/backups/db"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

pg_dump -U materials -d materials > $BACKUP_DIR/materials_$DATE.sql

# 保留最近7天的备份
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete

echo "Backup completed: materials_$DATE.sql"
```

```bash
# 设置定时任务
crontab -e

# 每天凌晨2点执行备份
0 2 * * * /opt/scripts/backup-db.sh >> /var/log/backup.log 2>&1
```

### 文件备份

```bash
# 备份上传文件
rsync -avz /opt/material-backend/static/uploads/ /opt/backups/uploads/
```

## 监控告警

### 系统监控

```bash
# 安装监控工具
sudo apt install -y htop iotop nethogs

# 查看系统状态
htop
iostat -x 1
```

### 日志监控

```bash
# 查看后端日志
tail -f /var/log/material-backend.log

# 查看Nginx日志
tail -f /var/log/nginx/material-access.log
tail -f /var/log/nginx/material-error.log

# 查看系统日志
tail -f /var/log/syslog
```

### 服务健康检查

```bash
# 检查后端服务
curl -s http://localhost:8088/api/app/version

# 检查数据库连接
psql -U materials -d materials -c "SELECT 1"
```

## 常见问题

### 1. 服务无法启动

检查日志：
```bash
journalctl -u material-backend -n 50
```

常见原因：
- 配置文件错误
- 端口被占用
- 数据库连接失败

### 2. 数据库连接失败

```bash
# 检查PostgreSQL状态
sudo systemctl status postgresql

# 测试连接
psql -U materials -d materials -h 127.0.0.1
```

### 3. Nginx 502错误

检查后端服务是否运行：
```bash
sudo systemctl status material-backend
```

检查Nginx配置：
```bash
sudo nginx -t
```

### 4. SSL证书问题

```bash
# 续期证书
sudo certbot renew

# 检查证书状态
sudo certbot certificates
```

## 性能优化

### 后端优化

1. 启用GZIP压缩
2. 数据库连接池
3. Redis缓存（可选）
4. CDN加速静态资源

### 数据库优化

1. 定期VACUUM
2. 创建合适的索引
3. 分区大表
4. 查询优化

### Nginx优化

1. 启用HTTP/2
2. 开启GZIP
3. 静态资源缓存
4. 连接池优化

## 升级流程

### 1. 备份数据

```bash
/opt/scripts/backup-db.sh
```

### 2. 拉取最新代码

```bash
cd /opt/material-backend
git pull origin main
```

### 3. 执行数据库迁移

```bash
psql -U materials -d materials -f migrations/new_migration.sql
```

### 4. 重新编译部署

```bash
go build -o server ./cmd/server
sudo systemctl restart material-backend
```

### 5. 更新前端

```bash
cd newstatic
npm install
npm run build
```

### 6. 验证部署

检查各功能是否正常工作。
