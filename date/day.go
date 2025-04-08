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

// 获取本周的最后一天时间
// @param week 获取周几的时间
func WeekDay(week ...int) time.Time {

	var card time.Weekday = 7
	if len(week) > 0 && week[0] > 0 && week[0] < 8 {
		card = time.Weekday(week[0])
	}

	now := time.Now()
	weekday := now.Weekday()
	daysUntilSunday := (time.Sunday - weekday + card) % card
	sunday := now.Add(time.Hour * 24 * time.Duration(daysUntilSunday))
	return time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 0, 0, 0, 0, sunday.Location())
}
