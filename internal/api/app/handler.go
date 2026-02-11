package app

import (
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

	response.Success(c, VersionCheckResponse{
		HasUpdate:     hasUpdate,
		LatestVersion: latest.Version,
		DownloadURL:   latest.DownloadURL,
		ForceUpdate:   latest.ForceUpdate,
		UpdateMessage: latest.UpdateMessage,
		ReleaseNotes:  latest.ReleaseNotes,
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
