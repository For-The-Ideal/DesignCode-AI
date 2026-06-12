package cos

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"frontend_api/config"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// Client COS 客户端单例
var client *cos.Client

// InitClient 初始化 COS 客户端
func InitClient() {
	cfg := config.AppConfig.COS
	if cfg.Bucket == "" || cfg.Region == "" || cfg.SecretID == "" || cfg.SecretKey == "" {
		log.Println("[COS] 配置不完整，跳过初始化")
		return
	}

	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.Bucket, cfg.Region))
	b := &cos.BaseURL{BucketURL: u}
	client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})
	log.Println("[COS] 客户端初始化成功")
}

// Upload 上传文件到 COS，返回可公开访问的 URL
func Upload(file *multipart.FileHeader) (string, error) {
	if client == nil {
		return "", fmt.Errorf("COS 客户端未初始化")
	}

	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	// 生成存储路径：uploads/年/月/日/uuid_原始文件名
	ext := filepath.Ext(file.Filename)
	now := time.Now()
	objectKey := fmt.Sprintf("uploads/%d/%02d/%02d/%d%s",
		now.Year(), now.Month(), now.Day(),
		now.UnixNano(), ext)

	_, err = client.Object.Put(context.Background(), objectKey, f, nil)
	if err != nil {
		return "", fmt.Errorf("上传到 COS 失败: %w", err)
	}

	// 返回 URL
	cfg := config.AppConfig.COS
	publicURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", cfg.Bucket, cfg.Region, objectKey)
	log.Printf("[COS] 上传成功: %s (%d bytes)", publicURL, file.Size)
	return publicURL, nil
}

// UploadBytes 上传字节数据到 COS
func UploadBytes(data []byte, filename string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("COS 客户端未初始化")
	}

	ext := filepath.Ext(filename)
	now := time.Now()
	objectKey := fmt.Sprintf("uploads/%d/%02d/%02d/%d%s",
		now.Year(), now.Month(), now.Day(),
		now.UnixNano(), ext)

	_, err := client.Object.Put(context.Background(), objectKey, bytes.NewReader(data), nil)
	if err != nil {
		return "", fmt.Errorf("上传到 COS 失败: %w", err)
	}

	cfg := config.AppConfig.COS
	publicURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", cfg.Bucket, cfg.Region, objectKey)
	log.Printf("[COS] 上传成功: %s (%d bytes)", publicURL, len(data))
	return publicURL, nil
}

// DecodeBase64Image 解码 base64 图片数据（去掉 data:image/xxx;base64, 前缀）
func DecodeBase64Image(base64Str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64Str)
}
