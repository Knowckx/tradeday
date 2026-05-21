package markets

import (
	"time"

	"github.com/Knowckx/tradeday/core/base"
	"github.com/Knowckx/tradeday/core/data"
)

var cnStockLocation = time.FixedZone("CST", 8*60*60)
var cnStockMinYear, cnStockMaxYear = mustBitmapYearRange(data.CNStockTradeBitmaps)

// cnStock 表示中国 A 股市场日历。
type cnStock struct{}

// NewCNStock 创建中国 A 股市场日历。
func NewCNStock() base.Calendar {
	return &cnStock{}
}

// IsTradeDay 判断给定日期是否为交易日。
func (c *cnStock) IsTradeDay(day base.Date) (bool, error) {
	calendarDay, err := newMarketDate(day, cnStockLocation, cnStockMinYear, cnStockMaxYear)
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
	calendarDay, err := newMarketDate(day, cnStockLocation, cnStockMinYear, cnStockMaxYear)
	if err != nil {
		return "", err
	}

	return data.CNStockTradeBitmaps.OffsetTradeDay(calendarDay, offset)
}

// ListTradeDays 返回闭区间 [start, end] 内的交易日列表。
func (c *cnStock) ListTradeDays(start, end base.Date) ([]base.Date, error) {
	startDay, err := newMarketDate(start, cnStockLocation, cnStockMinYear, cnStockMaxYear)
	if err != nil {
		return nil, err
	}

	endDay, err := newMarketDate(end, cnStockLocation, cnStockMinYear, cnStockMaxYear)
	if err != nil {
		return nil, err
	}

	if startDay.Time().After(endDay.Time()) {
		return nil, base.NewInvalidDateRangeError()
	}

	return data.CNStockTradeBitmaps.ListTradeDays(startDay, endDay)
}
