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

// 实现钩子的结构体
type TimeRotatingHook struct {
	mu        sync.Mutex //互斥锁 mutex通过mu.Lock()，defer mu.Unlock()来实现安全访问或修改共享资源
	dir       string
	file      *os.File
	filename  string
	lastWrite time.Time
}

// 需要日志文件路径，增加日志的时间间隔，结果是返回一个hook实例
func NewTimeRotatingHook(dir string, interval time.Duration) (*TimeRotatingHook, error) {
	hook := &TimeRotatingHook{
		dir:       dir,
		lastWrite: time.Now().Truncate(interval),
		//lastWrite被设置为当前时间向下取整到最接近的interval间隔
		//需要找到不大于当前时间的、最接近的时间，就是当前时间减去24小时但是又不超出当天的范围，只能是当天0点
	}
	//创建文件夹
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	hook.rotate()
	return hook, nil
}

func (h *TimeRotatingHook) rotate() error {
	//在修改文件或文件名之前，先获取互斥锁以确保线程安全。
	h.mu.Lock()
	defer h.mu.Unlock()
	//如果当前有打开的文件，则先关闭它。
	if h.file != nil {
		h.file.Close()
	}
	now := time.Now().Truncate(24 * time.Hour) // 每天一个文件
	//产生文件名
	h.filename = filepath.Join(h.dir, fmt.Sprintf("server-%s.log", now.Format("2006-01-02")))
	//尝试以追加模式打开（如果不存在则创建）新的日志文件。
	file, err := os.OpenFile(h.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	h.file = file
	return nil
}

// 钩子函数，开启log之前会调用
func (h *TimeRotatingHook) Fire(entry *log.Entry) error {
	now := time.Now().Truncate(24 * time.Hour) //每天一个文件
	//检查写入日志的当前时间是否晚于最后一次写入时间，如果是，则需要轮转日志（创建新的日志文件
	if now.After(h.lastWrite) {
		if err := h.rotate(); err != nil {
			return err
		}
		//更新上一条日志写入的时间
		h.lastWrite = now
	}
	//将日志条目写入到当前日志文件中。
	//str, _ := entry.String()
	//_, err := h.file.WriteString(str + "\n")
	return nil
}

func (h *TimeRotatingHook) Levels() []log.Level {
	return log.AllLevels
}

type CustomFormatter struct{}

// Format formats a log entry
func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	// Write the formatted log entry to the buffer
	//_, err := fmt.Fprintf(b, "[%s] %s %s\n", strings.ToUpper(entry.Level.String()), timestamp, fmt.Sprintf("%s %s", entry.Message, caller))
	_, err := fmt.Fprintf(b, "[%s] %s\n", strings.ToUpper(entry.Level.String()), timestamp, fmt.Sprintf("%s %s", entry.Message))
	return b.Bytes(), err
}
func InitLogger() (err error) {
	hook, err := NewTimeRotatingHook("./serverLogs", 24*time.Hour) //每天一个文件
	if err != nil {
		return err
	}
	log.AddHook(hook)
	log.SetOutput(hook.file)
	log.SetFormatter(&CustomFormatter{})
	//log.SetReportCaller(true)
	return nil
}
