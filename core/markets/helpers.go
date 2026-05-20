package markets

import (
	"time"

	"github.com/Knowckx/tradeday/core/base"
)

// calendarDateFactory 用于从 time.Time 构造日历内部日期对象。
type calendarDateFactory func(day time.Time) (*base.CalendarDate, error)

// calendarTradeDayChecker 用于复用基于日历内部日期对象的交易日判断逻辑。
type calendarTradeDayChecker func(day *base.CalendarDate) (bool, error)

// listTradeDaysInRange 遍历闭区间，收集交易日。
func listTradeDaysInRange(
	startDay time.Time,
	endDay time.Time,
	buildCalendarDate calendarDateFactory,
	checkTradeDay calendarTradeDayChecker,
) ([]base.Date, error) {
	tradeDays := make([]base.Date, 0)
	for current := startDay; !current.After(endDay); current = current.AddDate(0, 0, 1) {
		currentDate, err := buildCalendarDate(current)
		if err != nil {
			return nil, err
		}

		isTradeDay, err := checkTradeDay(currentDate)
		if err != nil {
			return nil, err
		}

		if isTradeDay {
			tradeDays = append(tradeDays, currentDate.Date())
		}
	}

	return tradeDays, nil
}

// offsetTradeDayFromTime 从给定日期开始按交易日偏移。
func offsetTradeDayFromTime(
	currentDay time.Time,
	offset int,
	buildCalendarDate calendarDateFactory,
	checkTradeDay calendarTradeDayChecker,
) (base.Date, error) {
	if offset == 0 {
		return "", base.NewInvalidOffsetError()
	}

	step := 1
	if offset < 0 {
		step = -1
		offset = -offset
	}

	var currentDate *base.CalendarDate
	for offset > 0 {
		currentDay = currentDay.AddDate(0, 0, step)
		currentDate, err := buildCalendarDate(currentDay)
		if err != nil {
			return "", err
		}

		isTradeDay, err := checkTradeDay(currentDate)
		if err != nil {
			return "", err
		}

		if isTradeDay {
			offset--
		}
	}

	return currentDate.Date(), nil
}
