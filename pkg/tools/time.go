package tools

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/logc"
)

// TimeTransformToWeek 时间转换成周
func TimeTransformToWeek(ct time.Time) string {
	// 获取当前时间
	currentDate := ct.Format("2006-01-02")
	// 解析日期字符串为时间对象
	date, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Time Transform To Week failed, err: %s", err.Error()))
		return ""
	}
	return date.Weekday().String()
}

// TimeTransformToSeconds // 时间转换成秒
func TimeTransformToSeconds(ct time.Time) int {
	cs := ct.Hour()*3600 + ct.Minute()*60
	return cs
}

// FormatTimeToUTC 格式化为 UTC 时间
func FormatTimeToUTC(t int64) string {
	utcTime := time.Unix(t, 0).UTC()
	utcTimeString := utcTime.Format("2006-01-02T15:04:05.999Z")
	return utcTimeString
}

// ParserDuration 获取时间区间的开始时间
func ParserDuration(curTime time.Time, logScope int, timeType string) time.Time {
	duration, err := time.ParseDuration(strconv.Itoa(logScope) + timeType)
	if err != nil {
		logrus.Error(err.Error())
		return time.Time{}
	}
	startsAt := curTime.Add(-duration)
	return startsAt
}

// ParseTime 解析时间
func ParseTime(month string) (int, time.Month, int) {
	parsedTime, err := time.Parse("2006-01", month)
	if err != nil {
		return 0, time.Month(0), 0
	}
	curYear, curMonth, curDay := parsedTime.Date()
	return curYear, curMonth, curDay
}

// GetWeekday 获取周几
func GetWeekday(date string) (time.Weekday, error) {
	t, err := time.Parse("2006-1-2", date)
	if err != nil {
		return 0, err
	}

	weekday := t.Weekday()
	return weekday, nil
}

// IsEndOfWeek 判断是否是周末
func IsEndOfWeek(dateStr string) bool {
	date, err := time.Parse("2006-1-2", dateStr)
	if err != nil {
		return false
	}
	return date.Weekday() == time.Sunday
}
