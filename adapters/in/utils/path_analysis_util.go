package utils

import (
	"fmt"
	"os"
)

type PathType int

const (
	SingleGameFolder PathType = iota
	ManyGameFolder
	ZipFile
	ErrType PathType = -1
)

// AnalysisPath 判断文件路径的类型：
func AnalysisPath(filePath string) (PathType, error) {
	// 获取文件信息
	info, err := os.Stat(filePath)
	if err != nil {
		return ErrType, fmt.Errorf("获取文件信息失败：%e", err)
	}

	// 如果是文件（不是目录），视为 ZipFile
	if !info.IsDir() {
		return ZipFile, nil
	}

	// 是目录，打开并读取其内容
	f, err := os.Open(filePath)
	if err != nil {
		return ErrType, fmt.Errorf("打开文件夹失败：%e", err)
	}
	defer f.Close()

	// 读取所有目录条目（包括文件和子目录）
	entries, err := f.ReadDir(-1) // -1 表示读取全部
	if err != nil {
		return ErrType, fmt.Errorf("读取所有目录失败：%e", err)
	}
	if len(entries) == 0 {
		return ErrType, fmt.Errorf("传入的是个空目录")
	}

	hasFile := false

	// 遍历条目，统计文件
	for _, entry := range entries {
		if !entry.IsDir() {
			hasFile = true
		}
	}

	// 根据统计结果返回相应类型
	if hasFile {
		return SingleGameFolder, nil
	} else {
		return ManyGameFolder, nil
	}
}
