package kafka

import (
	"github.com/IBM/sarama"
)

// 生产消息
// @author fuzeyu
// @date 2025/2/21
type Producer struct {

	// 消息配置
	config      *ClusterConfig
	producerCfg *sarama.Config

	// 客户端信息
	synProducer sarama.SyncProducer
}

// 获取生产者对象
func (obj *Producer) SynProducer() sarama.SyncProducer {
	return obj.synProducer
}

// 发送一条信息
// @param synProducer 消息信息
func (obj *Producer) Send(producer *sarama.ProducerMessage) (int32, int64, error) {
	return obj.synProducer.SendMessage(producer)
}

// 发送消息
// @param topic 主题
// @param msg 消息信息
func (obj *Producer) SendText(topic string, msg string, partitions ...int32) (int32, int64, error) {

	producer := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	// 确定分区信息
	if len(partitions) > 0 {
		producer.Partition = partitions[0]
	}
	return obj.Send(producer)
}

// 链接到集群
func (obj *Producer) connect() error {
	producer, err := sarama.NewSyncProducer(obj.config.Brokers, obj.producerCfg)
	if err != nil {
		return err
	}
	obj.synProducer = producer
	return nil
}

// @param config 集群配置
// @param producerCdg 消息发送配置
func NewProducer(config *ClusterConfig, producerCfg *sarama.Config) (*Producer, error) {
	c := &Producer{
		config:      config,
		producerCfg: producerCfg,
	}
	return c, c.connect()
}
