package markets

import (
	"time"

	"github.com/Knowckx/tradeday/core/base"
	"github.com/Knowckx/tradeday/core/data"
)

func newMarketDate(day base.Date, location *time.Location, minYear, maxYear int) (*base.CalendarDate, error) {
	calendarDay, err := day.ToCalendarDate(location)
	if err != nil {
		return nil, err
	}

	if !calendarDay.IsSupportedYear(minYear, maxYear) {
		return nil, base.NewDateOutOfRangeError()
	}

	return calendarDay, nil
}

func mustBitmapYearRange(bitmaps data.YearTradeBitmaps) (int, int) {
	minYear, maxYear, ok := bitmaps.SupportedYearRange()
	if !ok {
		panic("trade bitmap year range is empty")
	}

	return minYear, maxYear
}
