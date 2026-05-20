package core

// CalendarID 表示某一类市场日历的唯一标识。
type CalendarID string

const (
	// CalendarCNStock 表示中国 A 股市场日历。
	CalendarCNStock CalendarID = "cn_stock"
	// CalendarUSStock 表示美国股票市场日历。
	CalendarUSStock CalendarID = "us_stock"
)

// Calendar 表示某一类市场的交易日历。
// 所有日期参数都必须使用 2006-01-02 格式。
type Calendar interface {
	// IsTradeDay 判断给定日期是否为交易日。
	// day 必须为 2006-01-02 格式。
	// 当日期格式非法或超出该日历支持范围时，返回 error。
	IsTradeDay(day Date) (bool, error)

	// ListTradeDays 返回闭区间 [start, end] 内的交易日列表。
	// start 和 end 都必须为 2006-01-02 格式。
	// 当日期格式非法、日期区间非法或超出该日历支持范围时，返回 error。
	ListTradeDays(start, end Date) ([]Date, error)

	// OffsetTradeDay 返回给定日期前后第 N 个交易日。
	// day 必须为 2006-01-02 格式，且本身不必是交易日。
	// offset > 0 表示向后找第 N 个交易日。
	// offset < 0 表示向前找第 N 个交易日。
	// offset == 0 返回 error。
	OffsetTradeDay(day Date, offset int) (Date, error)

	// PrevTradeDay 返回给定日期的前一个交易日。
	// day 必须为 2006-01-02 格式，且本身不必是交易日。
	PrevTradeDay(day Date) (Date, error)

	// NextTradeDay 返回给定日期的后一个交易日。
	// day 必须为 2006-01-02 格式，且本身不必是交易日。
	NextTradeDay(day Date) (Date, error)
}

// New 根据日历标识创建对应的交易日历。
func New(calendarID CalendarID) (Calendar, error) {
	switch calendarID {
	case CalendarCNStock:
		return newCNStock(), nil
	case CalendarUSStock, "":
		return nil, newUnsupportedCalendarError()
	}

	return nil, newUnsupportedCalendarError()
}
