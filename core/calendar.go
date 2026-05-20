package core

import (
	"github.com/Knowckx/tradeday/core/base"
	"github.com/Knowckx/tradeday/core/markets"
)

var CalendarID = struct {
	CNStock string
	USStock string
}{
	CNStock: "cn_stock",
	USStock: "us_stock",
}

// New 根据日历标识创建对应的交易日历。
func New(calendarID string) (base.Calendar, error) {
	switch calendarID {
	case CalendarID.CNStock:
		return markets.NewCNStock(), nil
	case CalendarID.USStock:
		return usStockPlaceholder{}, nil
	}

	return nil, base.NewUnsupportedCalendarError()
}

// usStockPlaceholder 是美国交易日日历的占位实现。
type usStockPlaceholder struct{}

func (usStockPlaceholder) IsTradeDay(base.Date) (bool, error) {
	return false, base.NewUnsupportedCalendarError()
}

func (usStockPlaceholder) PrevTradeDay(base.Date) (base.Date, error) {
	return "", base.NewUnsupportedCalendarError()
}

func (usStockPlaceholder) NextTradeDay(base.Date) (base.Date, error) {
	return "", base.NewUnsupportedCalendarError()
}

func (usStockPlaceholder) OffsetTradeDay(base.Date, int) (base.Date, error) {
	return "", base.NewUnsupportedCalendarError()
}

func (usStockPlaceholder) ListTradeDays(base.Date, base.Date) ([]base.Date, error) {
	return nil, base.NewUnsupportedCalendarError()
}
