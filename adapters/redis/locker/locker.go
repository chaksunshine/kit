package locker

import (
	"context"
	"fmt"
	"github.com/chaksunshine/kit/adapters/redis"
	"github.com/chaksunshine/kit/thread"
	"github.com/chaksunshine/kit/unique"
	"github.com/chaksunshine/kit/unique/id"
	"time"
)

// redis分布式锁
// @author fuzeyu
// @date 2025/3/4
type Locker struct {
	lockName    string
	redisClient *redis.Client

	// 锁的版本号,用于去从
	version string

	// 上下文监听
	context context.Context
	cancel  context.CancelFunc

	// 配置
	config *Config
}

// 检查是否是当前携程获取到的锁
func (obj *Locker) currentIsLocker() bool {
	cmd := obj.redisClient.Get(thread.CtxRequestLocal(), obj.lockName)
	result, err := cmd.Result()
	return err == nil && result == obj.version
}

// 开始续期
func (obj *Locker) startListen() {

	for {

		// 检查是否取消
		if thread.IsCancel(obj.context) {
			return
		}

		// 等待
		time.Sleep(obj.config.RenewalTimer)

		// 检查是否要续期
		if obj.currentIsLocker() == false {
			obj.cancel()
		} else {
			obj.redisClient.Expire(thread.CtxRequestLocal(), obj.lockName, obj.config.LockTimeout)
		}
	}
}

// 加锁
func (obj *Locker) Lock() error {
	for {
		cmd := obj.redisClient.SetNX(thread.CtxRequestLocal(), obj.lockName, obj.version, obj.config.LockTimeout)
		result, err := cmd.Result()
		if err != nil {
			return err
		}
		if result {
			obj.context, obj.cancel = context.WithCancel(context.TODO())
			go obj.startListen()
			return nil
		}
		time.Sleep(obj.config.GetLockTimer)
	}
}

// 解锁
func (obj *Locker) Unlock() {

	if obj.currentIsLocker() == false {
		return
	}

	// 删除锁
	obj.redisClient.Del(thread.CtxRequestLocal(), obj.lockName)
	obj.cancel()
}

// @param redisClient redis客户端
// @param lockName 锁名称
func NewLocker(redisClient *redis.Client, lockName string, cfg ...*Config) *Locker {

	var cfgValue = defaultConfig
	if len(cfg) > 0 {
		cfgValue = cfg[0]
	}

	c := &Locker{
		lockName:    lockName,
		redisClient: redisClient,
		version:     fmt.Sprintf("%d-%s", id.Serial.Next(), unique.NewUUID()),
		config:      cfgValue,
	}
	return c
}
