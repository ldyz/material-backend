package system

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// pgConfig holds PostgreSQL connection configuration
type pgConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// parsePostgreSQLDSN parses the POSTGRES_DSN environment variable
// Supports both URL format (postgres://user:pass@host:port/db) and key-value format
func parsePostgreSQLDSN(dsn string) (*pgConfig, error) {
	config := &pgConfig{}

	// Try URL format first: postgres://user:pass@host:port/db
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		// Remove protocol prefix
		dsn = strings.TrimPrefix(dsn, "postgres://")
		dsn = strings.TrimPrefix(dsn, "postgresql://")

		// Split authentication@host/database
		parts := strings.Split(dsn, "@")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid DSN format")
		}

		// Parse user:password
		authParts := strings.SplitN(parts[0], ":", 2)
		config.User = authParts[0]
		if len(authParts) > 1 {
			config.Password = authParts[1]
		}

		// Parse host:port/database
		rest := parts[1]
		// Split by / to separate host from database
		slashIdx := strings.Index(rest, "/")
		var hostPort, database string
		if slashIdx >= 0 {
			hostPort = rest[:slashIdx]
			database = rest[slashIdx+1:]
		} else {
			hostPort = rest
		}

		// Remove query parameters if present
		if idx := strings.Index(database, "?"); idx >= 0 {
			database = database[:idx]
		}
		config.Database = database

		// Parse host:port
		if strings.Contains(hostPort, ":") {
			host, port, err := net.SplitHostPort(hostPort)
			if err != nil {
				return nil, fmt.Errorf("invalid host:port: %v", err)
			}
			config.Host = host
			config.Port = port
		} else {
			config.Host = hostPort
			config.Port = "5432"
		}

		return config, nil
	}

	// Try key-value format: key=value key=value...
	pairs := strings.Fields(dsn)
	config.Port = "5432" // default port
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.ToLower(kv[0])
		value := kv[1]

		switch key {
		case "host":
			config.Host = value
		case "port":
			config.Port = value
		case "user":
			config.User = value
		case "password":
			config.Password = value
		case "dbname":
			config.Database = value
		}
	}

	if config.Host == "" || config.User == "" || config.Database == "" {
		return nil, fmt.Errorf("incomplete DSN: missing required fields")
	}

	return config, nil
}

