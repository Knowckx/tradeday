package tradeday

import (
	"errors"
	"time"
)

var ErrUnsupportedCalendar = errors.New("tradeday: unsupported calendar")

type CalendarID string

const (
	CalendarCNStock CalendarID = "cn_stock"
	CalendarUSStock CalendarID = "us_stock"
)

// Calendar 表示某个交易所的交易日历。
type Calendar interface {
	// IsTradeDay 判断给定日期是否为交易日。
	IsTradeDay(day time.Time) bool
	// TradeDays 返回闭区间 [start, end] 内的交易日列表。
	TradeDays(start, end time.Time) []time.Time
	// OffsetTradeDay 返回给定日期前后第 N 个交易日。
	OffsetTradeDay(day time.Time, offset int) (time.Time, error)
	// PrevTradeDay 返回给定日期的前一个交易日。
	PrevTradeDay(day time.Time) (time.Time, error)
	// NextTradeDay 返回给定日期的后一个交易日。
	NextTradeDay(day time.Time) (time.Time, error)
}

// New 根据日历标识创建对应的交易日历。
func New(calendarID CalendarID) (Calendar, error) {
	if calendarID == "" {
		return nil, ErrUnsupportedCalendar
	}

	return nil, ErrUnsupportedCalendar
}
