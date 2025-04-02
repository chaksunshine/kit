package date

import (
	"errors"
	"fmt"
	"time"
)

// 解析时间
// @param value 时间格式
func ParserDate(value string) (time.Time, error) {

	// 按照固定长度解析
	var length = len(value)
	switch length {
	case 8:
		return time.ParseInLocation("20060102", value, time.UTC)

	case 14:
		return time.ParseInLocation("20060102150405", value, time.UTC)

	case 10:
		return time.ParseInLocation(time.DateOnly, value, time.UTC)

	case 19:
		return time.ParseInLocation(time.DateTime, value, time.UTC)

	default:
		return time.Time{}, errors.New(fmt.Sprintf("不支持的时间格式 %s 长度: %d", value, length))
	}
}
