package kafka

import (
	"github.com/IBM/sarama"
	"time"
)

// kafka集群配置
type ClusterConfig struct {
	Brokers []string
}

// 消息路由类型
type MsgRouterEventType int

const (
	ProduceEventTypeRoundLoop      MsgRouterEventType = iota // 按照轮询分区
	ProduceEventTypeHashKey                                  // 按照哈希key
	ProduceEventTypeSelfPartitions                           // 自定义实现分区
)

const (

	// 默认的普通消费组分组名称
	consumeStandGroupName = "x-default"
)

var (

	// kafka版本
	kafkaDefaultVersion = sarama.V3_1_0_0
)

// 消费者接口
// @param topic 消息主题
// @param msg 消息信息
type Consumer func(msg *sarama.ConsumerMessage)

// 格式化配置信息
// @param cfg 配置信息
func formatConfig(cfg *sarama.Config) {
	cfg.Metadata.RefreshFrequency = time.Second * 10 // 更新元数据
}
