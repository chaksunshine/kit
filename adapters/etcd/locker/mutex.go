package locker

import (
	"context"
	"github.com/chaksunshine/kit/thread"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"time"
)

// 分布式锁
// @author fuzeyu
// @date 2025/3/12
type Mutex struct {
	etcdClient *clientv3.Client
	prefix     string

	// 等待时间
	maxWait    time.Duration
	waitSecond int
}

// 执行任务
func (obj *Mutex) Do(call func()) error {

	// 设置有效期
	timeout, cancelFunc := context.WithTimeout(context.TODO(), obj.maxWait)
	defer cancelFunc()

	// 创建一个租约
	session, _ := concurrency.NewSession(obj.etcdClient, concurrency.WithTTL(obj.waitSecond))
	defer session.Close()
	mutex := concurrency.NewMutex(session, obj.prefix)

	// 创建锁
	if err := mutex.Lock(timeout); err != nil {
		return err
	}

	// 执行回调
	call()

	return mutex.Unlock(thread.CtxRequest())
}

// @param etcdClient 客户端
// @param prefix 前缀
// @param maxWait 最大等待时间
func NewMutex(etcdClient *clientv3.Client, prefix string, watSeconds int) *Mutex {
	c := &Mutex{
		etcdClient: etcdClient,
		prefix:     prefix,
		waitSecond: watSeconds,
		maxWait:    time.Second * time.Duration(watSeconds),
	}
	return c
}
