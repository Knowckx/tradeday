package markets

import (
	"time"

	"github.com/Knowckx/tradeday/core/base"
	"github.com/Knowckx/tradeday/core/data"
)

const (
	cnStockMinYear = 2015
	cnStockMaxYear = 2026
)

var cnStockLocation = time.FixedZone("CST", 8*60*60)

// cnStock 表示中国 A 股市场日历。
type cnStock struct{}

// NewCNStock 创建中国 A 股市场日历。
func NewCNStock() base.Calendar {
	return &cnStock{}
}

// newCNStockDate 创建 A 股市场可接受的日期对象。
func newCNStockDate(day base.Date) (*base.CalendarDate, error) {
	calendarDay, err := day.ToCalendarDate(cnStockLocation)
	if err != nil {
		return nil, err
	}

	if !calendarDay.IsSupportedYear(cnStockMinYear, cnStockMaxYear) {
		return nil, base.NewDateOutOfRangeError()
	}

	return calendarDay, nil
}

// IsTradeDay 判断给定日期是否为交易日。
func (c *cnStock) IsTradeDay(day base.Date) (bool, error) {
	calendarDay, err := newCNStockDate(day)
	if err != nil {
		return false, err
	}

	return c.isTradeDay(calendarDay)
}

func (c *cnStock) isTradeDay(day *base.CalendarDate) (bool, error) {
	return data.CNStockTradeBitmaps.IsTradeDay(day)
}

// PrevTradeDay 返回给定日期的前一个交易日。
func (c *cnStock) PrevTradeDay(day base.Date) (base.Date, error) {
	return c.OffsetTradeDay(day, -1)
}

// NextTradeDay 返回给定日期的后一个交易日。
func (c *cnStock) NextTradeDay(day base.Date) (base.Date, error) {
	return c.OffsetTradeDay(day, 1)
}

// OffsetTradeDay 返回交易日偏移结果。
// offset > 0 时，返回 day 之后的第 offset 个交易日，不包含 day 当天。
// offset < 0 时，返回 day 之前的第 -offset 个交易日，不包含 day 当天。
// offset == 0 时，仅当 day 当天是交易日时返回 day，否则返回 error。
func (c *cnStock) OffsetTradeDay(day base.Date, offset int) (base.Date, error) {
	calendarDay, err := newCNStockDate(day)
	if err != nil {
		return "", err
	}

	return data.CNStockTradeBitmaps.OffsetTradeDay(calendarDay, offset)
}

// ListTradeDays 返回闭区间 [start, end] 内的交易日列表。
func (c *cnStock) ListTradeDays(start, end base.Date) ([]base.Date, error) {
	startDay, err := newCNStockDate(start)
	if err != nil {
		return nil, err
	}

	endDay, err := newCNStockDate(end)
	if err != nil {
		return nil, err
	}

	if startDay.Time().After(endDay.Time()) {
		return nil, base.NewInvalidDateRangeError()
	}

	return data.CNStockTradeBitmaps.ListTradeDays(startDay, endDay)
}
