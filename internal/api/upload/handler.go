package upload

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"github.com/yourorg/material-backend/backend/internal/api/system"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// UploadConfig 上传配置
type UploadConfig struct {
	UploadDirectory string
	MaxFileSize     int64    // bytes
	AllowedTypes    map[string]bool
	MaxUploadCount  int
}

// GetUploadConfig 从数据库获取上传配置
func GetUploadConfig(db *gorm.DB) *UploadConfig {
	config := &UploadConfig{
		UploadDirectory: "static/uploads",
		MaxFileSize:     5 * 1024 * 1024, // 5MB
		AllowedTypes: map[string]bool{
			"png":  true,
			"jpg":  true,
			"jpeg": true,
			"gif":  true,
			"bmp":  true,
			"webp": true,
			"svg":  true,
		},
		MaxUploadCount: 10,
	}

	// 从数据库读取配置
	var systemConfigs []system.SystemConfig
	if err := db.Where("key IN ?",
		[]string{"upload_directory", "max_file_size", "allowed_file_types", "max_upload_count"}).
		Find(&systemConfigs).Error; err == nil {
		for _, sc := range systemConfigs {
			switch sc.Key {
			case "upload_directory":
				if sc.Value != "" {
					config.UploadDirectory = sc.Value
				}
			case "max_file_size":
				size := 5
				if sc.Value != "" {
					fmt.Sscanf(sc.Value, "%d", &size)
				}
				config.MaxFileSize = int64(size) * 1024 * 1024
			case "allowed_file_types":
				if sc.Value != "" {
					types := strings.Split(sc.Value, ",")
					config.AllowedTypes = make(map[string]bool)
					for _, t := range types {
						config.AllowedTypes[strings.TrimSpace(t)] = true
					}
				}
			case "max_upload_count":
				count := 10
				if sc.Value != "" {
					fmt.Sscanf(sc.Value, "%d", &count)
				}
				config.MaxUploadCount = count
			}
		}
	}

	return config
}

