package Logger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// 日志轮转钩子结构体
type TimeRotatingHook struct {
	mu        sync.Mutex
	dir       string
	file      *os.File
	filename  string
	lastWrite time.Time
	interval  time.Duration
}

// 创建一个新的轮转日志钩子
func NewTimeRotatingHook(dir string, interval time.Duration) (*TimeRotatingHook, error) {
	hook := &TimeRotatingHook{
		dir:       dir,
		interval:  interval,
		lastWrite: time.Now().Truncate(interval),
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	if err := hook.rotate(); err != nil {
		return nil, err
	}
	return hook, nil
}

// 轮转日志文件
func (h *TimeRotatingHook) rotate() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.file != nil {
		h.file.Close()
	}

	now := time.Now().Truncate(h.interval)
	h.filename = filepath.Join(h.dir, fmt.Sprintf("server-%s.log", now.Format("2006-01-02")))
	file, err := os.OpenFile(h.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	h.file = file
	return nil
}

// Fire 在每次日志调用时触发
func (h *TimeRotatingHook) Fire(entry *log.Entry) error {
	now := time.Now().Truncate(h.interval)
	if now.After(h.lastWrite) {
		if err := h.rotate(); err != nil {
			return err
		}
		h.lastWrite = now
	}

	line, err := entry.String()
	if err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err = h.file.WriteString(line)
	return err
}

// Levels 指定 Hook 触发的日志级别
func (h *TimeRotatingHook) Levels() []log.Level {
	return log.AllLevels
}

// 自定义日志格式
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	_, err := fmt.Fprintf(b, "[%s] %s %s\n", level, timestamp, entry.Message)
	return b.Bytes(), err
}

// 初始化日志系统
func InitLogger() error {
	hook, err := NewTimeRotatingHook("../serverLogs", 24*time.Hour) // 每天一个文件
	if err != nil {
		return err
	}

	log.SetFormatter(&CustomFormatter{})
	log.SetLevel(log.InfoLevel) // 设置日志级别
	log.SetOutput(os.Stdout)
	log.AddHook(hook) // 注册自定义 hook

	return nil
}
