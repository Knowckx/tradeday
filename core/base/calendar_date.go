package base

import (
	"time"
)

// CalendarDate 是日历内部使用的日期载体。
// 它同时保留规范化后的日期字符串和解析后的 time.Time，
// 方便后续业务逻辑按需直接取用。
type CalendarDate struct {
	date Date
	time time.Time
}

func NewCalendarDateFromTime(day time.Time) *CalendarDate {
	return &CalendarDate{
		date: Date(day.Format(DateLayout)),
		time: day,
	}
}

func (day *CalendarDate) Date() Date {
	return day.date
}

func (day *CalendarDate) Time() time.Time {
	return day.time
}

func (day *CalendarDate) IsSupportedYear(minYear, maxYear int) bool {
	year := day.time.Year()
	return year >= minYear && year <= maxYear
}

func (day *CalendarDate) IsWeekend() bool {
	weekday := day.time.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}
