package unzip

import (
	port "GameManager/domain/ports/out/unzip"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Unzip 使用 7z 命令行解压，支持分卷、密码（有/无）。
// path: 压缩文件路径（支持 .rar/.7z/.zip 等，第一个分卷即可）
// password: 密码，空字符串表示无密码
// 返回: 解压到的目录路径, error
func Unzip(path string, password string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("文件路径出错: %w", err)
	}

	// 2. 计算输出目录 = 输入文件的所在目录
	outputDir := filepath.Join(filepath.Dir(path), strings.Split(filepath.Base(path), ".")[0])

	// 3. 构建命令
	args := []string{"x", path, "-o" + outputDir, "-y"} // -y 自动覆盖

	// 密码处理（关键：必须紧跟 -p 无空格）
	if password != "" {
		args = append(args, "-p"+password)
	} else {
		// 即使密码为空，也传递 -p（空），这样如果文件需要密码就会直接报错而不是提示
		args = append(args, "-p")
	}

	// 4. 执行命令（非交互）
	cmd := exec.Command("7z", args...)

	// 捕获输出，用于错误判断
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := stdout.String() + stderr.String()

	// 5. 错误处理
	if err != nil {
		// 常见密码错误
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 2 {
			return "", port.ErrPassword
		}
		return "", fmt.Errorf("7z 执行失败: %v\n输出: %s", err, output)
	}

	// 成功
	return outputDir, nil
}
