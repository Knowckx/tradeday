package markets

import (
	"time"
	_ "time/tzdata"

	"github.com/Knowckx/tradeday/core/base"
	"github.com/Knowckx/tradeday/core/data"
)

var usStockLocation = mustLoadLocation("America/New_York")
var usStockMinYear, usStockMaxYear = mustBitmapYearRange(data.USStockTradeBitmaps)

// usStock 表示美国美股市场日历。
type usStock struct{}

func mustLoadLocation(name string) *time.Location {
	location, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}

	return location
}

// USCalendar 创建美国美股市场日历。
func USCalendar() base.Calendar {
	return &usStock{}
}

// IsTradeDay 判断给定日期是否为交易日。
func (c *usStock) IsTradeDay(day string) (bool, error) {
	calendarDay, err := newMarketDate(day, usStockLocation, usStockMinYear, usStockMaxYear)
	if err != nil {
		return false, err
	}

	return c.isTradeDay(calendarDay)
}

func (c *usStock) isTradeDay(day *base.CalendarDate) (bool, error) {
	return data.USStockTradeBitmaps.IsTradeDay(day)
}

// PrevTradeDay 返回给定日期的前一个交易日。
func (c *usStock) PrevTradeDay(day string) (string, error) {
	return c.OffsetTradeDay(day, -1)
}

// NextTradeDay 返回给定日期的后一个交易日。
func (c *usStock) NextTradeDay(day string) (string, error) {
	return c.OffsetTradeDay(day, 1)
}

// OffsetTradeDay 返回交易日偏移结果。
func (c *usStock) OffsetTradeDay(day string, offset int) (string, error) {
	calendarDay, err := newMarketDate(day, usStockLocation, usStockMinYear, usStockMaxYear)
	if err != nil {
		return "", err
	}

	targetDay, err := data.USStockTradeBitmaps.OffsetTradeDay(calendarDay, offset)
	if err != nil {
		return "", err
	}

	return string(targetDay), nil
}

// ListTradeDays 返回闭区间 [start, end] 内的交易日列表。
func (c *usStock) ListTradeDays(start, end string) ([]string, error) {
	startDay, err := newMarketDate(start, usStockLocation, usStockMinYear, usStockMaxYear)
	if err != nil {
		return nil, err
	}

	endDay, err := newMarketDate(end, usStockLocation, usStockMinYear, usStockMaxYear)
	if err != nil {
		return nil, err
	}

	if startDay.Time().After(endDay.Time()) {
		return nil, base.NewInvalidDateRangeError()
	}

	tradeDays, err := data.USStockTradeBitmaps.ListTradeDays(startDay, endDay)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(tradeDays))
	for i, tradeDay := range tradeDays {
		result[i] = string(tradeDay)
	}

	return result, nil
}
