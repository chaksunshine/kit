package kafka

import (
	"errors"
	"github.com/IBM/sarama"
	"time"
)

// 生产者参数
type ProduceParams struct {
	eventType MsgRouterEventType
	ack       sarama.RequiredAcks
	cfg       *sarama.Config
}

// 获取结果
func (obj *ProduceParams) Builder() *sarama.Config {
	return obj.cfg
}

// 设置发送有效期
// @param batchSize 发送大小
// @param lingerMs 超时时间
func (obj *ProduceParams) ExpiredAt(batchSize int, lingerMs time.Duration) *ProduceParams {
	obj.cfg.Producer.Flush.Bytes = batchSize    // 包大小batch.size = 100KB
	obj.cfg.Producer.Flush.Frequency = lingerMs // 超时时间linger.ms = 100
	return obj
}

// 设置事务
// @param transactionName 事务名称
func (obj *ProduceParams) Transaction(transactionName string) *ProduceParams {
	obj.cfg.Producer.Transaction.ID = transactionName // 设置事务名称
	obj.cfg.Producer.Idempotent = true                // 开启幂等性
	obj.cfg.Producer.RequiredAcks = sarama.WaitForAll // 开启了事务必须要全部节点成功才算成功
	obj.cfg.Net.MaxOpenRequests = 1                   // 最大未确定数不可以超过1
	return obj
}

// 初始化
func (obj *ProduceParams) init() {

	cfg := sarama.NewConfig()
	switch obj.eventType {
	case ProduceEventTypeRoundLoop:
		cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	case ProduceEventTypeHashKey:
		cfg.Producer.Partitioner = sarama.NewHashPartitioner

	case ProduceEventTypeSelfPartitions:
		cfg.Producer.Partitioner = sarama.NewManualPartitioner

	default:
		panic(errors.New("未知的分区类型"))
	}

	obj.cfg = cfg
	cfg.Producer.RequiredAcks = obj.ack  // 主节点收到消息就算成功
	cfg.Producer.Return.Successes = true // 返回成功状态
	cfg.Producer.Retry.Max = 3           // 最大重试次数
	formatConfig(cfg)

	// 设置消息有效期
	obj.ExpiredAt(1024*16, time.Millisecond*80)
}

// @param eventType 分区类型
func NewProduceParams(eventType MsgRouterEventType, ack ...sarama.RequiredAcks) *ProduceParams {

	var rAck = sarama.WaitForLocal
	if len(ack) > 0 {
		rAck = ack[0]
	}

	p := &ProduceParams{
		eventType: eventType,
		ack:       rAck,
	}
	p.init()
	return p
}
