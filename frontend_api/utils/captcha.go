package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// captchaStore 内存存储: captcha_id → answer (5分钟过期)
var (
	captchaStore = sync.Map{}
)

type captchaEntry struct {
	Answer    string
	ExpiresAt time.Time
}

// GenerateCaptcha 生成 SVG 验证码图片
// 返回 captcha_id、base64 图片、正确答案
func GenerateCaptcha() (captchaID string, imageBase64 string, answer string) {
	// 4位随机数字
	const digits = "0123456789"
	code := make([]byte, 4)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		code[i] = digits[n.Int64()]
	}
	answer = string(code)

	// 唯一 ID
	b := make([]byte, 12)
	rand.Read(b)
	captchaID = fmt.Sprintf("%x", b)

	// 存入内存，5分钟过期
	captchaStore.Store(captchaID, captchaEntry{
		Answer:    answer,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	// 生成 SVG（扭曲数字 + 噪点 + 干扰线）
	svg := buildCaptchaSVG(code)
	imageBase64 = "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg))

	return
}

// VerifyCaptcha 校验验证码（一次性消费，验证后删除）
func VerifyCaptcha(id, input string) bool {
	val, ok := captchaStore.Load(id)
	if !ok {
		return false
	}
	entry := val.(captchaEntry)
	if time.Now().After(entry.ExpiresAt) {
		captchaStore.Delete(id)
		return false
	}
	// 验证后删除，防止重复使用
	captchaStore.Delete(id)
	return entry.Answer == input
}

// buildCaptchaSVG 构建带干扰的 SVG 验证码图片
func buildCaptchaSVG(code []byte) string {
	const W, H = 140, 50

	// 背景
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">
<rect width="%d" height="%d" fill="#1a1a2e" rx="6"/>`, W, H, W, H)

	// 干扰线
	nLines, _ := rand.Int(rand.Reader, big.NewInt(3))
	for i := int64(0); i < nLines.Int64()+2; i++ {
		x1, _ := rand.Int(rand.Reader, big.NewInt(int64(W)))
		y1, _ := rand.Int(rand.Reader, big.NewInt(int64(H)))
		x2, _ := rand.Int(rand.Reader, big.NewInt(int64(W)))
		y2, _ := rand.Int(rand.Reader, big.NewInt(int64(H)))
		cIdx, _ := rand.Int(rand.Reader, big.NewInt(4))
		colors := []string{"#00ffff44", "#ff00ff44", "#ffffff22", "#00ff8844"}
		svg += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1"/>`, x1, y1, x2, y2, colors[cIdx.Int64()])
	}

	// 噪点
	nDots, _ := rand.Int(rand.Reader, big.NewInt(30))
	for i := int64(0); i < nDots.Int64()+15; i++ {
		cx, _ := rand.Int(rand.Reader, big.NewInt(int64(W)))
		cy, _ := rand.Int(rand.Reader, big.NewInt(int64(H)))
		r, _ := rand.Int(rand.Reader, big.NewInt(3))
		dotColors := []string{"#00ffff", "#ff00ff", "#ffffff", "#00ff88"}
		cIdx, _ := rand.Int(rand.Reader, big.NewInt(4))
		svg += fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d" fill="%s" opacity="0.6"/>`, cx, cy, r.Int64()+1, dotColors[cIdx.Int64()])
	}

	// 数字（不同旋转、颜色、偏移）
	fontColors := []string{"#00ffff", "#ff6b9d", "#00ff88", "#ffaa00"}
	for i, ch := range code {
		rot, _ := rand.Int(rand.Reader, big.NewInt(40))
		rotDeg := rot.Int64() - 20 // -20 ~ +19 度
		x := 18 + i*30
		y, _ := rand.Int(rand.Reader, big.NewInt(14))
		yOff := y.Int64() + 28 // 28~41
		size, _ := rand.Int(rand.Reader, big.NewInt(8))
		fontSize := size.Int64() + 24 // 24~31
		cIdx, _ := rand.Int(rand.Reader, big.NewInt(4))
		svg += fmt.Sprintf(`<text x="%d" y="%d" transform="rotate(%d, %d, %d)" font-size="%d" fill="%s" font-family="Arial, sans-serif" font-weight="bold" text-anchor="middle">%c</text>`,
			x, yOff, rotDeg, x, yOff, fontSize, fontColors[cIdx.Int64()], ch)
	}

	svg += `</svg>`
	return svg
}
