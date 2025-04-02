package catLog

import (
	"os"
	"time"
)

// 记录消息信息
// @param msg 信息
func Info(msg string) {
	_instance.info(msg)
}

// 记录错误信息
// @param msg 信息
func Error(msg string) {
	_instance.error(msg)
}

// 严重错误信息
func FatalError(msg string) {
	_instance.error(msg)
	time.Sleep(time.Second * 5)
	os.Exit(2)
}
