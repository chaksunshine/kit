package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/chaksunshine/kit/catLog"
	"sync"
	"time"
)

// 消费组消费
// @author fuzeyu
// @date 2025/2/24
type ConsumeGroup struct {
	config        *ClusterConfig
	workingNumber int
	topicName     string

	groupName    string
	kafkaVersion sarama.KafkaVersion

	// 设置设置了offset
	isOffset bool
	offset   int64

	// 消息缓存
	ctx context.Context
}

// 设置偏移量
func (obj *ConsumeGroup) SetOffset(offset int64) *ConsumeGroup {
	obj.offset = offset
	obj.isOffset = true
	return obj
}

// 创建客户端
// @param index 消费者下标
func (obj *ConsumeGroup) createClient(index int) (sarama.ConsumerGroup, error) {

	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = false // 关闭自动提交

	// 如果没有设置消费者变量,按照从头开始的方式消费
	config.Consumer.Offsets.Initial = obj.offset // 消费位置
	if obj.isOffset == false {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	config.Version = obj.kafkaVersion // 版本信息
	config.Consumer.Group.InstanceId = fmt.Sprintf("%s-%s-%d", consumeStandGroupName, obj.groupName, index)

	config.Consumer.Group.Session.Timeout = 15 * time.Second   // 会话超时时间
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second // 心跳间隔
	config.Consumer.MaxProcessingTime = 5 * time.Second        // 最大处理时间
	formatConfig(config)

	// 按照轮询的方式处理
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}

	// 创建消费者
	groupClient, err := sarama.NewConsumerGroup(obj.config.Brokers, obj.groupName, config)
	if err != nil {
		return nil, err
	}
	return groupClient, nil
}

// 开始消费
// @param ctx 上下文
// @param index 消费者下标
// @param call 回调方法
func (obj *ConsumeGroup) startConsume(ctx context.Context, index int, call Consumer) error {

	// 获取消费者
	client, err := obj.createClient(index)
	if err != nil {
		return err
	}
	defer client.Close()

	// 创建消费接口
	handler := newConsumeGroupHandler(ctx, obj.topicName, obj.groupName, index, newMsgCache(call))
	return client.Consume(ctx, []string{obj.topicName}, handler)
}

// 执行
// @param ctx 上下文
// @param call 回调
func (obj *ConsumeGroup) Run(ctx context.Context, call Consumer) {

	// 创建上下文
	localCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	obj.ctx = ctx

	// 批量消费
	wg := new(sync.WaitGroup)
	for index := 0; index < obj.workingNumber; index++ {
		wg.Add(1)
		go func(index int) {
			if err := obj.startConsume(localCtx, index, call); err != nil {
				catLog.Error(fmt.Sprintf("消费组消费者创建失败 %s", err.Error()))
			}
			wg.Done()
			cancel()
		}(index)
	}
	wg.Wait()

	// 等待结束
	<-localCtx.Done()
}

// @param config 配置信息
// @param topicName 主题名
// @param workingNumber 工作线程数
// @param groupName 消费者名称
func NewConsumeGroup(config *ClusterConfig, topicName string, workingNumber int, groupName string) *ConsumeGroup {
	c := &ConsumeGroup{
		config:        config,
		topicName:     topicName,
		workingNumber: workingNumber,
		groupName:     groupName,
		kafkaVersion:  kafkaDefaultVersion,
	}
	return c
}
