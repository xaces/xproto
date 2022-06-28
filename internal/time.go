package internal

import (
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

// TimeNow 当前时间
func TimeNow() string {
	return time.Now().Format(timeFormat)
}

// ToUnixTime 格式化字符串为时间戳
func ToUnixTime(dt string) int64 {
	tm, _ := time.Parse(timeFormat, dt)
	return tm.Unix()
}
