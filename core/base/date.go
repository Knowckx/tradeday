package base

import "time"

const DateLayout = "2006-01-02"

// Date 表示一个日期字符串，固定格式为 2006-01-02。
type Date = string

// NormalizeDate 将日期规范化为标准的 2006-01-02 格式。
func NormalizeDate(day Date) (Date, error) {
	parsedDate, err := parseDateInLocation(day, time.UTC)
	if err != nil {
		return "", err
	}

	return Date(parsedDate.Format(DateLayout)), nil
}

func parseDateInLocation(day Date, location *time.Location) (time.Time, error) {
	parsedDate, err := time.ParseInLocation(DateLayout, string(day), location)
	if err != nil {
		return time.Time{}, NewInvalidDateFormatError()
	}

	return parsedDate, nil
}
