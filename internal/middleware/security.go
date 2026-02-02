package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware 安全参数过滤中间件
// 防止 SQL 注入、XSS、路径遍历等攻击
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过静态文件请求（不需要安全检查）
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static/") ||
		   strings.HasPrefix(path, "/assets/") ||
		   strings.HasPrefix(path, "/mobile/") {
			c.Next()
			return
		}

		// 1. 过滤 Query 参数
		if err := sanitizeParams(c.Request.URL.Query()); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "请求参数包含非法字符",
			})
			return
		}

		// 2. 过滤 PostForm 参数
		if err := sanitizeParams(c.Request.PostForm); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "请求参数包含非法字符",
			})
			return
		}

		// 3. 检查 Path 参数
		for _, param := range c.Params {
			if err := sanitizeString(param.Value); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   fmt.Sprintf("路径参数 %s 包含非法字符", param.Key),
				})
				return
			}
		}

		// 4. 设置安全响应头
		setSecurityHeaders(c)

		c.Next()
	}
}

// sanitizeParams 清理参数
func sanitizeParams(params map[string][]string) error {
	for key, values := range params {
		// 检查参数名是否合法
		if err := validateParamName(key); err != nil {
			return err
		}

		// 清理每个值
		for i, value := range values {
			cleaned, err := sanitizeInput(value)
			if err != nil {
				return err
			}
			values[i] = cleaned
		}
	}
	return nil
}

// validateParamName 验证参数名是否合法
func validateParamName(name string) error {
	// 参数名只能包含字母、数字、下划线、中划线和点
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, name)
	if !matched {
		return fmt.Errorf("非法参数名: %s", name)
	}
	return nil
}

// sanitizeInput 清理输入字符串
func sanitizeInput(input string) (string, error) {
	// 1. 检查输入长度
	if utf8.RuneCountInString(input) > 10000 {
		return "", fmt.Errorf("输入过长")
	}

	// 2. 检测 SQL 注入模式
	if containsSQLInjection(input) {
		return "", fmt.Errorf("检测到 SQL 注入尝试")
	}

	// 3. 检测 XSS 攻击模式
	if containsXSS(input) {
		return "", fmt.Errorf("检测到 XSS 攻击尝试")
	}

	// 4. 检测路径遍历
	if containsPathTraversal(input) {
		return "", fmt.Errorf("检测到路径遍历尝试")
	}

	// 5. 过滤危险字符
	cleaned := filterDangerousChars(input)

	// 6. 再次检查清理后的长度
	if utf8.RuneCountInString(cleaned) > 10000 {
		return "", fmt.Errorf("清理后输入仍过长")
	}

	return cleaned, nil
}

// sanitizeString 清理单个字符串
func sanitizeString(input string) error {
	cleaned, err := sanitizeInput(input)
	if err != nil {
		return err
	}
	if cleaned != input {
		return fmt.Errorf("输入包含非法字符")
	}
	return nil
}

// containsSQLInjection 检测 SQL 注入模式
func containsSQLInjection(input string) bool {
	// 转换为小写进行检测
	lowerInput := strings.ToLower(input)

	// SQL 注入关键字列表
	sqlKeywords := []string{
		"select", "insert", "update", "delete", "drop", "union",
		"exec", "execute", "script", "javascript:", "onerror=", "onload=",
		"eval(", "expression(", "alert(", "document.cookie", "window.location",
		"--", ";--", "/*", "*/", "xp_", "sp_", "0x", "char(",
		"waitfor delay", "benchmark(", "sleep(",
	}

	// 检查 SQL 注入关键字
	for _, keyword := range sqlKeywords {
		if strings.Contains(lowerInput, keyword) {
			// 检查是否是独立的单词（而不是子字符串）
			if isWholeWord(lowerInput, keyword) {
				return true
			}
		}
	}

	// 检测单引号逃逸
	if strings.Contains(input, "\\") || strings.Contains(input, "\\'") {
		return true
	}

	return false
}

// containsXSS 检测 XSS 攻击模式
func containsXSS(input string) bool {
	lowerInput := strings.ToLower(input)

	// XSS 攻击模式
	xssPatterns := []string{
		"<script", "</script>", "javascript:", "vbscript:",
		"onload=", "onerror=", "onclick=", "onmouseover=",
		"onfocus=", "onblur=", "onkeydown=", "onkeyup=",
		"eval(", "expression(", "alert(", "confirm(",
		"document.cookie", "window.location", "document.location.href",
		"fromcharcode", ".innerHTML", "outerHTML",
		"iframe", "<object", "<embed",
	}

	for _, pattern := range xssPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}

	return false
}

// containsPathTraversal 检测路径遍历攻击
func containsPathTraversal(input string) bool {
	pathTraversalPatterns := []string{
		"../", "..\\", "~/", "~\\",
		"%2e%2e", "%252e", "..%2f", "..%5c",
		"/etc/", "/proc/", "c:\\", "d:\\",
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range pathTraversalPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}

	return false
}

