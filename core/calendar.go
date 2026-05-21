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
		return markets.NewUSStock(), nil
	}

	return nil, base.NewUnsupportedCalendarError()
}
