package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var once = sync.Once{}
var logInstance *zap.Logger

func initLog() {
	// 1. 控制台输出配置
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 2. 文件输出配置（使用自定义日期滚动 Writer）
	fileWriter := newDateWriter("./logs") // 日志存放到 ./logs 目录
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()) // 可改为 ConsoleEncoder
	fileSyncer := zapcore.AddSync(fileWriter) // fileWriter 已实现 Sync

	// 3. 设置日志级别（可改为 zapcore.InfoLevel）
	level := zapcore.DebugLevel

	// 4. 创建两个 Core
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, level)
	fileCore := zapcore.NewCore(fileEncoder, fileSyncer, level)

	// 5. 组合成 Tee Core，同时写入两处
	core := zapcore.NewTee(consoleCore, fileCore)

	// 6. 构建 Logger，附加调用者信息
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	logInstance = logger
}

// GetLogger 获取全局唯一日志实例（外部仅能通过此方法获取）
func GetLogger() *zap.Logger {
	once.Do(initLog)
	return logInstance
}

// SyncLogger 可选：提供Sync方法，供程序退出时刷新日志缓冲区
func SyncLogger() error {
	return logInstance.Sync()
}

// dateWriter 实现 io.Writer 和 Sync，按日期滚动文件
type dateWriter struct {
	mu      sync.Mutex
	dir     string
	file    *os.File
	dateStr string
}

// newDateWriter 创建日期滚动 Writer，dir 为日志目录
func newDateWriter(dir string) *dateWriter {
	return &dateWriter{dir: dir}
}

// Write 每次写入时检查日期，自动切换文件
func (w *dateWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	dateStr := now.Format("2006-01-02")

	if w.dateStr != dateStr {
		// 日期变化，切换文件
		if w.file != nil {
			w.file.Close()
		}
		// 确保目录存在
		if err := os.MkdirAll(w.dir, 0755); err != nil {
			return 0, err
		}
		filename := filepath.Join(w.dir, dateStr+".log")
		// 以追加模式打开，支持多次启动写入同一文件
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}
		w.file = f
		w.dateStr = dateStr
	}
	return w.file.Write(p)
}

// Sync 刷新文件缓冲区
func (w *dateWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}
