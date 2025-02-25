package kafka

import (
	"github.com/IBM/sarama"
	"sync"
	"time"
)

// 消息缓存对象
// @author fuzeyu
// @date 2025/2/24
type msgCache struct {

	// 消息缓存池
	length  int
	buffers []*sarama.ConsumerMessage

	// 锁
	mutex *sync.Mutex

	// 回调方法
	consumer Consumer

	// 执行回调
	execute chan struct{}

	batchSize  int           // 缓存池的数据达到了多少,必须执行一次
	lingerMs   time.Duration // 存在数据的话,多少秒必须执行一次
	checkTimer *time.Timer   // 强制执行时间
}

// 缓存时间信息
func (obj *msgCache) doTimer() *time.Timer {
	return obj.checkTimer
}

// 重置缓存信息
// @param timer 时间信息
func (obj *msgCache) resetBuffers() {
	obj.length = 0
	obj.buffers = make([]*sarama.ConsumerMessage, 0)
	obj.checkTimer.Reset(time.Minute)
}

// 执行一次消费信息
func (obj *msgCache) do() map[int32]int64 {

	var version = make(map[int32]int64)

	// 缓存信息
	obj.mutex.Lock()
	defer obj.mutex.Unlock()

	// 消费
	for _, buffer := range obj.buffers {

		// 回调方法
		obj.consumer(buffer)

		// 保存最新的版本信息
		// ! 报告消息是在另外回调中完成,可能会造成重复消费的情况
		if version[buffer.Partition] < buffer.Offset {
			version[buffer.Partition] = buffer.Offset
		}
	}

	// 更新版本信息
	obj.resetBuffers()
	return version
}

// 保存消息
// @param msg 消息信息
func (obj *msgCache) add(msg *sarama.ConsumerMessage) bool {

	var has bool

	obj.mutex.Lock()

	// 缓存信息
	obj.buffers = append(obj.buffers, msg)
	obj.length++

	// 是否达到了次数限制
	if obj.length >= obj.batchSize {
		has = true
	} else {
		obj.checkTimer.Reset(obj.lingerMs)
	}
	obj.mutex.Unlock()
	return has
}

// @param consumer 消息回调方法
func newMsgCache(consumer Consumer) *msgCache {
	c := &msgCache{
		mutex:      new(sync.Mutex),
		consumer:   consumer,
		execute:    make(chan struct{}),
		batchSize:  20,
		lingerMs:   time.Millisecond * 200,
		checkTimer: time.NewTimer(time.Hour),
	}

	// 初始化时间
	c.resetBuffers()
	return c
}
