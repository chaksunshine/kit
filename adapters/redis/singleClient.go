package redis

import (
	"errors"
	"fmt"
	"github.com/chaksunshine/kit/thread"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"time"
)

// 获取单节点redis客户端
// @param cfg 配置
// @param database 数据库
func NewSingleClient(cfg *Config) (*Client, error) {

	// 初始化配置
	if err := cfg.init(); err != nil {
		return nil, err
	}

	opt := &redis.Options{
		Addr:            cfg.Node[0],
		Password:        cfg.Pwd,
		PoolSize:        cfg.PoolSize,
		ConnMaxLifetime: time.Second * time.Duration(cfg.ConnectTimeout),
		DB:              cfg.Database,
	}
	rdb := redis.NewClient(opt)

	// 检查是否链接成功
	if err := rdb.Ping(thread.CtxRequestLocal()).Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("链接到redis单节点失败 %s", err.Error()))
	}
	return newClient(rdb), nil
}

// 按照字符串获取客户端信息
// @param configStr 配置
// @param database 数据库
func NewSingleClientByString(configStr string) (*Client, error) {
	var cfg *Config
	if err := yaml.Unmarshal([]byte(configStr), &cfg); err != nil {
		return nil, errors.New(fmt.Sprintf("解析配置文件失败 %s", err.Error()))
	}
	return NewSingleClient(cfg)
}
