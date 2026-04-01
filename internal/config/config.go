package config

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yourorg/material-backend/backend/internal/pkg/logger"
	"gopkg.in/yaml.v3"
)

// Config 应用程序配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Upload   UploadConfig   `yaml:"upload"`
	Log      LogConfig      `yaml:"log"`
	AI       AIConfig       `yaml:"ai"`
	CORS     CORSConfig     `yaml:"cors"`
	Wechat   WechatConfig   `yaml:"wechat"`
}

// WechatConfig 微信小程序配置
type WechatConfig struct {
	AppID     string `yaml:"app_id"`     // 小程序 AppID
	AppSecret string `yaml:"app_secret"` // 小程序 AppSecret
}

// AIConfig AI 配置
type AIConfig struct {
	// DeepSeek 配置
	DeepSeekAPIKey  string `yaml:"deepseek_api_key"`
	DeepSeekModel   string `yaml:"deepseek_model"`
	DeepSeekBaseURL string `yaml:"deepseek_base_url"`

	// 百度千帆配置 (Anthropic 兼容 API)
	BaiduAPIKey    string `yaml:"baidu_api_key"`     // 百度千帆 Auth Token (bce-v3/...)
	BaiduSecretKey string `yaml:"baidu_secret_key"`  // 百度千帆 Secret Key (原生 API 用)
	BaiduModel     string `yaml:"baidu_model"`       // 模型名称: glm-5, coding-plan 等
	BaiduBaseURL   string `yaml:"baidu_base_url"`    // API 基础 URL

	// OpenAI 配置 (可选，用于语音转文字)
	OpenAIAPIKey string `yaml:"openai_api_key"`

	// 本地 ASR 服务配置
	ASRServiceURL string `yaml:"asr_service_url"`
	ASREnabled    bool   `yaml:"asr_enabled"`

	// 默认使用的模型提供者: "baidu" 或 "deepseek"
	DefaultProvider string `yaml:"default_provider"`

	// Agent 安全配置
	AgentSecret string `yaml:"agent_secret"` // AI Agent 请求签名密钥
}

// AIProviderConfig 单个 AI 提供者的配置
type AIProviderConfig struct {
	Name    string `json:"name"`
	APIKey  string `json:"api_key,omitempty"`  // 不暴露敏感信息
	Model   string `json:"model"`
	BaseURL string `json:"base_url"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            int           `yaml:"port"`
	Mode            string        `yaml:"mode"`            // debug, release, test
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `yaml:"type"`     // postgresql, mysql, sqlite
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`  // disable, require, verify-ca, verify-full

	// 连接池配置
	MaxIdleConns int `yaml:"max_idle_conns"`
	MaxOpenConns int `yaml:"max_open_conns"`
	MaxLifetime  int `yaml:"max_lifetime"` // 秒
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `yaml:"secret"`
	ExpireTime time.Duration `yaml:"expire_time"`
	Issuer     string        `yaml:"issuer"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxFileSize      int64  `yaml:"max_file_size"`       // 字节
	MaxUploadCount   int    `yaml:"max_upload_count"`
	AllowedTypes     string `yaml:"allowed_types"`
	UploadDir        string `yaml:"upload_dir"`
	GenerateFileName bool   `yaml:"generate_file_name"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`       // debug, info, warn, error
	FileName   string `yaml:"file_name"`
	MaxSize    int    `yaml:"max_size"`    // MB
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`     // days
	Compress   bool   `yaml:"compress"`
}

// CORSConfig CORS 配置
type CORSConfig struct {
	Enabled         bool     `yaml:"enabled"`           // 是否启用 CORS 白名单
	AllowedOrigins  []string `yaml:"allowed_origins"`   // 允许的域名列表
	AllowCredentials bool    `yaml:"allow_credentials"` // 是否允许携带凭证
}

var (
	cfg     *Config
	cfgFile string
)