// RegisterRoutes 注册上传模块路由
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// 创建路由组
	r := rg.Group("/upload")
	// 使用JWT中间件进行身份验证
	r.Use(jwtpkg.TokenMiddleware())

	// 获取上传配置
	uploadConfig := GetUploadConfig(db)
	uploadFolder := uploadConfig.UploadDirectory

	// 确保上传目录存在
	if err := os.MkdirAll(uploadFolder, 0755); err != nil {
		fmt.Printf("创建上传目录失败: %v\n", err)
	}

	// 检查文件是否允许上传
	allowedFile := func(filename string) bool {
		ext := strings.ToLower(filepath.Ext(filename))
		if ext == "" || !strings.HasPrefix(ext, ".") {
			return false
		}
		ext = ext[1:] // 去掉点号
		_, ok := uploadConfig.AllowedTypes[ext]
		return ok
	}

	// ================== 通用图片上传接口 ==================
	r.POST("/image", func(c *gin.Context) {
		// 获取当前用户（可选，用于记录）
		_, _ = auth.GetCurrentUser(c, db)

		// 获取上传文件
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			response.BadRequest(c, "未选择文件")
			return
		}
		defer file.Close()

		// 检查文件名
		if header.Filename == "" {
			response.BadRequest(c, "未选择文件")
			return
		}

		// 检查文件扩展名
		if !allowedFile(header.Filename) {
			response.BadRequest(c, "文件类型不支持，仅支持图片格式（png, jpg, jpeg, gif, bmp, webp, svg）")
			return
		}

		// 生成唯一文件名（使用当前时间戳+随机数+文件扩展名）
		timestamp := time.Now().UnixNano() / int64(time.Millisecond)
		randomNum := timestamp % 10000
		ext := strings.ToLower(filepath.Ext(header.Filename))
		filename := fmt.Sprintf("%d_%d%s", timestamp, randomNum, ext)
		filePath := filepath.Join(uploadFolder, filename)

		// 确保上传目录存在（每次上传前都检查）
		if err := os.MkdirAll(uploadFolder, 0755); err != nil {
			response.InternalError(c, "创建上传目录失败")
			return
		}

		// 创建目标文件
		dst, err := os.Create(filePath)
		if err != nil {
			response.InternalError(c, "创建文件失败")
			return
		}
		defer dst.Close()

		// 复制文件内容
		written, err := io.Copy(dst, file)
		if err != nil {
			response.InternalError(c, "保存文件失败")
			return
		}
		fmt.Printf("文件上传成功 [路径: %s, 大小: %d 字节]\n", filePath, written)

		// 生成文件URL（相对路径）
		url := fmt.Sprintf("/static/uploads/%s", filename)

		response.SuccessWithMeta(c, map[string]interface{}{
			"url":      url,
			"filename": filename,
			"size":     header.Size,
		}, map[string]interface{}{
			"message": "上传成功",
		})
	})

	// ================== 批量图片上传接口 ==================
	r.POST("/images", func(c *gin.Context) {
		// 获取当前用户（可选，用于记录）
		_, _ = auth.GetCurrentUser(c, db)

		// 解析多部分表单
		form, err := c.MultipartForm()
		if err != nil {
			response.BadRequest(c, "解析表单失败")
			return
		}

		files := form.File["files"]
		if len(files) == 0 {
			response.BadRequest(c, "未选择文件")
			return
		}

		// 限制批量上传数量
		if len(files) > uploadConfig.MaxUploadCount {
			response.BadRequest(c, fmt.Sprintf("最多支持同时上传%d张图片", uploadConfig.MaxUploadCount))
			return
		}

		uploadedFiles := []map[string]interface{}{}
		failedFiles := []string{}

		for _, fileHeader := range files {
			// 检查文件扩展名
			if !allowedFile(fileHeader.Filename) {
				failedFiles = append(failedFiles, fmt.Sprintf("%s（文件类型不支持）", fileHeader.Filename))
				continue
			}

			// 打开上传的文件
			file, err := fileHeader.Open()
			if err != nil {
				failedFiles = append(failedFiles, fmt.Sprintf("%s（打开文件失败）", fileHeader.Filename))
				continue
			}

			// 生成唯一文件名
			timestamp := time.Now().UnixNano() / int64(time.Millisecond)
			randomNum := timestamp % 10000
			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
			filename := fmt.Sprintf("%d_%d%s", timestamp, randomNum, ext)
			filePath := filepath.Join(uploadFolder, filename)

			// 确保上传目录存在（每次上传前都检查）
			if err := os.MkdirAll(uploadFolder, 0755); err != nil {
				failedFiles = append(failedFiles, fmt.Sprintf("%s（创建上传目录失败）", fileHeader.Filename))
				file.Close()
				continue
			}

			// 创建目标文件
			dst, err := os.Create(filePath)
			if err != nil {
				file.Close()
				failedFiles = append(failedFiles, fmt.Sprintf("%s（创建文件失败）", fileHeader.Filename))
				continue
			}

			// 复制文件内容
			_, err = io.Copy(dst, file)
			dst.Close()
			file.Close()

			if err != nil {
				failedFiles = append(failedFiles, fmt.Sprintf("%s（保存失败）", fileHeader.Filename))
				continue
			}

			// 生成文件URL
			url := fmt.Sprintf("/static/uploads/%s", filename)
			uploadedFiles = append(uploadedFiles, map[string]interface{}{
				"url":      url,
				"filename": filename,
				"size":     fileHeader.Size,
			})
		}

		response.SuccessWithMessage(c, map[string]interface{}{
			"uploaded": uploadedFiles,
			"failed":   failedFiles,
		}, fmt.Sprintf("上传完成，成功%d个，失败%d个", len(uploadedFiles), len(failedFiles)))
	})
}
