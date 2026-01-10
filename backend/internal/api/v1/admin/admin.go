package admin

import (
	"app-platform-backend/internal/middleware"
	"app-platform-backend/internal/model"
	"app-platform-backend/internal/pkg/database"
	"app-platform-backend/internal/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 管理员登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "用户名和密码不能为空")
		return
	}

	// 验证用户名长度
	if len(req.Username) < 3 || len(req.Username) > 50 {
		response.ParamError(c, "用户名长度应在3-50个字符之间")
		return
	}

	var admin model.Admin
	if err := database.GetDB().Where("username = ?", req.Username).First(&admin).Error; err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	token, err := middleware.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		response.InternalError(c, "生成令牌失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"nickname": admin.Nickname,
			"avatar":   admin.Avatar,
		},
	})
}

// GetInfo 获取管理员信息
func GetInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var admin model.Admin
	if err := database.GetDB().First(&admin, userID).Error; err != nil {
		response.NotFound(c, "管理员不存在")
		return
	}

	response.Success(c, gin.H{
		"id":       admin.ID,
		"username": admin.Username,
		"nickname": admin.Nickname,
		"avatar":   admin.Avatar,
	})
}

// Logout 管理员登出
func Logout(c *gin.Context) {
	response.SuccessWithMessage(c, nil, "登出成功")
}

// UpdatePassword 更新密码
func UpdatePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "请输入旧密码和新密码")
		return
	}

	// 验证新密码长度
	if len(req.NewPassword) < 6 {
		response.ParamError(c, "新密码长度至少6个字符")
		return
	}

	if len(req.NewPassword) > 100 {
		response.ParamError(c, "新密码长度不能超过100个字符")
		return
	}

	var admin model.Admin
	if err := database.GetDB().First(&admin, userID).Error; err != nil {
		response.NotFound(c, "管理员不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.OldPassword)); err != nil {
		response.ParamError(c, "旧密码不正确")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "密码加密失败")
		return
	}

	if err := database.GetDB().Model(&admin).Update("password", string(hashedPassword)).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "密码修改成功")
}
