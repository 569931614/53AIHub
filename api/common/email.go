package common

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/53AI/53AIHub/model"
	"github.com/jordan-wright/email"
	"gorm.io/gorm"
)

// ValidateEmailFormat 验证基础邮箱格式（通用版）
func ValidateEmailFormat(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$` // 通用邮箱格式验证（如user@domain.com）
	return regexp.MustCompile(pattern).MatchString(email)
}

// sendQQMailWithTLS 使用手动TLS连接发送QQ邮件
func sendQQMailWithTLS(e *email.Email, auth smtp.Auth, host string, port int) error {
	// 创建 TLS 配置
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	// 连接到 SMTP 服务器
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return fmt.Errorf("TLS connection failed: %v", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("SMTP client creation failed: %v", err)
	}
	defer client.Quit()

	// 认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// 设置发件人和收件人
	if err = client.Mail(e.From); err != nil {
		return fmt.Errorf("sender setup failed: %v", err)
	}

	for _, recipient := range e.To {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("recipient setup failed: %v", err)
		}
	}

	// 写入邮件内容
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("data writer creation failed: %v", err)
	}
	defer wc.Close()

	// 构建邮件内容
	var msg strings.Builder
	// 添加From
	msg.WriteString(fmt.Sprintf("From: %s\r\n", e.From))
	// 添加To
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(e.To, ", ")))
	// 添加Subject
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", e.Subject))

	// 添加其他头部信息
	for key, values := range e.Headers {
		for _, value := range values {
			msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
	}

	// 添加空行分隔头部和正文
	msg.WriteString("\r\n")

	// 添加正文
	if len(e.HTML) > 0 {
		msg.Write(e.HTML)
	} else if len(e.Text) > 0 {
		msg.Write(e.Text)
	}

	// 发送邮件内容
	_, err = wc.Write([]byte(msg.String()))
	if err != nil {
		return fmt.Errorf("message sending failed: %v", err)
	}

	return nil
}

// SendEmail 使用jordan-wright/email库通过通用SMTP服务发送邮件
func SendEmail(e *email.Email, auth smtp.Auth, isSsl bool, host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)

	fmt.Printf("Sending email to %v via %s, SSL: %v, Port: %d\n", e.To, addr, isSsl, port)

	// 针对QQ邮箱等特殊处理
	isQQMail := strings.Contains(host, "qq.com")

	// 对于QQ邮箱且使用465端口，使用特殊处理方式
	if isQQMail && port == 465 {
		fmt.Printf("Using special handling for QQ mail on port 465\n")
		return sendQQMailWithTLS(e, auth, host, port)
	}

	if isSsl {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         host,
		}

		if err := e.SendWithTLS(addr, auth, tlsConfig); err != nil {
			return fmt.Errorf("failed to send email via %s:%d: %w", host, port, err)
		}
	} else {
		if err := e.Send(addr, auth); err != nil {
			return fmt.Errorf("failed to send email via %s:%d: %w", host, port, err)
		}
	}
	return nil
}

// VerifyEmailCode 验证邮箱验证码有效性
func VerifyEmailCode(email, code string) (bool, error) {
	if email == "" || code == "" {
		return false, errors.New("missing email or code parameter")
	}

	now := time.Now().UTC().UnixMilli()
	var storedCode string
	// 优先从Redis获取验证码
	storedCode, err := RedisGet("email_verification:" + email)
	if err != nil || storedCode == "" {
		// Redis不存在时查询数据库
		var vc model.VerificationCode
		err = model.DB.Where("target = ? AND type = 'email' AND code = ? AND expires_at > ?", email, code, now).First(&vc).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, errors.New("verification code expired or invalid")
			}
			return false, fmt.Errorf("database query error: %w", err)
		}
		storedCode = code
	}

	if storedCode != code {
		return false, errors.New(model.InvalidVerificationCode)
	}

	_ = RedisDel("email_verification:" + email)
	return true, nil
}

func GenerateRandomCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		bytes[i] = byte(48 + int(bytes[i])%10) // 0-9
	}
	return string(bytes), nil
}
