package logger

import (
	"gin-init/config" // 导入配置包
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
)

// 线程安全的 Writer 包装器
type syncWriter struct {
	io.Writer
	mu sync.Mutex
}

func (w *syncWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Writer.Write(p)
}

// Init 初始化日志切割（在 Gin 引擎创建前调用）
func Init() {
	// 在 logger.Init() 中添加
	if err := os.MkdirAll(filepath.Dir(config.Cfg.Log.Filename), 0755); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
	}

	// 根据配置初始化 lumberjack 切割器
	logWriter := &lumberjack.Logger{
		Filename:   config.Cfg.Log.Filename,   // 日志文件路径
		MaxSize:    config.Cfg.Log.MaxSize,    // 单个文件最大 MB
		MaxAge:     config.Cfg.Log.MaxAge,     // 保留天数
		MaxBackups: config.Cfg.Log.MaxBackups, // 最大备份数
		Compress:   config.Cfg.Log.Compress,   // 压缩旧日志
		LocalTime:  true,                      // 日志文件名使用本地时间
	}

	// 2. 重定向 stdout/stderr：终端 + 切割文件
	mwStdout := io.MultiWriter(os.Stdout, logWriter)
	mwStderr := io.MultiWriter(os.Stderr, logWriter)

	// 使用 syncWriter 包装 MultiWriter
	syncStdout := &syncWriter{Writer: mwStdout}
	syncStderr := &syncWriter{Writer: mwStderr}

	// 重定向标准输出和错误输出
	// 注意：不能直接赋值给 os.Stdout/os.Stderr，应该使用 os.RedirectStdout/os.RedirectStderr 或者
	// 直接设置 log 包的输出目标
	log.SetOutput(syncStdout)

	// 如果需要捕获所有程序的标准输出和错误输出，需要在程序开始时就设置
	// 这里我们主要确保日志能正确写入文件
	// 替换 Gin 的默认 writer（同时输出到控制台和切割文件，可选）
	gin.DefaultWriter = syncStdout
	gin.DefaultErrorWriter = syncStderr

	// 如需自定义日志格式（如添加时间、级别），可进一步封装
	// 例如结合 zap、logrus 等日志库，此处以基础用法为例
}
