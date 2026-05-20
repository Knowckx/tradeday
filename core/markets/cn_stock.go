package markets

import (
	"time"

	"github.com/Knowckx/tradeday/core/base"
	"github.com/Knowckx/tradeday/core/data"
)

const (
	cnStockMinYear           = 2015
	cnStockMaxYear           = 2026
	cnStockMaxTradeDaysPerYear = 260
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
	if _, ok := data.CNStockOpenDays[day.Date()]; ok {
		return true, nil
	}

	if day.IsWeekend() {
		return false, nil
	}

	if _, ok := data.CNStockClosedDays[day.Date()]; ok {
		return false, nil
	}

	return true, nil
}

// PrevTradeDay 返回给定日期的前一个交易日。
func (c *cnStock) PrevTradeDay(day base.Date) (base.Date, error) {
	return c.OffsetTradeDay(day, -1)
}

// NextTradeDay 返回给定日期的后一个交易日。
func (c *cnStock) NextTradeDay(day base.Date) (base.Date, error) {
	return c.OffsetTradeDay(day, 1)
}

// OffsetTradeDay 返回给定日期前后第 N 个交易日。
func (c *cnStock) OffsetTradeDay(day base.Date, offset int) (base.Date, error) {
	calendarDay, err := newCNStockDate(day)
	if err != nil {
		return "", err
	}

	if offset == 0 {
		return "", base.NewInvalidOffsetError()
	}

	if !canOffsetWithinSupportRange(calendarDay.Time(), offset) {
		return "", base.NewDateOutOfRangeError()
	}

	currentDate := calendarDay
	step := 1
	if offset < 0 {
		step = -1
		offset = -offset
	}

	var currentTradeDate *base.CalendarDate
	for offset > 0 {
		currentDate = currentDate.AddDays(step)
		currentTradeDate, err = newCNStockDateFromTime(currentDate.Time())
		if err != nil {
			return "", err
		}

		isTradeDay, err := c.isTradeDay(currentTradeDate)
		if err != nil {
			return "", err
		}

		if isTradeDay {
			offset--
		}
	}

	return currentTradeDate.Date(), nil
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

	return listTradeDaysInRange(
		startDay.Time(),
		endDay.Time(),
		newCNStockDateFromTime,
		c.isTradeDay,
	)
}

// newCNStockDateFromTime 为内部循环场景构造 A 股日期对象。
func newCNStockDateFromTime(day time.Time) (*base.CalendarDate, error) {
	cnStockDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, cnStockLocation)
	calendarDay := base.NewCalendarDateFromTime(cnStockDay)
	if !calendarDay.IsSupportedYear(cnStockMinYear, cnStockMaxYear) {
		return nil, base.NewDateOutOfRangeError()
	}

	return calendarDay, nil
}

func canOffsetWithinSupportRange(day time.Time, offset int) bool {
	if offset > 0 {
		remainingYears := cnStockMaxYear - day.Year() + 1
		return offset <= remainingYears*cnStockMaxTradeDaysPerYear
	}

	remainingYears := day.Year() - cnStockMinYear + 1
	return -offset <= remainingYears*cnStockMaxTradeDaysPerYear
}