// filterDangerousChars 过滤危险字符
func filterDangerousChars(input string) string {
	// 允许的字符：字母、数字、常见符号、中文
	cleaned := strings.Map(func(r rune) rune {
		// 允许 ASCII 字符
		if r < 128 {
			// 允许字母、数字、常见标点
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				return r
			}
			// 允许常见安全字符
			switch r {
			case ' ', '\t', '\n', '\r', // 空白字符
				'!', '@', '#', '$', '%', '^', '&', '*', '(', ')',
				'+', '=', '-', '_', '{', '}', '[', ']', '|',
				'\\', ':', ';', '"', '\'', '<', '>', ',', '.', '?',
				'/', '~', '`':
				return r
			default:
				// 其他 ASCII 字符被过滤
				return -1
			}
		}
		// 允许中文字符和其他 Unicode 字符
		if r >= 0x4E00 && r <= 0x9FFF { // CJK 统一汉字
			return r
		}
		// 允许常见 Unicode 标点
		if r >= 0x2000 && r <= 0x206F { // 通用标点
			return r
		}
		// 其他字符被过滤
		return -1
	}, input)

	return cleaned
}

// isWholeWord 检查是否是完整单词
func isWholeWord(input, word string) bool {
	// 使用正则表达式匹配完整单词
	pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(word) + `\b`)
	return pattern.MatchString(input)
}

// setSecurityHeaders 设置安全响应头
func setSecurityHeaders(c *gin.Context) {
	// 防止点击劫持
	c.Header("X-Frame-Options", "DENY")

	// 防止 MIME 类型嗅探
	c.Header("X-Content-Type-Options", "nosniff")

	// XSS 保护
	c.Header("X-XSS-Protection", "1; mode=block")

	// 内容安全策略（调整为适合 SPA 应用）
	// 允许同源、内联脚本和内联样式
	// 移动端需要从 alicdn.com 加载字体
	c.Header("Content-Security-Policy",
		"default-src 'self'; "+
		"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
		"style-src 'self' 'unsafe-inline' https://at.alicdn.com; "+
		"img-src 'self' data: blob: https:; "+
		"font-src 'self' data: https://at.alicdn.com; "+
		"connect-src 'self' https:; "+
		"frame-ancestors 'self';")

	// 严格传输安全（如果是 HTTPS）
	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}

	// 隐藏服务器信息
	c.Header("Server", "")
}

// ValidateID 验证 ID 参数（防止 ID 注入）
func ValidateID(c *gin.Context, param string) (uint, error) {
	idStr := c.Param(param)
	if idStr == "" {
		return 0, fmt.Errorf("参数 %s 不能为空", param)
	}

	// 检查是否为纯数字
	matched, _ := regexp.MatchString(`^\d+$`, idStr)
	if !matched {
		return 0, fmt.Errorf("参数 %s 必须是数字", param)
	}

	// 转换为数字并检查范围
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("参数 %s 格式错误", param)
	}

	if id > 4294967295 { // uint 最大值
		return 0, fmt.Errorf("参数 %s 值过大", param)
	}

	return uint(id), nil
}

// ValidatePage 验证分页参数
func ValidatePage(c *gin.Context) (page, perPage int) {
	// 获取分页参数，设置默认值
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")

	// 验证并转换 page
	if matched, _ := regexp.MatchString(`^\d+$`, pageStr); matched {
		page, _ = strconv.Atoi(pageStr)
	}
	if page < 1 {
		page = 1
	}
	if page > 10000 { // 防止过大值
		page = 10000
	}

	// 验证并转换 perPage
	if matched, _ := regexp.MatchString(`^\d+$`, perPageStr); matched {
		perPage, _ = strconv.Atoi(perPageStr)
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 { // 限制最大每页数量
		perPage = 100
	}

	return page, perPage
}

// SanitizeSQLLike 为 SQL LIKE 查询清理输入
func SanitizeSQLLike(input string) string {
	// 移除 LIKE 语句中的特殊字符
	input = strings.ReplaceAll(input, "\\", "\\\\")
	input = strings.ReplaceAll(input, "%", "\\%")
	input = strings.ReplaceAll(input, "_", "\\_")
	input = strings.ReplaceAll(input, "'", "''")
	return input
}

// RateLimitMiddleware 简单的速率限制中间件
// 使用 map 存储请求计数（生产环境建议使用 Redis）
func RateLimitMiddleware(requests int, duration int) gin.HandlerFunc {
	// 这里简化实现，生产环境应使用 Redis
	return func(c *gin.Context) {
		// TODO: 实现基于 Redis 的速率限制
		c.Next()
	}
}

// GetSafeString 安全地获取字符串参数
func GetSafeString(c *gin.Context, key, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	// 清理输入
	cleaned, err := sanitizeInput(value)
	if err != nil {
		return defaultValue
	}

	return cleaned
}

// GetSafeInt 安全地获取整数参数
func GetSafeInt(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	// 验证是否为纯数字
	matched, _ := regexp.MatchString(`^-?\d+$`, value)
	if !matched {
		return defaultValue
	}

	// 转换并限制范围
	num, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	// 限制范围
	if num < -2147483648 || num > 2147483647 {
		return defaultValue
	}

	return num
}

// ValidateJSON 验证 JSON 输入大小
func ValidateJSON(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" && c.Request.Method != "PUT" {
			c.Next()
			return
		}

		// 检查 Content-Length
		if c.Request.ContentLength > maxSize {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"success": false,
				"error":   "请求体过大",
			})
			return
		}

		c.Next()
	}
}

// LogSecurityEvent 记录安全事件（可扩展为发送到监控系统）
func LogSecurityEvent(c *gin.Context, eventType string, details string) {
	// 记录到日志
	fmt.Printf("[SECURITY] %s - IP: %s, Path: %s, Details: %s\n",
		eventType,
		c.ClientIP(),
		c.Request.URL.Path,
		details,
	)

	// TODO: 可以扩展为：
	// 1. 发送到安全监控系统
	// 2. 记录到数据库
	// 3. 触发告警
}
