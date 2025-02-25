package date

import (
	"context"
	"time"
)

// 获取今日第一秒时间
func TodayFirstTime() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// 获取明日凌晨上下文
// @param sub 延迟时间
func TomorrowContext(sub ...time.Duration) context.Context {
	firstTime := TodayFirstTime().Add(time.Hour * 24)
	if len(sub) > 0 {
		firstTime = firstTime.Add(sub[0])
	}
	deadline, _ := context.WithDeadline(context.TODO(), firstTime)
	return deadline
}
