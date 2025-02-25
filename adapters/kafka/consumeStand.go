package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/chaksunshine/kit/catLog"
)

// 独立消费组
// @author fuzeyu
// @date 2025/2/24
type ConsumeStand struct {
	config       *ClusterConfig
	kafkaVersion sarama.KafkaVersion

	// 消费的分区
	topicName  string
	partitions int32

	// 上一次的回调信息
	offset int64

	// 分区名称
	groupName string

	// 消息缓存对象
	msgCache *msgCache

	// 客户端
	kafkaClient   sarama.Client
	partition     sarama.PartitionOffsetManager
	offsetManager sarama.OffsetManager
}

// 打印日志
// @param msg 消息信息
func (obj *ConsumeStand) logger(msg string) {
	catLog.Info(fmt.Sprintf("kafka普通消费 %s/%s/%d %s", obj.topicName, obj.groupName, obj.partitions, msg))
}

// 设置名称
func (obj *ConsumeStand) SetGroupName(groupName string) {
	obj.groupName = groupName
}

// 获取上一次分区信息
func (obj *ConsumeStand) loadLastOffset() error {

	config := sarama.NewConfig()
	config.Version = obj.kafkaVersion

	// 获取管理客户端
	admin, err := sarama.NewClusterAdmin(obj.config.Brokers, config)
	if err != nil {
		return err
	}
	defer admin.Close()

	// 获取所有的消费组
	parNum := obj.partitions
	partitions := map[string][]int32{
		obj.topicName: {
			parNum,
		},
	}
	offsets, err := admin.ListConsumerGroupOffsets(obj.groupName, partitions)
	if err != nil {
		return err
	}
	resp := offsets.Blocks[obj.topicName][parNum]
	obj.offset = resp.Offset
	return nil
}

// 创建客户端链接
func (obj *ConsumeStand) createConnect() error {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = obj.kafkaVersion
	config.Consumer.Offsets.AutoCommit.Enable = false
	formatConfig(config)

	client, err := sarama.NewClient(obj.config.Brokers, config)
	if err != nil {
		return err
	}

	// 偏移量管理器
	offsetManager, err := sarama.NewOffsetManagerFromClient(obj.groupName, client)
	if err != nil {
		return err
	}

	// 分区管理器
	partition, err := offsetManager.ManagePartition(obj.topicName, obj.partitions)
	if err != nil {
		return err
	}
	obj.kafkaClient = client
	obj.offsetManager = offsetManager
	obj.partition = partition
	return nil
}

// 关闭客户端链接
func (obj *ConsumeStand) disconnect() {

	// kafka客户端
	if obj.kafkaClient != nil {
		_ = obj.kafkaClient.Close()
	}

	// 偏移量管理器
	if obj.offsetManager != nil {
		_ = obj.offsetManager.Close()
	}

	// 分区客户端
	if obj.partition != nil {
		_ = obj.partition.Close()
	}
}

// 报告消息信息
func (obj *ConsumeStand) report(versions map[int32]int64) {

	// 检查版本
	offset := versions[obj.partitions]
	if offset <= 0 {
		return
	}

	// 报告偏移量
	obj.partition.MarkOffset(offset+1, "")
	obj.offsetManager.Commit()
}

// 开始消费
// @param ctx 上下文
func (obj *ConsumeStand) consume(ctx context.Context) error {

	// 创建客户端
	client, err := sarama.NewConsumerFromClient(obj.kafkaClient)
	if err != nil {
		return err
	}
	defer client.Close()

	// 指定分区和下标开始消费
	partition, err := client.ConsumePartition(obj.topicName, obj.partitions, obj.offset)
	if err != nil {
		return err
	}
	defer partition.Close()

	for {

		select {

		// 消费完成
		case <-ctx.Done():
			return nil

		// 拉取消息
		case msg := <-partition.Messages():
			if obj.msgCache.add(msg) {
				obj.report(obj.msgCache.do())
			}

		// 按照时间
		case <-obj.msgCache.checkTimer.C:
			obj.report(obj.msgCache.do())
		}
	}
}

// 开始执行
// @param ctx 上下文
func (obj *ConsumeStand) Run(ctx context.Context, consumer Consumer) error {

	obj.msgCache = newMsgCache(consumer)

	// 获取上一次分区信息
	if err := obj.loadLastOffset(); err != nil {
		return err
	}

	// 加载客户端信息
	if err := obj.createConnect(); err != nil {
		return err
	}
	defer obj.disconnect()

	// 查询上下文
	localCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	go func() {
		if err := obj.consume(ctx); err != nil {
			obj.logger(fmt.Sprintf("消费分区失败  %s", err.Error()))
		}
		cancelFunc()
	}()

	// 等待消费完成
	<-localCtx.Done()
	return nil
}

// @param config 配置信息
// @param topicName 主题名
// @param partitions 分区信息
func NewConsumeStand(config *ClusterConfig, topicName string, partitions int32) *ConsumeStand {
	c := &ConsumeStand{
		config:       config,
		topicName:    topicName,
		partitions:   partitions,
		groupName:    consumeStandGroupName,
		kafkaVersion: kafkaDefaultVersion,
	}
	return c
}
