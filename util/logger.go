// Package util 提供带级别的、支持轮转的文件日志功能，并同步输出到控制台
package util

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
)

const (
	LevelError         = iota // 0
	LevelWarning              // 1
	LevelInformational        // 2
	LevelDebug                // 3
)

var (
	logger *Logger
	once   sync.Once
)

// Logger 是日志结构体
type Logger struct {
	level   int
	writers map[int]io.Writer // 每个级别对应一个 writer（含控制台 + 文件）
	closers []io.Closer       // 用于 Close 所有 lumberjack 实例
}

// newLumberjackLogger 创建一个 lumberjack 日志轮转器
func newLumberjackLogger(filename string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100, // MB
		MaxAge:     7,   // 天
		MaxBackups: 10,
		LocalTime:  true,
		Compress:   false,
	}
}

// BuildLogger 初始化全局日志器
// level: "error", "warning", "info", "debug"
func BuildLogger(level string) {
	// 确保 logs 目录存在（不存在则创建，存在则跳过）
	fmt.Println("正在初始化日志系统...")

	_ = os.MkdirAll("logs", 0755)

	intLevel := LevelInformational
	switch level {
	case "error":
		intLevel = LevelError
	case "warning":
		intLevel = LevelWarning
	case "debug":
		intLevel = LevelDebug
	case "info":
		intLevel = LevelInformational
	default:
		intLevel = LevelInformational
	}

	// 为每个级别创建对应的文件 logger
	errorLog := newLumberjackLogger("logs/error.log")
	warnLog := newLumberjackLogger("logs/warning.log")
	infoLog := newLumberjackLogger("logs/info.log")
	debugLog := newLumberjackLogger("logs/debug.log")

	// 所有日志都输出到控制台 + 对应级别文件
	errorWriter := io.MultiWriter(os.Stdout, errorLog)
	warnWriter := io.MultiWriter(os.Stdout, warnLog)
	infoWriter := io.MultiWriter(os.Stdout, infoLog)
	debugWriter := io.MultiWriter(os.Stdout, debugLog)

	logger = &Logger{
		level: intLevel,
		writers: map[int]io.Writer{
			LevelError:         errorWriter,
			LevelWarning:       warnWriter,
			LevelInformational: infoWriter,
			LevelDebug:         debugWriter,
		},
		closers: []io.Closer{errorLog, warnLog, infoLog, debugLog},
	}
}

// writeLog 内部统一写入方法
func (ll *Logger) writeLog(level int, prefix, format string, v ...interface{}) {
	if level > ll.level {
		return
	}
	writer, ok := ll.writers[level]
	if !ok || writer == nil {
		return
	}
	msg := fmt.Sprintf("%s %s %s\n",
		time.Now().Format(time.RFC3339),
		prefix,
		fmt.Sprintf(format, v...))
	fmt.Fprint(writer, msg)
}

// Panic 记录致命错误并退出程序
func (ll *Logger) Panic(format string, v ...interface{}) {
	ll.writeLog(LevelError, "[Panic]", format, v...)
	formatted := fmt.Sprintf(format, v...)
	panic(formatted)
}

// Error 记录错误日志
func (ll *Logger) Error(format string, v ...interface{}) {
	ll.writeLog(LevelError, "[E]", format, v...)
}

// Warning 记录警告日志
func (ll *Logger) Warning(format string, v ...interface{}) {
	ll.writeLog(LevelWarning, "[W]", format, v...)
}

// Info 记录信息日志
func (ll *Logger) Info(format string, v ...interface{}) {
	ll.writeLog(LevelInformational, "[I]", format, v...)
}

// Debug 记录调试日志
func (ll *Logger) Debug(format string, v ...interface{}) {
	ll.writeLog(LevelDebug, "[D]", format, v...)
}

// Log 返回全局日志实例（线程安全）
func Log() *Logger {
	once.Do(func() {
		if logger == nil {
			BuildLogger("debug")
		}
	})
	return logger
}

// CloseLogger 关闭所有日志文件（触发 flush）
func CloseLogger() error {
	if logger != nil {
		for _, closer := range logger.closers {
			if err := closer.Close(); err != nil {
				return err // 或者收集所有错误
			}
		}
	}
	return nil
}
