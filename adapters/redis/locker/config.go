package locker

import "time"

// 分布式锁配置
type Config struct {

	// 单次续期内耗时
	LockTimeout time.Duration

	// 抢锁失败后
	// 下一次重试时间
	GetLockTimer time.Duration

	// 续期时间
	RenewalTimer time.Duration
}

var defaultConfig = &Config{
	LockTimeout:  time.Second * 5,
	GetLockTimer: time.Millisecond * 200,
	RenewalTimer: time.Second,
}
