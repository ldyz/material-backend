# 数据库连接信息

## PostgreSQL
- **Host**: 127.0.0.1
- **Port**: 5432
- **User**: materials
- **Password**: julei1984
- **Database**: materials
- **SSL Mode**: disable

## 连接字符串
```
host=127.0.0.1 port=5432 user=materials password=julei1984 dbname=materials sslmode=disable
```

## psql 命令
```bash
psql -h 127.0.0.1 -p 5432 -U materials -d materials
```

## 在代码中使用 (GORM)
```go
dsn := "host=127.0.0.1 port=5432 user=materials password=julei1984 dbname=materials sslmode=disable"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```
