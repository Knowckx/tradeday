package base

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
