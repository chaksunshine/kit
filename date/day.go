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
// @param after 向后延迟的时间
func TomorrowContext(after ...time.Duration) context.Context {
	firstTime := TodayFirstTime().Add(time.Hour * 24)
	if len(after) > 0 {
		firstTime = firstTime.Add(after[0])
	}
	deadline, _ := context.WithDeadline(context.TODO(), firstTime)
	return deadline
}
