package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"github.com/53AI/53AIHub/common"
	"github.com/53AI/53AIHub/config"
	"github.com/53AI/53AIHub/model"
	"github.com/53AI/53AIHub/service"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	"gorm.io/gorm"
)

const (
	verificationCodeLength = 6
	duration               = 15
	codeExpiration         = duration * time.Minute
)

type SendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@126.com"` // 需要验证的邮箱地址
}

// @Summary 发送邮箱验证码
// @Description 向指定邮箱发送6位数字验证码（有效期10分钟）
// @Tags Email
// @Accept json
// @Produce json
// @Param data body SendVerificationEmailRequest true "邮箱验证请求"
// @Success 200 {object} model.CommonResponse{data=string} "成功响应：验证码已发送"
// @Failure 400 {object} model.CommonResponse "参数错误"
// @Failure 500 {object} model.CommonResponse "服务器内部错误"
// @Router /api/email/send_verification [post]
func SendVerificationEmail(c *gin.Context) {
	var req SendVerificationEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	code, err := common.GenerateRandomCode(verificationCodeLength)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	now := time.Now().UTC().UnixMilli()
	currentDayStart := now - (now % 86400000) // 计算当日0点时间戳（毫秒）

	var existingVC model.VerificationCode
	err = model.DB.Where("target = ? AND type = 'email' AND created_time >= ?", req.Email, currentDayStart).First(&existingVC).Error

	if err == nil {
		existingVC.Code = code
		existingVC.ExpiresAt = now + int64(codeExpiration.Milliseconds())
		existingVC.DailyCount++
		err = model.DB.Save(&existingVC).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		newVC := model.VerificationCode{
			Type:       model.VerificationCodeTypeEmail,
			Target:     req.Email,
			Code:       code,
			ExpiresAt:  now + int64(codeExpiration.Milliseconds()),
			DailyCount: 1,
		}
		err = model.DB.Create(&newVC).Error
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	_ = common.RedisSet("email_verification:"+req.Email, code, codeExpiration)

	eid := config.GetEID(c)
	e := email.NewEmail()
	e.To = []string{req.Email}
	e.Subject = "邮箱验证码"
	e.Text = []byte(fmt.Sprintf("您的验证码是：%s，有效期%d分钟", code, duration))

	auth, from, host, port, isSsl, err := service.GetSmtpConfig(eid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NetworkError.ToResponse(fmt.Errorf("failed to get SMTP auth: %w", err)))
		return
	}
	if from == "" {
		c.JSON(http.StatusInternalServerError, model.NetworkError.ToResponse(errors.New("SMTP from address is empty")))
		return
	}
	e.From = from

	if err := common.SendEmail(e, auth, isSsl, host, port); err != nil {
		c.JSON(http.StatusInternalServerError, model.NetworkError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse("Verification code has been sent"))
}

type SendTestEmailRequest struct {
	Host     string `json:"host" binding:"required" example:"smtp.126.com"`
	Port     int    `json:"port" binding:"required" example:"465"`
	Username string `json:"username" binding:"required" example:"user@126.com"`
	Password string `json:"password" binding:"required" example:"123456"`
	From     string `json:"from" binding:"required" example:"user@126.com"`
	IsSSL    bool   `json:"is_ssl" example:"true"`
	To       string `json:"to" binding:"required" example:"user@126.com"`
}

// @Summary 发送邮箱测试邮件
// @Description 向指定邮箱发送测试信息
// @Tags Email
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body SendTestEmailRequest true "邮箱测试请求"
// @Success 200 {object} model.CommonResponse{data=string} "成功响应：已发送"
// @Failure 400 {object} model.CommonResponse "参数错误"
// @Failure 500 {object} model.CommonResponse "服务器内部错误"
// @Router /api/email/send_test [post]
func SendTestEmail(c *gin.Context) {
	var req SendTestEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	e := email.NewEmail()
	e.To = []string{req.To}
	e.Subject = "53AI Hub SMTP设置测试邮件！"
	e.Text = []byte("收到此邮件表示配置无误")

	from := req.From
	auth := smtp.PlainAuth(
		"",
		req.Username,
		req.Password,
		req.Host,
	)

	// 使用配置中的值
	host := req.Host
	port := req.Port
	isSsl := req.IsSSL
	e.From = from

	if err := common.SendEmail(e, auth, isSsl, host, port); err != nil {
		c.JSON(http.StatusInternalServerError, model.NetworkError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse("Test email has been sent"))
}

// UpdateUserEmailRequest 更新用户邮箱请求结构体
type UpdateUserEmailRequest struct {
	Email string `json:"email" binding:"required,email"`      // 新邮箱地址
	Code  string `json:"code" binding:"required,min=6,max=6"` // 验证码（6位数字）
}

// UpdateUserEmail 绑定、更新用户邮箱
// @Summary 绑定、更新用户邮箱
// @Description 通过验证码验证后更新用户绑定的邮箱
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param data body UpdateUserEmailRequest true "新邮箱及验证码"
// @Success 200 {object} model.CommonResponse{data=model.User} "更新成功"
// @Failure 400 {object} model.CommonResponse "参数错误"
// @Failure 401 {object} model.CommonResponse "验证码无效"
// @Failure 409 {object} model.CommonResponse "邮箱已被绑定"
// @Router /api/users/{id}/email [patch]
func UpdateUserEmail(c *gin.Context) {
	// 解析路径参数ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	var req UpdateUserEmailRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	// 验证验证码
	valid, err := common.VerifyEmailCode(req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.AuthFailed.ToResponse(err))
		return
	}
	if !valid {
		c.JSON(http.StatusUnauthorized, model.AuthFailed.ToNewErrorResponse(model.InvalidVerificationCode))
		return
	}

	// 获取当前用户
	user, err := model.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.NotFound.ToResponse(err))
		return
	}

	// 检查邮箱是否已被其他用户绑定
	existingUser, err := model.GetUserByEmail(user.Eid, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}
	if existingUser.UserID > 0 && existingUser.UserID != id {
		err := errors.New("This email has been bound by another user")
		c.JSON(http.StatusConflict, model.AuthFailed.ToErrorResponse(err))
		return
	}

	// 更新邮箱
	user.Email = req.Email
	if err := user.Update(false); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(user))
}
