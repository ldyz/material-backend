package upload

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
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

// 压缩配置
const (
	maxWidth       = 1920  // 最大宽度
	maxHeight      = 1080  // 最大高度
	jpegQuality    = 80    // JPEG 压缩质量 (1-100)
	compressThreshold = 500 * 1024 // 超过 500KB 才压缩
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

// compressImage 压缩图片
// 返回压缩后的数据和文件扩展名
func compressImage(data []byte, ext string, originalSize int64) ([]byte, string, error) {
	// 小于阈值不压缩
	if originalSize < compressThreshold {
		return data, ext, nil
	}

	ext = strings.ToLower(strings.TrimPrefix(ext, "."))

	// SVG 不压缩
	if ext == "svg" {
		return data, ext, nil
	}

	// 解码图片
	var img image.Image
	var err error
	reader := bytes.NewReader(data)

	switch ext {
	case "jpg", "jpeg":
		img, err = jpeg.Decode(reader)
	case "png":
		img, err = png.Decode(reader)
	default:
		// 其他格式尝试用通用解码器
		img, _, err = image.Decode(reader)
	}

	if err != nil {
		// 解码失败，返回原始数据
		fmt.Printf("图片解码失败，保留原文件: %v\n", err)
		return data, ext, nil
	}

	// 检查尺寸，如果太大则缩放
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// 计算缩放比例
	scale := 1.0
	if imgWidth > maxWidth {
		scale = float64(maxWidth) / float64(imgWidth)
	}
	if imgHeight > maxHeight {
		hScale := float64(maxHeight) / float64(imgHeight)
		if hScale < scale {
			scale = hScale
		}
	}

	// 如果需要缩放
	if scale < 1.0 {
		newWidth := int(float64(imgWidth) * scale)
		newHeight := int(float64(imgHeight) * scale)
		img = resizeImage(img, newWidth, newHeight)
		fmt.Printf("图片缩放: %dx%d -> %dx%d\n", imgWidth, imgHeight, newWidth, newHeight)
	}

	// 编码为 JPEG（压缩效果最好）
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality}); err != nil {
		fmt.Printf("JPEG 编码失败: %v\n", err)
		return data, ext, nil
	}

	// 如果压缩后反而更大，返回原始数据
	if buf.Len() >= len(data) {
		fmt.Printf("压缩后大小(%d) >= 原始大小(%d)，保留原文件\n", buf.Len(), len(data))
		return data, ext, nil
	}

	fmt.Printf("图片压缩成功: %d -> %d 字节 (压缩率: %.1f%%)\n",
		originalSize, buf.Len(), float64(buf.Len())/float64(originalSize)*100)

	return buf.Bytes(), "jpg", nil
}

// resizeImage 简单的图片缩放（使用最近邻插值）
func resizeImage(img image.Image, newWidth, newHeight int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	// 创建新图片
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 简单的最近邻插值
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := x * srcWidth / newWidth
			srcY := y * srcHeight / newHeight
			dst.Set(x, y, img.At(srcX+bounds.Min.X, srcY+bounds.Min.Y))
		}
	}

	return dst
}

// processAndSaveImage 处理并保存图片（压缩）
func processAndSaveImage(fileData []byte, ext string, uploadFolder string, originalSize int64) (string, int64, error) {
	// 压缩图片
	compressedData, newExt, err := compressImage(fileData, ext, originalSize)
	if err != nil {
		return "", 0, err
	}

	// 生成唯一文件名
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	randomNum := timestamp % 10000
	filename := fmt.Sprintf("%d_%d.%s", timestamp, randomNum, newExt)
	filePath := filepath.Join(uploadFolder, filename)

	// 保存文件
	if err := os.WriteFile(filePath, compressedData, 0644); err != nil {
		return "", 0, err
	}

	return filename, int64(len(compressedData)), nil
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

		// 确保上传目录存在（每次上传前都检查）
		if err := os.MkdirAll(uploadFolder, 0755); err != nil {
			response.InternalError(c, "创建上传目录失败")
			return
		}

		// 读取文件内容
		fileData, err := io.ReadAll(file)
		if err != nil {
			response.InternalError(c, "读取文件失败")
			return
		}

		ext := strings.ToLower(filepath.Ext(header.Filename))

		// 处理并保存图片（压缩）
		filename, compressedSize, err := processAndSaveImage(fileData, ext, uploadFolder, header.Size)
		if err != nil {
			response.InternalError(c, "保存文件失败")
			return
		}

		fmt.Printf("文件上传成功 [文件: %s, 原始: %d 字节, 压缩后: %d 字节]\n",
			filename, header.Size, compressedSize)

		// 生成文件URL（相对路径）
		url := fmt.Sprintf("/uploads/%s", filename)

		response.SuccessWithMeta(c, map[string]interface{}{
			"url":           url,
			"filename":      filename,
			"original_size": header.Size,
			"compressed_size": compressedSize,
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

			// 读取文件内容
			fileData, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				failedFiles = append(failedFiles, fmt.Sprintf("%s（读取文件失败）", fileHeader.Filename))
				continue
			}

			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

			// 处理并保存图片（压缩）
			filename, compressedSize, err := processAndSaveImage(fileData, ext, uploadFolder, fileHeader.Size)
			if err != nil {
				failedFiles = append(failedFiles, fmt.Sprintf("%s（保存失败）", fileHeader.Filename))
				continue
			}

			// 生成文件URL
			url := fmt.Sprintf("/uploads/%s", filename)
			uploadedFiles = append(uploadedFiles, map[string]interface{}{
				"url":             url,
				"filename":        filename,
				"original_size":   fileHeader.Size,
				"compressed_size": compressedSize,
			})

			fmt.Printf("批量上传成功 [文件: %s, 原始: %d 字节, 压缩后: %d 字节]\n",
				filename, fileHeader.Size, compressedSize)
		}

		response.SuccessWithMessage(c, map[string]interface{}{
			"uploaded": uploadedFiles,
			"failed":   failedFiles,
		}, fmt.Sprintf("上传完成，成功%d个，失败%d个", len(uploadedFiles), len(failedFiles)))
	})
}
