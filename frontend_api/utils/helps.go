package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateCaptcha 生成验证码，返回 captcha_id 和 4位随机数字
func GenerateCaptcha() map[string]string {
	// 唯一 ID
	b := make([]byte, 16)
	rand.Read(b)

	// 4位随机数字
	const digits = "0123456789"
	code := make([]byte, 4)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		code[i] = digits[n.Int64()]
	}

	return map[string]string{
		"captcha_id": fmt.Sprintf("%x", b),
		"captcha":    string(code),
	}
}
