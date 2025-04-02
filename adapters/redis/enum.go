package redis

type RedisType int

const (

	// redis单机服务
	RedisSingle RedisType = iota

	// redis哨兵服务
	RedisSentinel
)
