package core

import (
	"time"

	"github.com/Knowckx/tradeday/core/data"
)

const (
	cnStockMinYear = 2024
	cnStockMaxYear = 2026
)

var cnStockLocation = time.FixedZone("CST", 8*60*60)

// cnStock 表示中国 A 股市场日历。
type cnStock struct{}

// newCNStock 创建中国 A 股市场日历。
func newCNStock() *cnStock {
	return &cnStock{}
}

// IsTradeDay 判断给定日期是否为交易日。
func (c *cnStock) IsTradeDay(day Date) (bool, error) {
	normalizedDay, parsedDay, err := parseCNStockDate(day)
	if err != nil {
		return false, err
	}

	if !isSupportedCNStockDate(parsedDay) {
		return false, newDateOutOfRangeError()
	}

	if _, ok := data.CNStockOpenDays[string(normalizedDay)]; ok {
		return true, nil
	}

	if isWeekend(parsedDay) {
		return false, nil
	}

	if _, ok := data.CNStockClosedDays[string(normalizedDay)]; ok {
		return false, nil
	}

	return true, nil
}

// ListTradeDays 返回闭区间 [start, end] 内的交易日列表。
func (c *cnStock) ListTradeDays(start, end Date) ([]Date, error) {
	_, startDay, err := parseCNStockDate(start)
	if err != nil {
		return nil, err
	}

	_, endDay, err := parseCNStockDate(end)
	if err != nil {
		return nil, err
	}

	if startDay.After(endDay) {
		return nil, newInvalidDateRangeError()
	}

	if !isSupportedCNStockDate(startDay) || !isSupportedCNStockDate(endDay) {
		return nil, newDateOutOfRangeError()
	}

	tradeDays := make([]Date, 0)
	for current := startDay; !current.After(endDay); current = current.AddDate(0, 0, 1) {
		currentDate := Date(current.Format(DateLayout))
		isTradeDay, checkErr := c.IsTradeDay(currentDate)
		if checkErr != nil {
			return nil, checkErr
		}

		if isTradeDay {
			tradeDays = append(tradeDays, currentDate)
		}
	}

	return tradeDays, nil
}

// OffsetTradeDay 返回给定日期前后第 N 个交易日。
func (c *cnStock) OffsetTradeDay(day Date, offset int) (Date, error) {
	_, currentDay, err := parseCNStockDate(day)
	if err != nil {
		return "", err
	}

	if !isSupportedCNStockDate(currentDay) {
		return "", newDateOutOfRangeError()
	}

	if offset == 0 {
		return "", newInvalidOffsetError()
	}

	step := 1
	if offset < 0 {
		step = -1
		offset = -offset
	}

	for offset > 0 {
		currentDay = currentDay.AddDate(0, 0, step)
		if !isSupportedCNStockDate(currentDay) {
			return "", newDateOutOfRangeError()
		}

		currentDate := Date(currentDay.Format(DateLayout))
		isTradeDay, checkErr := c.IsTradeDay(currentDate)
		if checkErr != nil {
			return "", checkErr
		}

		if isTradeDay {
			offset--
		}
	}

	return Date(currentDay.Format(DateLayout)), nil
}

// PrevTradeDay 返回给定日期的前一个交易日。
func (c *cnStock) PrevTradeDay(day Date) (Date, error) {
	return c.OffsetTradeDay(day, -1)
}

// NextTradeDay 返回给定日期的后一个交易日。
func (c *cnStock) NextTradeDay(day Date) (Date, error) {
	return c.OffsetTradeDay(day, 1)
}

func parseCNStockDate(day Date) (Date, time.Time, error) {
	normalizedDay, err := day.Normalize()
	if err != nil {
		return "", time.Time{}, err
	}

	parsedDay, err := normalizedDay.parseInLocation(cnStockLocation)
	if err != nil {
		return "", time.Time{}, err
	}

	return normalizedDay, parsedDay, nil
}

func isWeekend(day time.Time) bool {
	weekday := day.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func isSupportedCNStockDate(day time.Time) bool {
	year := day.Year()
	return year >= cnStockMinYear && year <= cnStockMaxYear
}
