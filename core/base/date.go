package base

import "time"

const DateLayout = "2006-01-02"

// Date 表示一个日期字符串，固定格式为 2006-01-02。
type Date string

// ToCalendarDate 对日期做预检查并转换为内部日期对象。
func (day Date) ToCalendarDate(location *time.Location) (*CalendarDate, error) {
	if len(day) != len(DateLayout) {
		return nil, NewInvalidDateFormatError()
	}

	parsedDay, err := time.ParseInLocation(DateLayout, string(day), location)
	if err != nil {
		return nil, NewInvalidDateFormatError()
	}

	normalizedDay := Date(parsedDay.Format(DateLayout))
	if normalizedDay != day {
		return nil, NewInvalidDateFormatError()
	}

	return &CalendarDate{
		date: normalizedDay,
		time: parsedDay,
	}, nil
}
