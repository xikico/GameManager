package utils

import (
	"GameManager/domain/entity"
	"GameManager/domain/ports/out/utils"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.design/x/clipboard"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type Utils struct {
}

func GetUtils() utils.Utils {
	return &Utils{}
}

func (u *Utils) OpenGame(game *entity.Game) error {
	ext := strings.ToLower(filepath.Ext(game.StartPath))
	switch ext {
	case ".exe":
		return runAsAdmin(game.StartPath, game.Path)
	case ".bat":
		return errors.New("暂未实现 bat 启动")
	default:
		return fmt.Errorf("不支持的启动文件类型: %s", ext)
	}
}

func (u *Utils) DeleteGame(game *entity.Game) error {
	return os.RemoveAll(game.Path)
}

func (u *Utils) Img2Base64(ctx context.Context) {
	// 1. 初始化剪贴板
	err := clipboard.Init()
	if err != nil {
		return
	}

	timer := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-timer:
			// 2. 尝试读取剪贴板中的图片数据（PNG格式）
			//    注意：Read(FmtImage) 总是返回 PNG 编码的数据[reference:2]
			imgData := clipboard.Read(clipboard.FmtImage)
			if imgData == nil {
				continue
			}
			var compressed []byte = nil
			if len(imgData) > 100*1024 {
				img, decodeErr := png.Decode(bytes.NewReader(imgData))
				if decodeErr == nil {
					buf := bytes.NewBuffer(nil)
					jpeg.Encode(buf, img, &jpeg.Options{Quality: 60})
					compressed = buf.Bytes()
				}
			}
			if compressed != nil {
				imgData = compressed
			}
			base64Str := base64.StdEncoding.EncodeToString(imgData)
			dataURL := fmt.Sprintf("data:image/jpeg;base64,%s", base64Str)
			// 写入剪贴板
			clipboard.Write(clipboard.FmtText, []byte(dataURL))
		case <-ctx.Done():
			return
		}
	}
}

// runAsAdmin 在gamePath目录下以管理员模式打开exePath
func runAsAdmin(exePath string, gamePath string) error {
	shell32 := syscall.NewLazyDLL("shell32.dll")
	shellExecute := shell32.NewProc("ShellExecuteW")

	// 参数说明：
	// hwnd = 0 (no parent window)
	// verb = "runas" (request admin)
	// file = exe path
	// params = nil (no args)
	// dir = nil (use default)
	// showCmd = SW_SHOW (1)
	runas, _ := syscall.UTF16PtrFromString("runas")
	runPath, _ := syscall.UTF16PtrFromString(exePath)
	dir, _ := syscall.UTF16PtrFromString(gamePath)

	ret, _, err := shellExecute.Call(
		0,                                // hwnd
		uintptr(unsafe.Pointer(runas)),   // lpVerb ("runas")
		uintptr(unsafe.Pointer(runPath)), // lpFile (可执行文件路径)
		0,                                // lpParameters (命令行参数，可为 nil)
		uintptr(unsafe.Pointer(dir)),     // lpDirectory 设置工作目录
		1,                                // nShowCmd (SW_SHOWNORMAL)
	)
	if ret <= 32 {
		return err
	}
	return nil
}
