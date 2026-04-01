package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	"github.com/yourorg/material-backend/backend/internal/config"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// wechatLoginReq 微信登录请求
type wechatLoginReq struct {
	Code string `json:"code" binding:"required"` // 微信登录凭证
}

// wechatBindReq 绑定微信请求
type wechatBindReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

// wechatSessionResponse 微信 session 响应
type wechatSessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// RegisterWechatRoutes 注册微信登录相关路由
func RegisterWechatRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/auth/wechat-login", wechatLoginHandler(db))
	r.POST("/auth/bind-wechat", wechatBindHandler(db))
}

// wechatLoginHandler 微信登录处理
func wechatLoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req wechatLoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "缺少 code 参数")
			return
		}

		// 获取微信 openid
		openID, err := getWechatOpenID(req.Code)
		if err != nil {
			response.InternalError(c, "微信登录失败: "+err.Error())
			return
		}

		// 查找绑定了该 openid 的用户
		var user User
		if err := db.Preload("Roles").Where("wechat_open_id = ?", openID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 用户未绑定，返回需要绑定
				response.Success(c, map[string]interface{}{
					"need_bind": true,
					"openid":    openID,
					"message":   "该微信未绑定账号，请先绑定",
				})
				return
			}
			response.InternalError(c, "查询用户失败")
			return
		}

		if !user.IsActive {
			response.Unauthorized(c, "账户已被禁用")
			return
		}

		// 更新最后登录时间
		now := time.Now().UTC()
		db.Model(&User{}).Where("id = ?", user.ID).Update("last_login", now)

		// 生成 token
		token, err := jwtpkg.GenerateToken(user.ID, user.Username)
		if err != nil {
			response.InternalError(c, "token 创建失败")
			return
		}

		response.SuccessWithMeta(c, user.ToDTO(), map[string]interface{}{
			"token": token,
		})
	}
}

// wechatBindHandler 绑定微信处理
func wechatBindHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req wechatBindReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "参数错误")
			return
		}

		// 获取微信 openid
		openID, err := getWechatOpenID(req.Code)
		if err != nil {
			response.InternalError(c, "微信登录失败: "+err.Error())
			return
		}

		// 验证用户名密码
		var user User
		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.Unauthorized(c, "用户名或密码错误")
				return
			}
			response.InternalError(c, "查询用户失败")
			return
		}

		if !user.CheckPassword(req.Password) {
			response.Unauthorized(c, "用户名或密码错误")
			return
		}

		if !user.IsActive {
			response.Unauthorized(c, "账户已被禁用")
			return
		}

		// 检查该 openid 是否已被其他用户绑定
		var existUser User
		if err := db.Where("wechat_open_id = ? AND id != ?", openID, user.ID).First(&existUser).Error; err == nil {
			response.BadRequest(c, "该微信已绑定其他账号")
			return
		}

		// 绑定 openid
		if err := db.Model(&user).Update("wechat_open_id", openID).Error; err != nil {
			response.InternalError(c, "绑定失败")
			return
		}

		// 重新加载用户信息
		db.Preload("Roles").First(&user, user.ID)

		// 更新最后登录时间
		now := time.Now().UTC()
		db.Model(&User{}).Where("id = ?", user.ID).Update("last_login", now)

		// 生成 token
		token, err := jwtpkg.GenerateToken(user.ID, user.Username)
		if err != nil {
			response.InternalError(c, "token 创建失败")
			return
		}

		response.SuccessWithMeta(c, user.ToDTO(), map[string]interface{}{
			"token": token,
		})
	}
}

// getWechatOpenID 通过 code 获取微信 openid
func getWechatOpenID(code string) (string, error) {
	cfg := config.Get()
	appID := cfg.Wechat.AppID
	appSecret := cfg.Wechat.AppSecret

	if appID == "" || appSecret == "" {
		return "", fmt.Errorf("微信小程序配置缺失")
	}

	// 调用微信接口获取 openid
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求微信接口失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var sessionResp wechatSessionResponse
	if err := json.Unmarshal(body, &sessionResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if sessionResp.ErrCode != 0 {
		return "", fmt.Errorf("微信接口错误: %d - %s", sessionResp.ErrCode, sessionResp.ErrMsg)
	}

	return sessionResp.OpenID, nil
}
