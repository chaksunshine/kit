package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/chaksunshine/kit/catLog"
)

// 按照消费组的方式创建消息
// @author fuzeyu
// @date 2025/2/24
type consumeGroupHandler struct {
	msgCache *msgCache

	index int
	group string
	ctx   context.Context
	topic string
}

// 输出日志
// @param msg 消息信息
func (obj *consumeGroupHandler) logger(msg string) {
	catLog.Info(fmt.Sprintf("消费者信息 %s/%s-%d %s", obj.topic, obj.group, obj.index, msg))
}

func (obj *consumeGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	obj.logger(`消费者初始化`)
	return nil
}

func (obj *consumeGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	obj.logger(`消费者关闭`)
	return nil
}

// 报告消息信息
// @param result 消息结果
func (obj *consumeGroupHandler) report(session sarama.ConsumerGroupSession, result map[int32]int64) {
	if len(result) <= 0 {
		return
	}
	for partitions, offset := range result {
		session.MarkOffset(obj.topic, partitions, offset+1, "")
	}
	session.Commit()
}

// 消费者消费
// @param session 客户端
// @param claim 消息渠道
func (obj *consumeGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for {

		select {

		// 消费完成
		case <-obj.ctx.Done():
			return nil

		// 拉取消息
		case msg := <-claim.Messages():
			if obj.msgCache.add(msg) {
				obj.report(session, obj.msgCache.do())
			}

		// 按照时间
		case <-obj.msgCache.checkTimer.C:
			obj.report(session, obj.msgCache.do())
		}
	}
}

// @param ctx 上下文
// @param topic 主题
// @param group 分组
// @param index 下标
// @param cacheBuffers 缓存信息
func newConsumeGroupHandler(ctx context.Context, topic string, group string, index int, msgCache *msgCache) *consumeGroupHandler {
	c := &consumeGroupHandler{
		ctx:      ctx,
		topic:    topic,
		group:    group,
		index:    index,
		msgCache: msgCache,
	}
	return c
}
