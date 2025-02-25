package redisPool

import (
	"github.com/chaksunshine/kit/thread"
	"github.com/redis/go-redis/v9"
)

// 返回一个redisClient
// @author fuzeyu
// @date 2025/2/17
type Client struct {
	redis.Cmdable
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
