package thread

import (
	"context"
	"time"
)

// 本地类数据请求
func CtxRequestLocal() context.Context {
	return CtxRequest(2)
}

// 获取正常请求
func CtxRequest(seconds ...int) context.Context {
	var value = 3
	if len(seconds) > 0 {
		value = seconds[0]
	}
	timeout, _ := context.WithTimeout(context.TODO(), time.Second*time.Duration(value))
	return timeout
}

// 检查一个上下文是否取消
// @param ctx 上下文
func IsCancel(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