// Load 从文件加载配置
func Load(configFile string) (*Config, error) {
	cfgFile = configFile
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	cfg = &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 设置默认值
	setDefaults()

	// 验证配置
	if err := validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadOrDefault 从文件加载配置，如果失败则使用默认配置
func LoadOrDefault(configFile string) *Config {
	cfg, err := Load(configFile)
	if err != nil {
		// 使用默认配置
		cfg = &Config{}
		setDefaults()
	}
	return cfg
}

// Get 获取当前配置
func Get() *Config {
	return cfg
}

// Watch 监听配置文件变化
func Watch(onChange func(*Config)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if newCfg, err := Load(cfgFile); err == nil {
						cfg = newCfg
						onChange(cfg)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Warnf("配置文件监听错误: %v", err)
			}
		}
	}()

	return watcher.Add(cfgFile)
}

// setDefaults 设置默认值
func setDefaults() {
	// 服务器默认配置
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8088
	}
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "debug"
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 60 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 60 * time.Second
	}
	if cfg.Server.ShutdownTimeout == 0 {
		cfg.Server.ShutdownTimeout = 10 * time.Second
	}

	// 数据库默认配置
	if cfg.Database.Type == "" {
		cfg.Database.Type = "postgresql"
	}
	if cfg.Database.Host == "" {
		cfg.Database.Host = "127.0.0.1"
	}
	if cfg.Database.Port == 0 {
		cfg.Database.Port = 5432
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 10
	}
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 100
	}
	if cfg.Database.MaxLifetime == 0 {
		cfg.Database.MaxLifetime = 3600 // 1小时
	}

	// JWT默认配置
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "your-secret-key-change-in-production"
	}
	if cfg.JWT.ExpireTime == 0 {
		cfg.JWT.ExpireTime = 24 * time.Hour
	}
	if cfg.JWT.Issuer == "" {
		cfg.JWT.Issuer = "material-backend"
	}

	// 上传默认配置
	if cfg.Upload.MaxFileSize == 0 {
		cfg.Upload.MaxFileSize = 5 * 1024 * 1024 // 5MB
	}
	if cfg.Upload.MaxUploadCount == 0 {
		cfg.Upload.MaxUploadCount = 10
	}
	if cfg.Upload.AllowedTypes == "" {
		cfg.Upload.AllowedTypes = "jpg,jpeg,png,gif,bmp,webp,svg,pdf,doc,docx,xls,xlsx"
	}
	if cfg.Upload.UploadDir == "" {
		cfg.Upload.UploadDir = "static/uploads"
	}
	if cfg.Upload.GenerateFileName {
		cfg.Upload.GenerateFileName = true
	}

	// 日志默认配置
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.MaxSize == 0 {
		cfg.Log.MaxSize = 100 // MB
	}
	if cfg.Log.MaxBackups == 0 {
		cfg.Log.MaxBackups = 3
	}
	if cfg.Log.MaxAge == 0 {
		cfg.Log.MaxAge = 7 // days
	}

	// AI 默认配置
	if cfg.AI.DeepSeekModel == "" {
		cfg.AI.DeepSeekModel = "deepseek-chat"
	}
	if cfg.AI.DeepSeekBaseURL == "" {
		cfg.AI.DeepSeekBaseURL = "https://api.deepseek.com/v1"
	}
	// 百度千帆默认配置
	if cfg.AI.BaiduModel == "" {
		cfg.AI.BaiduModel = "coding-plan" // 默认使用 coding-plan 模型
	}
	if cfg.AI.BaiduBaseURL == "" {
		cfg.AI.BaiduBaseURL = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat"
	}
	if cfg.AI.ASRServiceURL == "" {
		cfg.AI.ASRServiceURL = "http://localhost:8089"
	}
}

// validate 验证配置
func validate() error {
	if cfg.Server.Port < 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("服务器端口配置错误: %d", cfg.Server.Port)
	}

	if cfg.Database.Type != "postgresql" && cfg.Database.Type != "mysql" && cfg.Database.Type != "sqlite" {
		return fmt.Errorf("不支持的数据库类型: %s", cfg.Database.Type)
	}

	if cfg.Database.User == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}

	if cfg.Database.Database == "" {
		return fmt.Errorf("数据库名不能为空")
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	switch c.Type {
	case "postgresql":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User, c.Password, c.Host, c.Port, c.Database)
	case "sqlite":
		return c.Database
	default:
		return ""
	}
}

// GetMaxFileSizeMB 获取最大文件大小(MB)
func (c *UploadConfig) GetMaxFileSizeMB() int {
	return int(c.MaxFileSize / 1024 / 1024)
}