// createSQLBackup creates a SQL backup using the database/sql driver
// This avoids pg_dump version compatibility issues
func createSQLBackup(dsn string, backupPath string) error {
	config, err := parsePostgreSQLDSN(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Build connection string for sql.DB
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create backup file
	file, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	// Write header
	header := fmt.Sprintf("-- PostgreSQL Database Backup\n")
	header += fmt.Sprintf("-- Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	header += fmt.Sprintf("-- Database: %s\n\n", config.Database)
	header += fmt.Sprintf("-- Disabling triggers and foreign keys for faster import\n")
	header += fmt.Sprintf("SET session_replication_role = 'replica';\n\n")
	if _, err := file.WriteString(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Get all table names
	tables, err := getTables(db)
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	// Backup each table
	for _, table := range tables {
		// Skip system tables
		if strings.HasPrefix(table, "pg_") || strings.HasPrefix(table, "sql_") {
			continue
		}

		if err := backupTable(db, file, table); err != nil {
			return fmt.Errorf("failed to backup table %s: %w", table, err)
		}
	}

	// Write footer
	footer := "\n-- Re-enabling triggers and foreign keys\n"
	footer += "SET session_replication_role = 'origin';\n"
	if _, err := file.WriteString(footer); err != nil {
		return fmt.Errorf("failed to write footer: %w", err)
	}

	return nil
}

// getTables returns all table names in the current database
func getTables(db *sql.DB) ([]string, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, rows.Err()
}

// backupTable backs up a single table's schema and data
func backupTable(db *sql.DB, file *os.File, tableName string) error {
	// Quote table name to handle reserved keywords
	quotedTableName := fmt.Sprintf("\"%s\"", tableName)

	// Write comment
	if _, err := file.WriteString(fmt.Sprintf("\n-- Table: %s\n\n", tableName)); err != nil {
		return err
	}

	// Get table schema - skip schema generation for simplicity and reliability
	// Schema will need to be recreated separately or through migrations
	if _, err := file.WriteString(fmt.Sprintf("-- Schema and data for table: %s\n", tableName)); err != nil {
		return err
	}

	// Get column names for INSERT statement
	columnsQuery := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position
	`

	columnRows, err := db.Query(columnsQuery, tableName)
	if err != nil {
		return err
	}

	var columns []string
	for columnRows.Next() {
		var col string
		if err := columnRows.Scan(&col); err != nil {
			columnRows.Close()
			return err
		}
		columns = append(columns, col)
	}
	columnRows.Close()

	if len(columns) == 0 {
		return nil // Skip tables with no columns
	}

	// Quote column names to handle reserved keywords
	quotedColumns := make([]string, len(columns))
	for i, col := range columns {
		quotedColumns[i] = fmt.Sprintf("\"%s\"", col)
	}

	dataQuery := fmt.Sprintf("SELECT %s FROM %s", strings.Join(quotedColumns, ", "), quotedTableName)

	rows, err := db.Query(dataQuery)
	if err != nil {
		return err
	}

	// Prepare INSERT statement header
	insertHeader := fmt.Sprintf("INSERT INTO %s (%s) VALUES\n", quotedTableName, strings.Join(quotedColumns, ", "))

	var valueStrings []string
	var firstRow = true

	for rows.Next() {
		// Create slice of interfaces for scanning
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			rows.Close()
			return err
		}

		// Convert values to SQL literals
		valueStrs := make([]string, len(values))
		for i, v := range values {
			valueStrs[i] = formatSQLValue(v)
		}

		valueStrings = append(valueStrings, "("+strings.Join(valueStrs, ", ")+")")

		// Write in batches of 100 to avoid huge statements
		if len(valueStrings) >= 100 {
			if !firstRow {
				if _, err := file.WriteString(",\n"); err != nil {
					rows.Close()
					return err
				}
			}
			if _, err := file.WriteString(insertHeader + strings.Join(valueStrings, ",\n") + ";\n"); err != nil {
				rows.Close()
				return err
			}
			valueStrings = nil
			firstRow = true
		} else {
			firstRow = false
		}
	}
	rows.Close()

	// Write remaining values
	if len(valueStrings) > 0 {
		if _, err := file.WriteString(insertHeader + strings.Join(valueStrings, ",\n") + ";\n"); err != nil {
			return err
		}
	}

	return nil
}

// formatSQLValue formats a value for SQL INSERT statement
func formatSQLValue(v interface{}) string {
	if v == nil {
		return "NULL"
	}

	switch val := v.(type) {
	case []byte:
		str := string(val)
		return "'" + strings.ReplaceAll(escapeSQLString(str), "'", "''") + "'"
	case string:
		return "'" + strings.ReplaceAll(escapeSQLString(val), "'", "''") + "'"
	case time.Time:
		return "'" + val.Format("2006-01-02 15:04:05") + "'"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// escapeSQLString escapes special characters in SQL strings
func escapeSQLString(s string) string {
	replacer := strings.NewReplacer(
		"\x00", "\\x00",
		"\n", "\\n",
		"\r", "\\r",
		"\\", "\\\\",
	)
	return replacer.Replace(s)
}

// restoreDatabase restores the database from a backup SQL file
// This is a dangerous operation that will delete all existing data
func restoreDatabase(config *pgConfig, backupContent string) error {
	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Get all table names
	tables, err := getTables(db)
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	// Disable foreign key constraints temporarily
	if _, err := db.Exec("SET session_replication_role = 'replica'"); err != nil {
		return fmt.Errorf("failed to disable triggers: %w", err)
	}

	// Delete all data from each table
	for _, table := range tables {
		// Skip system tables
		if strings.HasPrefix(table, "pg_") || strings.HasPrefix(table, "sql_") {
			continue
		}

		quotedTableName := fmt.Sprintf("\"%s\"", table)
		if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", quotedTableName)); err != nil {
			return fmt.Errorf("failed to delete from table %s: %w", table, err)
		}

		// Reset sequences
		if _, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", quotedTableName)); err != nil {
			// If TRUNCATE fails (e.g., due to foreign keys), try DELETE instead
			if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", quotedTableName)); err != nil {
				return fmt.Errorf("failed to truncate table %s: %w", table, err)
			}
		}
	}

	// Split the backup SQL into individual statements
	// Simple approach: split by semicolon, but handle multi-line CREATE statements
	statements := parseSQLStatements(backupContent)

	// Execute each statement
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		if _, err := db.Exec(stmt); err != nil {
			// Log the error but continue with other statements
			fmt.Printf("Warning: failed to execute statement: %v\nStatement: %s\n", err, stmt)
		}
	}

	// Re-enable foreign key constraints
	if _, err := db.Exec("SET session_replication_role = 'origin'"); err != nil {
		return fmt.Errorf("failed to re-enable triggers: %w", err)
	}

	return nil
}

// parseSQLStatements parses SQL content into individual statements
// Handles multi-line INSERT statements and comments
func parseSQLStatements(content string) []string {
	var statements []string
	var currentStmt strings.Builder
	lines := strings.Split(content, "\n")

	inInsert := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip comments
		if strings.HasPrefix(trimmed, "--") {
			// If we're building a statement, finish it first
			if currentStmt.Len() > 0 {
				statements = append(statements, currentStmt.String())
				currentStmt.Reset()
				inInsert = false
			}
			continue
		}

		// Detect INSERT statement start
		if strings.HasPrefix(strings.ToUpper(trimmed), "INSERT INTO") {
			inInsert = true
		}

		// Add line to current statement
		if currentStmt.Len() > 0 {
			currentStmt.WriteString(" ")
		}
		currentStmt.WriteString(trimmed)

		// Check if statement ends with semicolon
		if strings.HasSuffix(trimmed, ";") {
			stmt := currentStmt.String()

			// For INSERT statements, include the entire multi-line statement
			if inInsert || strings.HasPrefix(strings.ToUpper(stmt), "INSERT INTO") ||
			   strings.HasPrefix(strings.ToUpper(stmt), "CREATE") ||
			   strings.HasPrefix(strings.ToUpper(stmt), "ALTER") ||
			   strings.HasPrefix(strings.ToUpper(stmt), "SET") {
				statements = append(statements, stmt)
				currentStmt.Reset()
				inInsert = false
			} else {
				// Simple statement, add it
				statements = append(statements, stmt)
				currentStmt.Reset()
			}
		}
	}

	// Add any remaining statement
	if currentStmt.Len() > 0 {
		statements = append(statements, currentStmt.String())
	}

	return statements
}

// RegisterRoutes 注册系统模块路由
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// 公开接口 - 不需要认证
	// 获取基本系统设置（包括系统名称等公开信息）
	rg.GET("/system/public/settings", func(c *gin.Context) {
		var configs []SystemConfig
		if err := db.Where("\"key\" IN (?)", []string{"system_name", "system_version"}).Find(&configs).Error; err != nil {
			response.InternalError(c, fmt.Sprintf("获取系统设置失败: %v", err))
			return
		}

		settings := map[string]any{}
		for _, config := range configs {
			settings[config.Key] = config.Value
		}

		response.Success(c, settings)
	})

	// 创建路由组
	r := rg.Group("/system")
	// 使用JWT中间件进行身份验证
	r.Use(jwtpkg.TokenMiddleware())

	// 注册AI相关路由
	RegisterAIRoutes(r, db)

	// 注册报表相关路由
	RegisterReportRoutes(r, db)

	// ================== 系统日志接口 ==================
	// 仪表板使用的日志接口，所有登录用户都可以访问
	r.GET("/logs", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		level := c.Query("level")
		date := c.Query("date")
		startTime := c.Query("start_time")
		endTime := c.Query("end_time")

		query := db.Model(&SystemLog{})

		// 筛选条件
		if level != "" {
			query = query.Where("level = ?", level)
		}
		if date != "" {
			// 仅日期，筛选当天
			if t, err := time.Parse("2006-01-02", date); err == nil {
				dtEnd := t.Add(24*time.Hour - time.Second)
				query = query.Where("created_at >= ? AND created_at <= ?", t, dtEnd)
			}
		}
		if startTime != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
				query = query.Where("created_at >= ?", t)
			}
		}
		if endTime != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
				query = query.Where("created_at <= ?", t)
			}
		}

		var total int64
		query.Count(&total)

		var logs []SystemLog
		query.Order("created_at DESC").
			Offset((page - 1) * pageSize).
			Limit(pageSize).
			Find(&logs)

		logList := make([]map[string]any, 0, len(logs))
		for _, log := range logs {
			logEntry := map[string]any{
				"id":         log.ID,
				"level":      log.Level,
				"message":    log.Message,
				"module":     log.Module,
				"user_id":    log.UserID,
				"ip_address": log.IPAddress,
				"created_at": log.CreatedAt.Format("2006-01-02 15:04:05"),
			}

			// 查询用户名
			if log.UserID > 0 {
				var user struct {
					ID       uint
					Username string
				}
				if err := db.Model(&struct{}{}).Table("users").Where("id = ?", log.UserID).
					Select("id, username").Scan(&user).Error; err == nil && user.ID > 0 {
					logEntry["user"] = user.Username
				} else {
					logEntry["user"] = "未知用户"
				}
			} else {
				logEntry["user"] = "系统"
			}

			logList = append(logList, logEntry)
		}

		response.SuccessWithPagination(c, logList, int64(page), int64(pageSize), total)
	})

	r.POST("/logs/clear", auth.PermissionMiddleware(db, "system_log"), func(c *gin.Context) {
		if err := db.Where("1 = 1").Delete(&SystemLog{}).Error; err != nil {
			response.InternalError(c, fmt.Sprintf("清空日志失败: %v", err))
			return
		}

		// 记录日志
		currentUser, _ := auth.GetCurrentUser(c, db)
		var userID uint
		if currentUser != nil {
			userID = currentUser.ID
		}
		db.Create(&SystemLog{
			Level:     "WARNING",
			Message:   "清空系统日志",
			Module:    "system",
			UserID:    userID,
			IPAddress: c.ClientIP(),
		})

		response.SuccessWithMessage(c, nil, "系统日志已清空")
	})

	r.DELETE("/logs", auth.PermissionMiddleware(db, "system_log"), func(c *gin.Context) {
		if err := db.Where("1 = 1").Delete(&SystemLog{}).Error; err != nil {
			response.InternalError(c, fmt.Sprintf("清空日志失败: %v", err))
			return
		}

		// 记录日志
		currentUser, _ := auth.GetCurrentUser(c, db)
		var userID uint
		if currentUser != nil {
			userID = currentUser.ID
		}
		db.Create(&SystemLog{
			Level:     "WARNING",
			Message:   "清空系统日志",
			Module:    "system",
			UserID:    userID,
			IPAddress: c.ClientIP(),
		})

		response.SuccessWithMessage(c, nil, "系统日志已清空")
	})

	// ================== 系统统计接口 ==================
	// 仪表板使用的统计接口，所有登录用户都可以访问
	r.GET("/stats", func(c *gin.Context) {
		// 初始化统计数据
		stats := map[string]any{
			"total_materials":      0,
			"total_stock_value":    0,
			"pending_requisitions": 0,
			"low_stock_count":      0,
			"total_users":          0,
			"total_projects":       0,
			"monthly_requisitions": 0,
			"monthly_inbound":      0,
		}

		// 获取用户总数
		var totalUsers int64
		db.Table("users").Count(&totalUsers)
		stats["total_users"] = totalUsers

		// 获取物资总数
		var totalMaterials int64
		db.Table("material_master").Count(&totalMaterials)
		stats["total_materials"] = totalMaterials

		// 获取库存相关统计
		var totalStockItems, lowStockCount int64
		db.Table("stocks").Count(&totalStockItems)
		db.Table("stocks").Where("quantity <= 10").Count(&lowStockCount)
		stats["total_stock_value"] = totalStockItems
		stats["low_stock_count"] = lowStockCount

		// 获取领料单统计
		var pendingRequisitions, monthlyRequisitions int64
		db.Table("requisitions").Where("status = ?", "pending").Count(&pendingRequisitions)
		// 获取本月领料单数量
		currentMonth := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -int(time.Now().Day())+1)
		db.Table("requisitions").Where("created_at >= ?", currentMonth).Count(&monthlyRequisitions)
		stats["pending_requisitions"] = pendingRequisitions
		stats["monthly_requisitions"] = monthlyRequisitions

		// 获取项目总数
		var totalProjects int64
		db.Table("projects").Count(&totalProjects)
		stats["total_projects"] = totalProjects

		// 获取入库单统计
		var monthlyInbound int64
		db.Table("inbound_orders").Where("created_at >= ?", currentMonth).Count(&monthlyInbound)
		stats["monthly_inbound"] = monthlyInbound

		response.Success(c, stats)
	})

	// ================== 系统信息接口 ==================
	r.GET("/info", auth.PermissionMiddleware(db, "system_statistics"), func(c *gin.Context) {
		// 获取系统基本信息
		systemInfo := map[string]any{
			"system":           "Linux", // 简化实现，假设在Linux上运行
			"platform":         "Linux",
			"python_version":   "3.10+", // Go后端不需要Python版本
			"current_time":     time.Now().Format("2006-01-02 15:04:05"),
			"app_version":      "1.0.0",
			"cpu_count":        "N/A",
			"memory_total":     "N/A",
			"memory_available": "N/A",
			"disk_total":       "N/A",
			"disk_free":        "N/A",
		}

		// 尝试获取磁盘信息
		if stat, err := os.Stat("."); err == nil {
			if stat.Sys() != nil {
				// 简化实现，不使用psutil包
			}
		}

		response.Success(c, systemInfo)
	})

	// ================== 备份管理接口 ==================

	// 创建备份
	r.POST("/backup", auth.PermissionMiddleware(db, "system_backup"), func(c *gin.Context) {
		// 生成备份文件名
		timestamp := time.Now().Format("20060102_150405")
		backupFilename := fmt.Sprintf("backup_%s.sql", timestamp)
		backupPath := filepath.Join(".", backupFilename)

		// 获取PostgreSQL连接信息
		pgDSN := os.Getenv("POSTGRES_DSN")
		if pgDSN == "" {
			response.InternalError(c, "PostgreSQL连接信息未配置")
			return
		}

		// 使用Go原生备份功能（不依赖pg_dump）
		if err := createSQLBackup(pgDSN, backupPath); err != nil {
			response.InternalError(c, fmt.Sprintf("备份失败: %v", err))
			return
		}

		// 获取备份文件大小
		stat, err := os.Stat(backupPath)
		if err != nil {
			response.InternalError(c, fmt.Sprintf("获取备份文件信息失败: %v", err))
			return
		}

		// 记录备份信息到数据库
		currentUser, _ := auth.GetCurrentUser(c, db)
		createdBy := "system"
		var userID uint = 0
		if currentUser != nil {
			createdBy = currentUser.Username
			userID = currentUser.ID
		}

		backupRecord := SystemBackup{
			Filename:    backupFilename,
			Filepath:    backupPath,
			Size:        stat.Size(),
			CreatedBy:   createdBy,
			Description: fmt.Sprintf("PostgreSQL数据库备份 - %s", timestamp),
		}
		// Explicitly set ID to 0 to ensure GORM treats it as unset for auto-increment
		backupRecord.ID = 0

		if err := db.Omit("ID").Create(&backupRecord).Error; err != nil {
			response.InternalError(c, fmt.Sprintf("记录备份信息失败: %v", err))
			return
		}

		// 记录系统日志
		logRecord := SystemLog{
			Level:     "INFO",
			Message:   fmt.Sprintf("执行数据库备份: %s", backupFilename),
			Module:    "backup",
			UserID:    userID,
			IPAddress: c.ClientIP(),
		}
		db.Create(&logRecord)

		response.Created(c, map[string]any{
			"filename":    backupFilename,
			"backup_file": backupFilename,
			"size":        stat.Size(),
			"created_at":  time.Now().Format("2006-01-02 15:04:05"),
		}, "数据库备份成功")
	})

	// 获取备份历史
	r.GET("/backup/history", auth.PermissionMiddleware(db, "system_backup"), func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

		// 限制每页数量
		if pageSize > 100 {
			pageSize = 100
		}

		var total int64
		if err := db.Model(&SystemBackup{}).Count(&total).Error; err != nil {
			response.InternalError(c, fmt.Sprintf("获取备份总数失败: %v", err))
			return
		}

		var backups []SystemBackup
		if err := db.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&backups).Error; err != nil {
			response.InternalError(c, fmt.Sprintf("获取备份历史失败: %v", err))
			return
		}

		// 转换为DTO
		backupList := make([]map[string]any, len(backups))
		for i, backup := range backups {
			backupList[i] = backup.ToDTO()
		}

		response.SuccessWithPagination(c, backupList, int64(page), int64(pageSize), total)
	})

	// 获取备份列表
	r.GET("/backup", auth.PermissionMiddleware(db, "system_backup"), func(c *gin.Context) {
		// 获取数据库中的备份记录
		var backups []SystemBackup
		if err := db.Order("created_at desc").Find(&backups).Error; err != nil {
			response.InternalError(c, fmt.Sprintf("获取备份记录失败: %v", err))
			return
		}

		// 同时扫描文件系统中的备份文件
		backupFiles := []map[string]any{}
		files, err := filepath.Glob("backup_*.sql")
		if err == nil {
			for _, file := range files {
				if stat, err := os.Stat(file); err == nil {
					backupFiles = append(backupFiles, map[string]any{
						"name":       file,
						"filename":   file,
						"size":       stat.Size(),
						"date":       stat.ModTime().Format("2006-01-02 15:04:05"),
						"created_at": stat.ModTime().Format("2006-01-02 15:04:05"),
					})
				}
			}
		}

		// 合并数据库记录和文件系统扫描结果
		backupList := []map[string]any{}
		dbFilenames := map[string]bool{}

		// 添加数据库记录
		for _, backup := range backups {
			backupList = append(backupList, backup.ToDTO())
			dbFilenames[backup.Filename] = true
		}

		// 添加文件系统中存在但数据库中没有记录的文件
		for _, fileInfo := range backupFiles {
			if filename, ok := fileInfo["filename"].(string); ok {
				if !dbFilenames[filename] {
					backupList = append(backupList, fileInfo)
				}
			}
		}

		// 按时间排序
		sort.Slice(backupList, func(i, j int) bool {
			ti, _ := time.Parse("2006-01-02 15:04:05", backupList[i]["created_at"].(string))
			tj, _ := time.Parse("2006-01-02 15:04:05", backupList[j]["created_at"].(string))
			return ti.After(tj)
		})

		response.Success(c, backupList)
	})

	// 下载备份文件 (RESTful方式)
	r.GET("/backup/:backup_name/download", auth.PermissionMiddleware(db, "system_backup"), func(c *gin.Context) {
		backupName := c.Param("backup_name")
		if backupName == "" || backupName == "undefined" {
			response.BadRequest(c, "备份文件名不能为空")
			return
		}

		backupPath := filepath.Join(".", backupName)
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			response.NotFound(c, "备份文件不存在")
			return
		}

		c.FileAttachment(backupPath, backupName)
	})

	// 删除备份文件 (RESTful方式)
	r.DELETE("/backup/:backup_name", auth.PermissionMiddleware(db, "system_backup"), func(c *gin.Context) {
		backupName := c.Param("backup_name")
		if backupName == "" || backupName == "undefined" {
			response.BadRequest(c, "备份文件名不能为空")
			return
		}

		backupPath := filepath.Join(".", backupName)

		// 检查文件是否存在
		fileExists := true
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			fileExists = false
		}

		// 如果文件存在，删除文件
		if fileExists {
			if err := os.Remove(backupPath); err != nil {
				response.InternalError(c, fmt.Sprintf("删除文件失败: %v", err))
				return
			}
		}

		// 从数据库中删除记录
		db.Where("filename = ?", backupName).Delete(&SystemBackup{})

		message := "备份文件删除成功"
		if !fileExists {
			message = "备份记录已删除（文件不存在）"
		}

		response.SuccessWithMessage(c, nil, message)
	})

	// 传统方式：删除备份文件 (使用POST body)
	r.POST("/backup/delete", auth.PermissionMiddleware(db, "system_backup"), func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数错误")
			return
		}

		if req.Name == "" {
			response.BadRequest(c, "备份文件名不能为空")
			return
		}

		backupPath := filepath.Join(".", req.Name)

		// 检查文件是否存在
		fileExists := true
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			fileExists = false
		}

		// 如果文件存在，删除文件
		if fileExists {
			if err := os.Remove(backupPath); err != nil {
				response.InternalError(c, fmt.Sprintf("删除文件失败: %v", err))
				return
			}
		}

		// 从数据库中删除记录
		db.Where("filename = ?", req.Name).Delete(&SystemBackup{})

		message := "备份文件删除成功"
		if !fileExists {
			message = "备份记录已删除（文件不存在）"
		}

		response.SuccessWithMessage(c, nil, message)
	})

	// 恢复备份 - 危险操作，会清空当前数据库
	r.POST("/backup/restore", auth.PermissionMiddleware(db, "system_backup"), func(c *gin.Context) {
		var req struct {
			BackupName string `json:"backup_name" binding:"required"`
			Confirm    bool   `json:"confirm" binding:"required"` // 必须明确确认
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求数据格式错误")
			return
		}

		// 必须明确确认
		if !req.Confirm {
			response.BadRequest(c, "必须明确确认恢复操作，此操作将清空当前数据库")
			return
		}

		// 验证备份文件存在
		backupPath := filepath.Join(".", req.BackupName)
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			response.NotFound(c, "备份文件不存在")
			return
		}

		// 读取备份文件内容
		backupContent, err := os.ReadFile(backupPath)
		if err != nil {
			response.InternalError(c, fmt.Sprintf("读取备份文件失败: %v", err))
			return
		}

		// 获取数据库连接信息
		pgDSN := os.Getenv("POSTGRES_DSN")
		if pgDSN == "" {
			response.InternalError(c, "PostgreSQL连接信息未配置")
			return
		}

		config, err := parsePostgreSQLDSN(pgDSN)
		if err != nil {
			response.InternalError(c, fmt.Sprintf("解析数据库配置失败: %v", err))
			return
		}

		// 获取当前用户信息用于日志
		currentUser, _ := auth.GetCurrentUser(c, db)
		var userID uint = 0
		if currentUser != nil {
			userID = currentUser.ID
		}

		// 执行恢复操作
		if err := restoreDatabase(config, string(backupContent)); err != nil {
			// 记录失败日志
			db.Create(&SystemLog{
				Level:     "ERROR",
				Message:   fmt.Sprintf("数据库恢复失败: %s - %v", req.BackupName, err),
				Module:    "backup",
				UserID:    userID,
				IPAddress: c.ClientIP(),
			})

			response.InternalError(c, fmt.Sprintf("数据库恢复失败: %v", err))
			return
		}

		// 记录成功日志
		db.Create(&SystemLog{
			Level:     "WARNING",
			Message:   fmt.Sprintf("数据库恢复成功: %s", req.BackupName),
			Module:    "backup",
			UserID:    userID,
			IPAddress: c.ClientIP(),
		})

		response.SuccessWithMessage(c, nil, "数据库恢复成功，请刷新页面")
	})

	// ================== 活动记录接口 ==================
	r.GET("/recent-activities", auth.PermissionMiddleware(db, "system_activities"), func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if limit <= 0 {
			limit = 10
		}

		activities := []map[string]any{}

		// 获取最近的库存操作记录
		type StockLogWithStock struct {
			StockLogID   uint      `json:"stock_log_id"`
			StockID      uint      `json:"stock_id"`
			Type         string    `json:"type"`
			Quantity     float64   `json:"quantity"`
			Unit         string    `json:"unit"`
			Recipient    string    `json:"recipient"`
			Remark       string    `json:"remark"`
			Time         time.Time `json:"time"`
			StockName    string    `json:"stock_name"`
		}

		var stockLogs []StockLogWithStock
		db.Raw(`
			SELECT sl.id as stock_log_id, sl.stock_id, sl.type, sl.quantity, s.unit, sl.recipient, sl.remark, sl.time, m.name as stock_name
			FROM stock_logs sl
			JOIN stocks s ON sl.stock_id = s.id
			JOIN materials m ON s.material_id = m.id
			ORDER BY sl.time DESC
			LIMIT ?
		`, limit).Scan(&stockLogs)

		for _, log := range stockLogs {
			activityType := "stock_in"
			if log.Type != "in" {
				activityType = "stock_out"
			}

			title := fmt.Sprintf("%s: %s", map[string]string{"in": "入库", "out": "出库"}[log.Type], log.StockName)
			description := fmt.Sprintf("数量: %g %s", log.Quantity, log.Unit)
			if log.Remark != "" {
				description += fmt.Sprintf(", 备注: %s", log.Remark)
			}

			activities = append(activities, map[string]any{
				"id":         fmt.Sprintf("stock_%d", log.StockLogID),
				"title":      title,
				"description": description,
				"type":       activityType,
				"created_at": log.Time.Format("2006-01-02T15:04:05Z"),
				"user":       log.Recipient,
			})
		}

		// 按时间排序，取最新的记录
		sort.Slice(activities, func(i, j int) bool {
			return activities[i]["created_at"].(string) > activities[j]["created_at"].(string)
		})

		// 限制数量
		if len(activities) > limit {
			activities = activities[:limit]
		}

		response.SuccessWithMeta(c, activities, map[string]any{"total": len(activities)})
	})

	// ================== 物资分类统计接口 ==================
	r.GET("/material-category-stats", auth.PermissionMiddleware(db, "system_statistics"), func(c *gin.Context) {
		// 按物资分类统计数量
		type CategoryStat struct {
			Category string `json:"category"`
			Count    int64  `json:"count"`
		}

		var stats []CategoryStat
		if err := db.Table("material_master").
			Select("COALESCE(category, '未分类') as category, COUNT(*) as count").
			Group("category").
			Order("count DESC").
			Find(&stats).Error; err != nil {
			response.InternalError(c, fmt.Sprintf("获取物资分类统计失败: %v", err))
			return
		}

		response.Success(c, stats)
	})

	// ================== 项目物资统计接口 ==================
	r.GET("/project-material-stats", auth.PermissionMiddleware(db, "system_statistics"), func(c *gin.Context) {
		// 按项目统计物资数量
		type ProjectMaterialStat struct {
			ProjectID   uint   `json:"project_id"`
			ProjectName string `json:"project"`
			Count       int64  `json:"count"`
		}

		var stats []ProjectMaterialStat
		if err := db.Table("materials AS m").
			Select("COALESCE(m.project_id, 0) as project_id, COALESCE(p.name, '未分配') as project, COUNT(*) as count").
			Joins("LEFT JOIN projects AS p ON m.project_id = p.id").
			Group("m.project_id, p.name").
			Order("count DESC").
			Find(&stats).Error; err != nil {
			response.InternalError(c, fmt.Sprintf("获取项目物资统计失败: %v", err))
			return
		}

		response.Success(c, stats)
	})

	// ================== 系统设置接口 ==================
	r.GET("/settings", auth.PermissionMiddleware(db, "system_config"), func(c *gin.Context) {
		var configs []SystemConfig
		if err := db.Find(&configs).Error; err != nil {
			response.InternalError(c, fmt.Sprintf("获取系统设置失败: %v", err))
			return
		}

		settings := map[string]any{}
		for _, config := range configs {
			// 根据类型转换值
			switch config.Type {
			case "boolean":
				value, _ := strconv.ParseBool(config.Value)
				settings[config.Key] = value
			case "integer":
				value, _ := strconv.Atoi(config.Value)
				settings[config.Key] = value
			case "float":
				value, _ := strconv.ParseFloat(config.Value, 64)
				settings[config.Key] = value
			default:
				settings[config.Key] = config.Value
			}
		}

		response.Success(c, settings)
	})

	// 更新系统设置
	r.PUT("/settings", auth.PermissionMiddleware(db, "system_config"), func(c *gin.Context) {
		var updateData map[string]any
		if err := c.ShouldBindJSON(&updateData); err != nil {
			response.BadRequest(c, fmt.Sprintf("无效的请求数据: %v", err))
			return
		}

		// 获取当前所有配置
		var configs []SystemConfig
		if err := db.Find(&configs).Error; err != nil {
			response.InternalError(c, fmt.Sprintf("获取系统设置失败: %v", err))
			return
		}

		// 创建配置键到ID的映射
		configMap := map[string]uint{}
		for _, config := range configs {
			configMap[config.Key] = config.ID
		}

		// 处理更新
		for key, value := range updateData {
			// 转换值为字符串
			var strValue string
			var configType string

			switch v := value.(type) {
			case bool:
				strValue = strconv.FormatBool(v)
				configType = "boolean"
			case float64:
				// 检查是否为整数
				if v == float64(int(v)) {
					strValue = strconv.Itoa(int(v))
					configType = "integer"
				} else {
					strValue = strconv.FormatFloat(v, 'f', -1, 64)
					configType = "float"
				}
			default:
				strValue = fmt.Sprintf("%v", v)
				configType = "string"
			}

			// 检查配置是否存在
			if id, exists := configMap[key]; exists {
				// 更新现有配置
				if err := db.Model(&SystemConfig{}).Where("id = ?", id).Updates(map[string]any{
					"value": strValue,
					"type":  configType,
				}).Error; err != nil {
					response.InternalError(c, fmt.Sprintf("更新配置项 %s 失败: %v", key, err))
					return
				}
			} else {
				// 创建新配置
				newConfig := SystemConfig{
					Key:   key,
					Value: strValue,
					Type:  configType,
				}
				if err := db.Create(&newConfig).Error; err != nil {
					response.InternalError(c, fmt.Sprintf("创建配置项 %s 失败: %v", key, err))
					return
				}
			}
		}

		response.SuccessWithMessage(c, nil, "系统设置更新成功")
	})
}