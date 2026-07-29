package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var pre = "data:image/jpeg;base64,"

const maxNoCompressSize = 100 * 1024 // 100KB

func Img2Base64(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开图片失败: %w", err)
	}
	defer f.Close()

	// 读取全部内容到内存（一般图片都可接受）
	data, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("读取图片失败: %w", err)
	}

	// 小于 xxKB 直接返回，不压缩
	if len(data) <= maxNoCompressSize {
		return pre + base64.StdEncoding.EncodeToString(data), nil
	}

	// ---------------- 大图开始压缩 ----------------
	ext := strings.ToLower(filepath.Ext(filePath))

	var img image.Image
	var decodeErr error

	switch ext {
	case ".jpg", ".jpeg":
		img, decodeErr = jpeg.Decode(bytes.NewReader(data))
	case ".png":
		img, decodeErr = png.Decode(bytes.NewReader(data))
	default:
		// 不支持的格式 → 直接返回原始（也可以选择返回错误）
		return "", fmt.Errorf("不支持的格式:%s", ext)
	}

	if decodeErr != nil {
		return "", fmt.Errorf("解码图片失败: %w", decodeErr)
	}

	// 准备压缩输出缓冲
	buf := bytes.NewBuffer(nil)

	// 根据文件类型决定用什么编码方式压缩
	switch ext {
	case ".jpg", ".jpeg", ".png":
		// jpeg 压缩，质量 75（可调：60~85 之间）
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 60})

		if err != nil {
			return "", fmt.Errorf("图片编码失败: %w", err)
		}

		compressed := buf.Bytes()

		// 可选：打印压缩前后对比（调试用）
		GetLogger().Debug("压缩成功", zap.Any("压缩前", len(data)), zap.Any("压缩后", len(compressed)))

		return pre + base64.StdEncoding.EncodeToString(compressed), nil
	default:
		// 正常来说到不了这里
		return "", fmt.Errorf("逻辑错误，请检查代码逻辑，文件后缀:%s", ext)
	}
}

func Byte2Base64(imgBytes []byte) (string, error) {
	if imgBytes == nil {
		return "", errors.New("传入图片为空")
	}
	base64Str := base64.StdEncoding.EncodeToString(imgBytes)
	return base64Str, nil
}
