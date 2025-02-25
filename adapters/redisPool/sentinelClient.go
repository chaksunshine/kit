package redisPool

import (
	"errors"
	"fmt"
	"github.com/chaksunshine/kit/thread"
	"github.com/redis/go-redis/v9"
	"time"
)

// 哨兵配置
// @author fuzeyu
// @date 2025/2/17
type SentinelClient struct {
	master *Client
	slaver *Client
	cfg    *Config
}

// 获取链接
func (obj *SentinelClient) Master() *Client {
	return obj.master
}

// 获取链接
func (obj *SentinelClient) Slaver() *Client {
	return obj.slaver
}

// 创建一个链接
// @param master 写节点
// @param read 是否只读
func (obj *SentinelClient) createConnect(read bool) (*Client, error) {

	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:      obj.cfg.Sentinel,
		SentinelAddrs:   obj.cfg.Node,
		ReplicaOnly:     read,
		Password:        obj.cfg.Pwd,
		DB:              obj.cfg.database,
		PoolSize:        obj.cfg.PoolSize,
		ConnMaxLifetime: time.Second * time.Duration(obj.cfg.ConnectTimeout),
	})

	if err := rdb.Ping(thread.CtxRequestLocal()).Err(); err != nil {
		return nil, err
	}
	return newClient(rdb), nil
}

// 注册节点
// @param master 写节点
// @param slave 读节点
func (obj *SentinelClient) connect() error {

	// 获取master链接
	if clients, err := obj.createConnect(false); err != nil {
		return errors.New(fmt.Sprintf("链接到redis哨兵集群master服务失败 %s", err.Error()))
	} else {
		obj.master = clients
	}

	// 获取slaver链接
	if clients, err := obj.createConnect(true); err != nil {
		return errors.New(fmt.Sprintf("链接到redis哨兵集群slaver服务失败 %s", err.Error()))
	} else {
		obj.slaver = clients
	}
	return nil
}

// 获取哨兵主从链接配置
// @param cfg 配置 主节点
// @param database 库
func NewSentinelClient(cfg *Config, database int) (*SentinelClient, error) {

	// 初始化配置
	if err := cfg.init(database); err != nil {
		return nil, err
	}

	var self = &SentinelClient{
		cfg: cfg,
	}
	return self, self.connect()
}
