package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"frontend_api/config"
)

// ImageHTTPClient 下载生成任务图片的专用客户端
//
//   - 60s 超时：避免恶意/慢速地址挂死 worker
//   - 重定向逐跳校验：防止 302 跳到内网地址绕过域名白名单
var ImageHTTPClient = &http.Client{
	Timeout: 60 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 3 {
			return errors.New("重定向次数过多")
		}
		return ValidateImageURL(req.URL.String())
	},
}

// ValidateImageURL 校验图片下载地址，防 SSRF：
//   - 必须是 https
//   - 域名必须命中白名单：COS 桶域名（自动推导）或 server.allowed_image_hosts
//     （支持 *.example.com 通配子域）
func ValidateImageURL(raw string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Hostname() == "" {
		return errors.New("URL 格式不正确")
	}
	if u.Scheme != "https" {
		return errors.New("仅支持 https 图片地址")
	}

	host := strings.ToLower(u.Hostname())
	for _, rule := range allowedImageHosts() {
		if matchHost(host, rule) {
			return nil
		}
	}
	return fmt.Errorf("图片域名 %s 不在允许列表", host)
}

// allowedImageHosts 汇总允许的图片域名：COS 桶域名 + 配置额外放行项
func allowedImageHosts() []string {
	var hosts []string
	if config.AppConfig == nil {
		return hosts
	}

	if cos := config.AppConfig.COS; cos.Bucket != "" && cos.Region != "" {
		hosts = append(hosts, fmt.Sprintf("%s.cos.%s.myqcloud.com",
			strings.ToLower(cos.Bucket), strings.ToLower(cos.Region)))
	}
	for _, h := range config.AppConfig.Server.AllowedImageHosts {
		if h = strings.ToLower(strings.TrimSpace(h)); h != "" {
			hosts = append(hosts, h)
		}
	}
	return hosts
}

// matchHost 域名匹配：精确相等，或通配规则 "*.example.com" 匹配其子域
func matchHost(host, rule string) bool {
	if host == rule {
		return true
	}
	if strings.HasPrefix(rule, "*.") {
		return strings.HasSuffix(host, rule[1:]) // rule[1:] = ".example.com"
	}
	return false
}
