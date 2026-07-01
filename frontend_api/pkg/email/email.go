package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"frontend_api/config"
)

// Sender SMTP 发送器
type Sender struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
	UseTLS   bool // true=SSL(465), false=STARTTLS(587)
}

// GetSender 根据 provider 获取对应的 SMTP 发送器
func GetSender(provider string) (*Sender, error) {
	cfg, ok := config.AppConfig.Email.Providers[provider]
	if !ok {
		return nil, fmt.Errorf("不支持的邮箱服务商: %s", provider)
	}
	return &Sender{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		From:     cfg.From,
		UseTLS:   cfg.Port == "465",
	}, nil
}

// SendResetEmail 发送密码重置邮件
func (s *Sender) SendResetEmail(to, token string) error {
	resetURL := fmt.Sprintf("%s/resetPassword?token=%s", config.AppConfig.Server.SiteURL, token)

	body := fmt.Sprintf(`From: %s
To: %s
Subject: =?UTF-8?B?5a+G56CB6YeN572u6YKu5Lu2?=
MIME-Version: 1.0
Content-Type: text/html; charset=UTF-8

<!DOCTYPE html>
<html>
<body style="background:#0a0a0f;padding:40px 0;margin:0;font-family:Arial,sans-serif">
<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:480px;margin:0 auto">
  <tr><td style="padding:32px;background:#12121a;border-radius:16px;border:1px solid rgba(0,255,255,0.15)">
    <h2 style="color:#00ffff;margin:0 0 16px;font-size:22px">🔐 密码重置</h2>
    <p style="color:#888;font-size:14px;line-height:1.8;margin:0 0 24px">
      我们收到了你的密码重置请求。点击下方按钮即可设置新密码，链接 <b style="color:#ff00ff">30 分钟内</b> 有效。
    </p>
    <a href="%s" style="display:block;text-align:center;padding:14px 0;background:linear-gradient(135deg,rgba(0,255,255,0.2),rgba(255,0,255,0.2));border:1px solid rgba(0,255,255,0.4);border-radius:28px;color:#fff;text-decoration:none;font-size:15px;font-weight:600">重置密码</a>
    <p style="color:#555;font-size:12px;margin:20px 0 0;line-height:1.6">
      如果按钮无法点击，请复制以下链接到浏览器：<br/>
      <span style="color:#00ffff;word-break:break-all">%s</span>
    </p>
  </td></tr>
</table>
</body>
</html>`, s.From, to, resetURL, resetURL)

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.User, s.Password, s.Host)

	if s.UseTLS {
		// SSL 直连（465）
		tlsConfig := &tls.Config{
			ServerName:         s.Host,
			InsecureSkipVerify: false,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("TLS 连接失败: %w", err)
		}
		client, err := smtp.NewClient(conn, s.Host)
		if err != nil {
			return fmt.Errorf("SMTP 客户端创建失败: %w", err)
		}
		defer client.Quit()
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP 认证失败: %w", err)
		}
		if err := client.Mail(s.From); err != nil {
			return fmt.Errorf("MAIL FROM 失败: %w", err)
		}
		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("RCPT TO 失败: %w", err)
		}
		wc, err := client.Data()
		if err != nil {
			return fmt.Errorf("DATA 失败: %w", err)
		}
		_, err = fmt.Fprint(wc, body)
		if err != nil {
			return fmt.Errorf("写入邮件内容失败: %w", err)
		}
		return wc.Close()
	}

	// STARTTLS（587）
	err := smtp.SendMail(addr, auth, s.From, []string{to}, []byte(body))
	if err != nil {
		return fmt.Errorf("邮件发送失败: %w", err)
	}
	return nil
}

// ValidateEmail 校验邮箱格式并识别 provider
func ValidateEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || !strings.Contains(email, "@") {
		return "", fmt.Errorf("邮箱格式不正确")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("邮箱格式不正确")
	}
	domain := parts[1]

	switch {
	case strings.Contains(domain, "qq.com"):
		return "qq", nil
	case strings.Contains(domain, "163.com"):
		return "163", nil
	case strings.Contains(domain, "gmail.com"):
		return "gmail", nil
	default:
		return "", fmt.Errorf("暂不支持该邮箱服务商，请使用 QQ/163/Gmail 邮箱")
	}
}
