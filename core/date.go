package core

import "time"

const DateLayout = "2006-01-02"

// Date 表示一个日期字符串，固定格式为 2006-01-02。
type Date string

// ParseDate 解析一个 2006-01-02 格式的日期字符串。
func ParseDate(text string) (Date, error) {
	date := Date(text)
	if err := date.Validate(); err != nil {
		return "", err
	}

	return date, nil
}

// Validate 校验日期格式是否合法。
func (d Date) Validate() error {
	_, err := d.parseInLocation(time.UTC)
	return err
}

// Normalize 将日期规范化为标准的 2006-01-02 格式。
func (d Date) Normalize() (Date, error) {
	parsedDate, err := d.parseInLocation(time.UTC)
	if err != nil {
		return "", err
	}

	return Date(parsedDate.Format(DateLayout)), nil
}

func (d Date) parseInLocation(location *time.Location) (time.Time, error) {
	parsedDate, err := time.ParseInLocation(DateLayout, string(d), location)
	if err != nil {
		return time.Time{}, newInvalidDateFormatError()
	}

	return parsedDate, nil
}
