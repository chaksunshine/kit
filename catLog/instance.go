package catLog

import (
	"fmt"
	"github.com/chaksunshine/kit/date"
	"github.com/chaksunshine/kit/path"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 日志信息
// @author fuzeyu
// @date 2025/2/23
type instance struct {
	rwMutex *sync.RWMutex

	// 上次日志对象
	logger     *zap.Logger
	lastLogger *os.File
}

// 清理文件
func (obj *instance) clearFile() {
	var current = time.Now()
	_ = filepath.Walk(logDir, func(path string, entry os.FileInfo, err error) error {
		if err == nil && entry.IsDir() == false && strings.HasSuffix(entry.Name(), logExt) {
			if sub := current.Sub(entry.ModTime()); sub >= logSaveDay {
				os.Remove(path)
				println(fmt.Sprintf("删除日志文件 %s", path))
			}
		}
		return nil
	})
}

// 切换任务
func (obj *instance) autoEvent() {

	for {

		// 清理文件
		obj.clearFile()

		// 等到明天
		context := date.TomorrowContext()
		<-context.Done()

		// 切换日志
		obj.reConnect()
	}
}

// 重置日志信息
func (obj *instance) reConnect() {

	obj.rwMutex.Lock()
	defer obj.rwMutex.Unlock()

	// 关闭上一次
	if obj.lastLogger != nil {
		obj.lastLogger.Close()
		_ = obj.logger.Sync()
	}

	// 初始化目录
	base := filepath.Base(os.Args[0])
	if strings.Contains(base, ".") {
		base = strings.TrimSuffix(base, filepath.Ext(base))
	}

	// 创建文件
	filePath := filepath.Join(logDir, fmt.Sprintf("%s.%s%s", base, time.Now().Format(time.DateOnly), logExt))
	if err := path.FolderMustCreate(logDir); err != nil {
		println(fmt.Sprintf("切换日志失败 %s", err.Error()))
		return
	}

	// 创建文件夹聚丙
	resource, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		println(fmt.Sprintf("切换日志失败 %s", err.Error()))
		return
	}
	obj.lastLogger = resource
	_, _ = resource.Write([]byte(fmt.Sprintf("\n\nlog init in %s\n", time.Now())))

	// 2. 定义日志编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder, // 日志级别小写
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	// 3. 创建两个输出目标的核心
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 控制台输出为普通文本格式
		zapcore.Lock(os.Stdout),                  // 输出到控制台
		zapcore.DebugLevel,                       // 设置日志级别
	)

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 文件输出为 JSON 格式
		zapcore.AddSync(resource),             // 输出到文件
		zapcore.InfoLevel,                     // 设置日志级别
	)

	core := zapcore.NewTee(consoleCore, fileCore)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	obj.logger = logger
}

// 保存基本信息
// @param msg 消息信息
func (obj *instance) info(msg string) {
	obj.rwMutex.RLock()
	obj.logger.Info(msg)
	obj.rwMutex.RUnlock()
}

// 记录错误日志
// @param msg 消息信息
func (obj *instance) error(msg string) {
	obj.rwMutex.RLock()
	obj.logger.Error(msg)
	obj.rwMutex.RUnlock()
}

// 创建客户端
func newInstance() *instance {
	c := &instance{
		rwMutex: new(sync.RWMutex),
	}
	c.reConnect()
	go c.autoEvent()
	return c
}

var _instance = newInstance()
