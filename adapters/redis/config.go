package redis

import (
	"errors"
	"fmt"
)

// redis 链接配置
type Config struct {

	// 集群名称
	Sentinel string `yaml:"sentinel"`

	// 节点配置
	Node []string `yaml:"node"`

	// 密码
	Pwd string `yaml:"pwd"`

	// 连接池大小
	PoolSize int `yaml:"poolSize"`

	// 单个连接时间
	// 单位:秒
	ConnectTimeout int `yaml:"connectTimeout"`

	// 使用的库文件
	Database int `yaml:"database"`
}

// 初始化
// @param database 库文件
func (obj *Config) init() error {

	if len(obj.Node) <= 0 {
		return errors.New(fmt.Sprintf("未知的链接地址"))
	}

	if obj.PoolSize == 0 {
		obj.PoolSize = 200
	}

	// 单个链接时间
	// 如果时间太长会导致当某一个阶段宕机后,长时间无法恢复
	if obj.ConnectTimeout == 0 {
		obj.ConnectTimeout = 10
	}
	return nil
}
