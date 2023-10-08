package utils

import (
	"time"

	"github.com/jinzhu/now"
)

// 判断时间点处于今天
func IsToday(d int64) bool {
	tt := time.Unix(d, 0)
	now := time.Now()
	return tt.Year() == now.Year() && tt.Month() == now.Month() && tt.Day() == now.Day()
}

//是否是下一天
func IsNextToday(d int64) bool {
	d += 24 * 3600
	tt := time.Unix(d, 0)
	now := time.Now()
	return tt.Year() == now.Year() && tt.Month() == now.Month() && tt.Day() == now.Day()
}

// 判断是否是出于同一周
func IsSameWeek(d int64) bool {
	// 将时间戳转换为 time.Time 类型
	time1 := time.Unix(d, 0)
	time2 := time.Now()

	// 获取时间戳所属的年份和周数
	year1, week1 := time1.ISOWeek()
	year2, week2 := time2.ISOWeek()

	// 判断是否同一年同一周
	if year1 == year2 && week1 == week2 {
		return true
	} else {
		return false
	}
}

// 判断是否大于1周
func IsAfterWeek(d int64) bool {
	tt := time.Unix(d, 0)
	nowt := time.Now()
	if !tt.Before(nowt) {
		return false
	}
	at := now.With(tt).AddDate(0, 0, 7).Unix()
	return nowt.Unix() >= at
}

// 获取今天零点时间戳
func GetTodayZeroTime(curTime int64) int64 {
	currentTime := time.Unix(curTime, 0)
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	return startTime.Unix()
}

// 获取当前时间戳下一天0点时间戳
func GetZeroTime(curTime int64) int64 {
	currentTime := time.Unix(curTime, 0)
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	return startTime.Unix() + 86400 //3600*24
}

//是否是昨天
func IsYestoday(timestamp int64) bool {
	tt := time.Unix(timestamp, 0)
	yesTime := time.Now().AddDate(0, 0, -1)
	return tt.Year() == yesTime.Year() && tt.Month() == yesTime.Month() && tt.Day() == yesTime.Day()
}

//计算自然天数
func DiffDays(t1, t2 int64) int {
	if t1 == t2 {
		return -1
	}

	if t1 > t2 {
		t1, t2 = t2, t1
	}

	secOfDay := 3600 * 24
	diffDays := 0
	secDiff := t2 - t1
	if secDiff > int64(secOfDay) {
		tmpDays := int(secDiff / int64(secOfDay))
		t1 += int64(tmpDays) * int64(secOfDay)
		diffDays += tmpDays
	}
	st := time.Unix(t1, 0)
	et := time.Unix(t2, 0)
	dateformat := "20060102"
	if st.Format(dateformat) != et.Format(dateformat) {
		diffDays += 1
	}
	return diffDays
}
