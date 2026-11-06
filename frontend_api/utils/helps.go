package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"math"
	"runtime"
	"time"
)

// ═══════════════════════════════════════════════
//  公共工具方法
//  与业务无关、可被多个包引用的通用能力放这里
// ═══════════════════════════════════════════════

// 图片压缩参数：长边上限与 JPEG 质量
const (
	// MaxImageEdge 图片长边上限（超过则等比缩小）
	MaxImageEdge = 1600
	// jpegQuality JPEG 编码质量
	jpegQuality = 85
)

// white 白色填充色（PNG 透明区域转 JPEG 前铺垫，否则透明会变黑）
var white = color.RGBA{R: 255, G: 255, B: 255, A: 255}

// HeapMB 返回当前 Go 堆占用（MB），用于任务级内存观测与扩容评估
func HeapMB() float64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return float64(ms.HeapAlloc) / (1024 * 1024)
}

// RetryDelay 指数退避等待时间：1s / 2s / 4s / …，上限 8s（attempt 从 1 开始）
func RetryDelay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	d := time.Duration(1<<uint(attempt-1)) * time.Second
	if d > 8*time.Second {
		d = 8 * time.Second
	}
	return d
}

// CompressImageToBase64 解码图片、等比缩放到长边 ≤ MaxImageEdge 并编码为 JPEG base64
// 返回 data URL 格式：data:image/jpeg;base64,...
// 说明：设计稿原图直接 base64 会让内存峰值达到几十 MB，压缩后通常只剩几百 KB
func CompressImageToBase64(raw []byte) (string, error) {
	src, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	img := ResizeToMaxEdge(src, MaxImageEdge)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return "", err
	}
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// ResizeToMaxEdge 等比缩放到长边不超过 maxEdge（双线性插值，透明区域铺白底）；小图不放大
func ResizeToMaxEdge(src image.Image, maxEdge int) image.Image {
	flat := flattenOnWhite(src)
	b := flat.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= maxEdge && h <= maxEdge {
		return flat
	}

	scale := float64(maxEdge) / float64(w)
	if h > w {
		scale = float64(maxEdge) / float64(h)
	}
	nw, nh := int(float64(w)*scale+0.5), int(float64(h)*scale+0.5)
	if nw < 1 {
		nw = 1
	}
	if nh < 1 {
		nh = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
	xRatio := float64(w) / float64(nw)
	yRatio := float64(h) / float64(nh)
	for y := 0; y < nh; y++ {
		sy := (float64(y)+0.5)*yRatio - 0.5
		y0 := int(math.Floor(sy))
		fy := sy - float64(y0)
		y1 := y0 + 1
		if y0 < 0 {
			y0 = 0
		}
		if y1 >= h {
			y1 = h - 1
		}
		row := dst.Pix[y*dst.Stride : y*dst.Stride+nw*4]
		for x := 0; x < nw; x++ {
			sx := (float64(x)+0.5)*xRatio - 0.5
			x0 := int(math.Floor(sx))
			fx := sx - float64(x0)
			x1 := x0 + 1
			if x0 < 0 {
				x0 = 0
			}
			if x1 >= w {
				x1 = w - 1
			}

			c00 := flat.RGBAAt(x0, y0)
			c10 := flat.RGBAAt(x1, y0)
			c01 := flat.RGBAAt(x0, y1)
			c11 := flat.RGBAAt(x1, y1)

			w00 := (1 - fx) * (1 - fy)
			w10 := fx * (1 - fy)
			w01 := (1 - fx) * fy
			w11 := fx * fy

			i := x * 4
			row[i] = uint8(float64(c00.R)*w00 + float64(c10.R)*w10 + float64(c01.R)*w01 + float64(c11.R)*w11)
			row[i+1] = uint8(float64(c00.G)*w00 + float64(c10.G)*w10 + float64(c01.G)*w01 + float64(c11.G)*w11)
			row[i+2] = uint8(float64(c00.B)*w00 + float64(c10.B)*w10 + float64(c01.B)*w01 + float64(c11.B)*w11)
			row[i+3] = 255
		}
	}
	return dst
}

// flattenOnWhite 将图片绘制到白色背景的 RGBA 画布（处理 PNG 透明区域）
func flattenOnWhite(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(white), image.Point{}, draw.Src)
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Over)
	return dst
}
