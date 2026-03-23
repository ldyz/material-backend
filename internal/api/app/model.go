package app

import "time"

// AppVersion 应用版本信息
type AppVersion struct {
	ID            uint      `json:"id"`
	Platform      string    `json:"platform"`       // android, ios
	Version       string    `json:"version"`        // 版本号，如 1.0.0
	BuildNumber   int       `json:"build_number"`   // 构建号
	DownloadURL   string    `json:"download_url"`   // 下载链接
	ForceUpdate   bool      `json:"force_update"`   // 是否强制更新
	UpdateMessage string    `json:"update_message"` // 更新提示
	ReleaseNotes  string    `json:"release_notes"`  // 更新日志
	PublishedAt   time.Time `json:"published_at"`   // 发布时间
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// VersionCheckResponse 版本检查响应
type VersionCheckResponse struct {
	HasUpdate     bool   `json:"has_update"`      // 是否有更新
	LatestVersion string `json:"latest_version"`   // 最新版本号
	DownloadURL   string `json:"download_url"`     // 下载地址
	ForceUpdate   bool   `json:"force_update"`     // 是否强制更新
	UpdateMessage string `json:"update_message"`   // 更新提示
	ReleaseNotes  string `json:"release_notes"`    // 更新日志
}
