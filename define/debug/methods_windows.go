package debug

import (
	"errors"
	"fmt"
	"github.com/chaksunshine/kit/catLog"
)

// 标记项目不允许发布
// @param data 选项信息
func NotRelease(data ...interface{}) error {
	if IsWindows == false {
		return errors.New("当前非windows环境,不允许发布")
	}

	for _, datum := range data {
		catLog.Error(fmt.Sprintf("项目收到不允许发布标记 %#v", datum))
	}
	return errors.New("项目标记了不允许发布")
}
