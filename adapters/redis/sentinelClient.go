package redis

import (
	"errors"
	"fmt"
	"github.com/chaksunshine/kit/thread"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"time"
)

// 哨兵配置
// @author fuzeyu
// @date 2025/2/17
type SentinelClient struct {
	client *Client
	cfg    *Config
}

// 获取链接
func (obj *SentinelClient) Client() *Client {
	return obj.client
}

// 创建一个链接
// @param master 写节点
// @param read 是否只读
func (obj *SentinelClient) createConnect() (*Client, error) {
	rdb := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:       obj.cfg.Sentinel,
		SentinelAddrs:    obj.cfg.Node,
		Password:         obj.cfg.Pwd,
		SentinelPassword: obj.cfg.Pwd,
		DB:               obj.cfg.Database,
		PoolSize:         obj.cfg.PoolSize,
		ConnMaxLifetime:  time.Second * time.Duration(obj.cfg.ConnectTimeout),
		//RouteByLatency:   true, // 读操作自动选择低延迟的从节点
		RouteRandomly: true, // 读操作随机选择从节点实现负载均衡
		ReplicaOnly:   true,
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

	connect, err := obj.createConnect()
	if err != nil {
		return errors.New(fmt.Sprintf("链接到哨兵集群信息失败 %s", err.Error()))
	}
	obj.client = connect
	return nil
}

// 获取哨兵主从链接配置
// @param cfg 配置 主节点
// @param database 库
func NewSentinelClient(cfg *Config) (*SentinelClient, error) {

	// 初始化配置
	if err := cfg.init(); err != nil {
		return nil, err
	}

	var self = &SentinelClient{
		cfg: cfg,
	}
	return self, self.connect()
}

// 解析配置
func NewSentinelClientByString(configStr string) (*SentinelClient, error) {
	var cfg *Config
	if err := yaml.Unmarshal([]byte(configStr), &cfg); err != nil {
		return nil, errors.New(fmt.Sprintf("解析配置文件失败 %s", err.Error()))
	}
	return NewSentinelClient(cfg)
}
