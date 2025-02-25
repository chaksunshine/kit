package redisPool

import (
	"errors"
	"fmt"
	"github.com/chaksunshine/kit/thread"
	"github.com/redis/go-redis/v9"
	"time"
)

// 获取单节点redis客户端
// @param cfg 配置
// @param database 数据库
func NewSingleClient(cfg *Config, databases int) (*Client, error) {

	// 初始化配置
	if err := cfg.init(databases); err != nil {
		return nil, err
	}

	opt := &redis.Options{
		Addr:            cfg.Node[0],
		Password:        cfg.Pwd,
		PoolSize:        cfg.PoolSize,
		ConnMaxLifetime: time.Second * time.Duration(cfg.ConnectTimeout),
		DB:              cfg.database,
	}
	rdb := redis.NewClient(opt)

	// 检查是否链接成功
	if err := rdb.Ping(thread.CtxRequestLocal()).Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("链接到redis单节点失败 %s", err.Error()))
	}
	return newClient(rdb), nil
}
