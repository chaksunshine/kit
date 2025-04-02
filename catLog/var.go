package catLog

import (
	"github.com/chaksunshine/kit/define"
	"path/filepath"
	"time"
)

var (

	// 获取日志文件路径
	logDir = filepath.Join(define.RootPath, "logs")

	// 日志路径
	logExt = ".log"

	// 日志保留天数
	logSaveDay = 7 * time.Hour * 24
)
