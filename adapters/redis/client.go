package redis

import (
	"fmt"
	"github.com/chaksunshine/kit/catLog"
	"github.com/chaksunshine/kit/thread"
	"github.com/redis/go-redis/v9"
)

// 返回一个redisClient
// @author fuzeyu
// @date 2025/2/17
type Client struct {
	redis.Cmdable
}

// 关闭客户端
func (obj *Client) Close() {

	switch obj.Cmdable.(type) {

	// 关闭节点
	case *redis.Client:
		if err := obj.Cmdable.(*redis.Client).Close(); err != nil {
			catLog.Error(fmt.Sprintf("关闭redis节点失败 %s", err.Error()))
		}

	default:
		catLog.Error(fmt.Sprintf("关闭redis客户端的时候发现未知的客户端 %T", obj.Cmdable))
	}
}

// 获取当前节点信息
func (obj *Client) NodePort() (string, error) {
	local := thread.CtxRequestLocal()
	result, err := obj.Cmdable.ConfigGet(local, "port").Result()
	if err != nil {
		return "", err
	}
	return result["port"], nil
}

// @param client 客户端
func newClient(client redis.Cmdable) *Client {
	return &Client{
		Cmdable: client,
	}
}
