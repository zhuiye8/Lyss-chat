package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel 表示日志级别
type LogLevel int

const (
	// DEBUG 调试级别
	DEBUG LogLevel = iota
	// INFO 信息级别
	INFO
	// WARN 警告级别
	WARN
	// ERROR 错误级别
	ERROR
	// FATAL 致命错误级别
	FATAL
)

// Logger 表示日志记录器
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

// New 创建一个新的日志记录器
func New(level string) *Logger {
	return &Logger{
		level:  parseLevel(level),
		logger: log.New(os.Stdout, "", 0),
	}
}

// parseLevel 解析日志级别字符串
func parseLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

// levelToString 将日志级别转换为字符串
func levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// log 记录日志
func (l *Logger) log(level LogLevel, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := levelToString(level)
	
	var logMsg string
	if len(args) > 0 {
		if err, ok := args[0].(error); ok && err != nil {
			logMsg = fmt.Sprintf("%s [%s] %s: %v", timestamp, levelStr, msg, err)
		} else {
			logMsg = fmt.Sprintf("%s [%s] %s: %v", timestamp, levelStr, msg, args)
		}
	} else {
		logMsg = fmt.Sprintf("%s [%s] %s", timestamp, levelStr, msg)
	}

	l.logger.Println(logMsg)
	
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug 记录调试级别日志
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)
}

// Info 记录信息级别日志
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)
}

// Infof 记录格式化的信息级别日志
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(INFO, fmt.Sprintf(format, args...))
}

// Warn 记录警告级别日志
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)
}

// Error 记录错误级别日志
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
}

// Fatal 记录致命错误级别日志
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(FATAL, msg, args...)
}
