package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// checkVersion 检查应用版本
func (h *Handler) CheckVersion(c *gin.Context) {
	platform := c.Query("platform")
	if platform == "" {
		platform = "android"
	}

	currentVersion := c.Query("current_version")

	// 获取该平台的最新版本
	var latest AppVersion
	err := h.DB.Where("platform = ?", platform).
		Order("published_at DESC, created_at DESC").
		First(&latest).Error

	if err != nil {
		// 如果没有版本信息，返回无需更新
		response.Success(c, VersionCheckResponse{
			HasUpdate:     false,
			LatestVersion: "",
			DownloadURL:   "",
			ForceUpdate:   false,
		})
		return
	}

	// 比较版本号
	hasUpdate := compareVersions(currentVersion, latest.Version)

	// 生成下载 URL，使用带版本号的文件名
	// 添加时间戳参数避免浏览器缓存
	downloadURL := fmt.Sprintf("https://home.mbed.org.cn:9090/mobile-updates/%s/material-management-%s.apk?t=%d",
		platform, latest.Version, time.Now().Unix())

	response.Success(c, VersionCheckResponse{
		HasUpdate:     hasUpdate,
		LatestVersion: latest.Version,
		DownloadURL:   downloadURL,
		ForceUpdate:   latest.ForceUpdate,
		UpdateMessage: latest.UpdateMessage,
		ReleaseNotes:  latest.ReleaseNotes,
	})
}

// DownloadAPK 下载 APK 文件
func (h *Handler) DownloadAPK(c *gin.Context) {
	version := c.Query("version")
	platform := c.Query("platform")
	if platform == "" {
		platform = "android"
	}

	// 构建文件路径（使用绝对路径）
	baseDir, _ := os.Getwd()
	var filename string
	if version != "" {
		filename = filepath.Join(baseDir, "mobile-app-updates", platform, "material-management-"+version+".apk")
	} else {
		// 如果没有指定版本，从数据库获取最新版本
		var latest AppVersion
		err := h.DB.Where("platform = ?", platform).
			Order("published_at DESC, created_at DESC").
			First(&latest).Error
		if err == nil {
			filename = filepath.Join(baseDir, "mobile-app-updates", platform, "material-management-"+latest.Version+".apk")
		} else {
			// 回退到 latest.apk
			filename = filepath.Join(baseDir, "mobile-app-updates", platform, "latest.apk")
		}
	}

	// 解析符号链接
	realPath, err := filepath.EvalSymlinks(filename)
	if err == nil && realPath != "" {
		filename = realPath
	}

	// 检查文件是否存在
	fileInfo, err := os.Stat(filename)
	if err != nil || os.IsNotExist(err) {
		response.Error(c, http.StatusNotFound, "文件不存在: version="+version+", file="+filename)
		return
	}

	// 验证是文件而不是目录
	if fileInfo.IsDir() {
		response.Error(c, http.StatusNotFound, "路径是目录而非文件")
		return
	}

	// 设置下载的文件名
	downloadName := "material-management-" + version + ".apk"
	if version == "" {
		downloadName = "material-management.apk"
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Header("Content-Type", "application/vnd.android.package-archive")
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 发送文件
	c.File(filename)
}

// GetLatestAPKInfo 获取最新 APK 信息（返回 JSON）
func (h *Handler) GetLatestAPKInfo(c *gin.Context) {
	platform := c.Query("platform")
	if platform == "" {
		platform = "android"
	}

	// 获取该平台的最新版本
	var latest AppVersion
	err := h.DB.Where("platform = ?", platform).
		Order("published_at DESC, created_at DESC").
		First(&latest).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "未找到版本信息",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"platform":       latest.Platform,
		"version":        latest.Version,
		"download_url":   latest.DownloadURL,
		"api_download_url": "/api/app/download-apk?platform=" + platform + "&version=" + latest.Version,
		"update_message": latest.UpdateMessage,
		"release_notes":  latest.ReleaseNotes,
		"published_at":   latest.PublishedAt,
	})
}

// compareVersions 比较两个版本号
// 返回 true 表示 v2 比 v1 新
func compareVersions(v1, v2 string) bool {
	if v1 == "" {
		return true
	}

	// 简单的版本号比较
	// 例如: 1.0.0 < 1.0.1 < 1.1.0 < 2.0.0
	v1Parts := parseVersion(v1)
	v2Parts := parseVersion(v2)

	for i := 0; i < 3; i++ {
		if v2Parts[i] > v1Parts[i] {
			return true
		}
		if v2Parts[i] < v1Parts[i] {
			return false
		}
	}

	return false // 版本相同
}

func parseVersion(version string) [3]int {
	var parts [3]int
	current := 0
	partIndex := 0

	for _, ch := range version {
		if ch >= '0' && ch <= '9' {
			current = current*10 + int(ch-'0')
		} else if ch == '.' && partIndex < 2 {
			parts[partIndex] = current
			partIndex++
			current = 0
		}
	}

	if partIndex < 3 {
		parts[partIndex] = current
	}

	return parts
}
