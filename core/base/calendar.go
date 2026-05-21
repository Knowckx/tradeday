package base

// Calendar 表示某一类市场的交易日历。
// 所有日期参数都必须使用 2006-01-02 格式。
type Calendar interface {
	// IsTradeDay 判断给定日期是否为交易日。
	// day 进入实现后会先被规范化。
	// 当日期格式非法或超出该日历支持范围时，返回 error。
	IsTradeDay(day Date) (bool, error)

	// PrevTradeDay 返回给定日期的前一个交易日。
	// day 进入实现后会先被规范化，且本身不必是交易日。
	PrevTradeDay(day Date) (Date, error)

	// NextTradeDay 返回给定日期的后一个交易日。
	// day 进入实现后会先被规范化，且本身不必是交易日。
	NextTradeDay(day Date) (Date, error)

	// OffsetTradeDay 返回交易日偏移结果。
	// offset > 0 时，返回 day 之后的第 offset 个交易日，不包含 day 当天。
	// offset < 0 时，返回 day 之前的第 -offset 个交易日，不包含 day 当天。
	// offset == 0 时，仅当 day 当天是交易日时返回 day，否则返回 error。
	// day 进入实现后会先被规范化，且本身不必是交易日。
	OffsetTradeDay(day Date, offset int) (Date, error)

	// ListTradeDays 返回闭区间 [start, end] 内的交易日列表。
	// start 和 end 进入实现后都会先被规范化。
	// start 和 end 的 Location 必须一致。
	// 当日期格式非法、日期区间非法或超出该日历支持范围时，返回 error。
	ListTradeDays(start, end Date) ([]Date, error)
}
