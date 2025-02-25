package kafka

import (
	"github.com/IBM/sarama"
)

// kafka 管理器
// @author fuzeyu
// @date 2025/2/21
type Manager struct {
	config *ClusterConfig
	client sarama.ClusterAdmin
}

// 创建链接
func (obj *Manager) connect() error {
	admin, err := sarama.NewClusterAdmin(obj.config.Brokers, sarama.NewConfig())
	if err != nil {
		return err
	}
	obj.client = admin
	return nil
}

// 创建分区的参数
// @param partitions 分区数量
// @param replication 副本数量
func (obj *Manager) CrateTopicParams(partitions, replication int) *sarama.TopicDetail {
	var policy = "delete"
	return &sarama.TopicDetail{
		NumPartitions:     int32(partitions),
		ReplicationFactor: int16(replication),
		ConfigEntries: map[string]*string{
			"cleanup.policy": &policy,
		},
	}
}

// 创建主题
// @param topic 主题名称
// @param params 参数
func (obj *Manager) CreateTopic(topic string, params *sarama.TopicDetail) error {
	return obj.client.CreateTopic(topic, params, false)
}

// 删除主题
// @param topic 主题名称
func (obj *Manager) DeleteTopic(topic string) error {
	return obj.client.DeleteTopic(topic)
}

// 集群配置
// @param config 配置信息
func NewManager(config *ClusterConfig) (*Manager, error) {
	c := &Manager{
		config: config,
	}
	return c, c.connect()
}
