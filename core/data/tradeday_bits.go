package data

import (
	"fmt"
	"math/bits"
	"time"

	"github.com/Knowckx/tradeday/core/base"
)

// YearTradeBitmap 表示某一年的交易日真值位图。
// 6 个 uint64 共 384 bit，足够覆盖一年 366 天。
// Bits 中 bit 0 表示当年 1 月 1 日，bit 1 表示 1 月 2 日，以此类推。
// bit=1 表示交易日，bit=0 表示非交易日。
type YearTradeBitmap struct {
	Year int
	Bits [6]uint64
}

// YearTradeBitmaps 表示按年份索引的交易日真值位图集合。
type YearTradeBitmaps map[int]YearTradeBitmap

// IsTradeDay 判断给定日期是否为交易日。
func (bitmaps YearTradeBitmaps) IsTradeDay(day *base.CalendarDate) (bool, error) {
	if day == nil {
		return false, base.NewDateOutOfRangeError()
	}

	t := day.Time()
	year := t.Year()
	dayIndex := t.YearDay() - 1

	bitmap, ok := bitmaps[year]
	if !ok {
		return false, base.NewDateOutOfRangeError()
	}

	if dayIndex < 0 || dayIndex >= len(bitmap.Bits)*64 {
		return false, base.NewDateOutOfRangeError()
	}

	word := bitmap.Bits[dayIndex/64]
	bit := uint(dayIndex % 64)
	return (word>>bit)&1 == 1, nil
}

// OffsetTradeDay 返回交易日偏移结果。
// offset > 0 时，返回 day 之后的第 offset 个交易日，不包含 day 当天。
// offset < 0 时，返回 day 之前的第 -offset 个交易日，不包含 day 当天。
// offset == 0 时，仅当 day 当天是交易日时返回 day，否则返回 error。
func (bitmaps YearTradeBitmaps) OffsetTradeDay(day *base.CalendarDate, offset int) (base.Date, error) {
	if day == nil {
		return "", base.NewDateOutOfRangeError()
	}

	if offset == 0 {
		isTradeDay, err := bitmaps.IsTradeDay(day)
		if err != nil {
			return "", err
		}

		if isTradeDay {
			return day.Date(), nil
		}

		return "", base.NewInvalidOffsetError()
	}

	t := day.Time()
	year := t.Year()
	dayIndex := t.YearDay() - 1
	location := t.Location()

	if offset > 0 {
		return bitmaps.offsetTradeDayByStep(year, dayIndex, offset, 1, location)
	}

	return bitmaps.offsetTradeDayByStep(year, dayIndex, -offset, -1, location)
}

// ListTradeDays 返回闭区间 [start, end] 内的交易日列表。
func (bitmaps YearTradeBitmaps) ListTradeDays(start, end *base.CalendarDate) ([]base.Date, error) {
	if start == nil || end == nil {
		return nil, base.NewDateOutOfRangeError()
	}

	startTime := start.Time()
	endTime := end.Time()
	if startTime.After(endTime) {
		return nil, base.NewInvalidDateRangeError()
	}
	if startTime.Location() != endTime.Location() {
		return nil, base.NewInvalidDateRangeError()
	}

	location := startTime.Location()
	startYear := startTime.Year()
	endYear := endTime.Year()

	tradeDays := make([]base.Date, 0, (endYear-startYear+1)*250)
	for year := startYear; year <= endYear; year++ {
		bitmap, ok := bitmaps[year]
		if !ok {
			return nil, base.NewDateOutOfRangeError()
		}

		daysInYear := tradeYearDayCount(year, location)
		startBit := 0
		if year == startYear {
			startBit = startTime.YearDay() - 1
		}

		endBit := daysInYear - 1
		if year == endYear {
			endBit = endTime.YearDay() - 1
		}

		tradeBitmapRange(bitmap.Bits, startBit, endBit, false, func(wordIndex int, word uint64) bool {
			for word != 0 {
				trailing := bits.TrailingZeros64(word)
				tradeDays = append(tradeDays, tradeDateFromYearDay(year, wordIndex*64+trailing, location))
				word &= word - 1
			}

			return true
		})
	}

	return tradeDays, nil
}

func (bitmaps YearTradeBitmaps) offsetTradeDayByStep(year, dayIndex, offset, step int, location *time.Location) (base.Date, error) {
	for {
		bitmap, ok := bitmaps[year]
		if !ok {
			return "", base.NewDateOutOfRangeError()
		}

		daysInYear := tradeYearDayCount(year, location)
		available, targetIndex, ok := tradeDayOffsetWithinYear(bitmap.Bits, dayIndex, offset, step, daysInYear)
		if ok {
			return tradeDateFromYearDay(year, targetIndex, location), nil
		}

		offset -= available
		year += step
		if year < 1 {
			return "", base.NewDateOutOfRangeError()
		}

		if step > 0 {
			dayIndex = -1
			continue
		}

		dayIndex = tradeYearDayCount(year, location)
	}
}

func tradeDayOffsetWithinYear(bitmap [6]uint64, dayIndex, offset, step, daysInYear int) (int, int, bool) {
	if step > 0 {
		available := countTradeDaysAfter(bitmap, dayIndex, daysInYear)
		if offset <= available {
			targetIndex, ok := findNthTradeDayAfter(bitmap, dayIndex, offset, daysInYear)
			if !ok {
				return 0, 0, false
			}

			return available, targetIndex, true
		}

		return available, 0, false
	}

	available := countTradeDaysBefore(bitmap, dayIndex, daysInYear)
	if offset <= available {
		targetIndex, ok := findNthTradeDayBefore(bitmap, dayIndex, offset, daysInYear)
		if !ok {
			return 0, 0, false
		}

		return available, targetIndex, true
	}

	return available, 0, false
}

func countTradeDaysAfter(bitmap [6]uint64, dayIndex, daysInYear int) int {
	count := 0
	tradeBitmapRange(bitmap, dayIndex+1, daysInYear-1, false, func(_ int, word uint64) bool {
		count += bits.OnesCount64(word)
		return true
	})

	return count
}

func countTradeDaysBefore(bitmap [6]uint64, dayIndex, daysInYear int) int {
	count := 0
	tradeBitmapRange(bitmap, 0, dayIndex-1, true, func(_ int, word uint64) bool {
		count += bits.OnesCount64(word)
		return true
	})

	return count
}

func findNthTradeDayAfter(bitmap [6]uint64, dayIndex, nth, daysInYear int) (int, bool) {
	if nth <= 0 {
		return 0, false
	}

	found := false
	targetIndex := 0
	tradeBitmapRange(bitmap, dayIndex+1, daysInYear-1, false, func(wordIndex int, word uint64) bool {
		for word != 0 {
			trailing := bits.TrailingZeros64(word)
			nth--
			if nth == 0 {
				targetIndex = wordIndex*64 + trailing
				found = true
				return false
			}

			word &= word - 1
		}

		return true
	})

	return targetIndex, found
}

func findNthTradeDayBefore(bitmap [6]uint64, dayIndex, nth, daysInYear int) (int, bool) {
	if nth <= 0 {
		return 0, false
	}

	found := false
	targetIndex := 0
	tradeBitmapRange(bitmap, 0, dayIndex-1, true, func(wordIndex int, word uint64) bool {
		for word != 0 {
			leading := 63 - bits.LeadingZeros64(word)
			nth--
			if nth == 0 {
				targetIndex = wordIndex*64 + leading
				found = true
				return false
			}

			word &= ^(uint64(1) << uint(leading))
		}

		return true
	})

	return targetIndex, found
}

func tradeBitmapRange(bitmap [6]uint64, startBit, endBit int, reverse bool, visit func(wordIndex int, word uint64) bool) {
	if startBit > endBit {
		return
	}

	maxBit := len(bitmap)*64 - 1
	if endBit < 0 || startBit > maxBit {
		return
	}

	if startBit < 0 {
		startBit = 0
	}

	if endBit > maxBit {
		endBit = maxBit
	}

	firstWord := startBit / 64
	lastWord := endBit / 64
	if reverse {
		for wordIndex := lastWord; wordIndex >= firstWord; wordIndex-- {
			word := bitmap[wordIndex] & tradeBitmapWordMask(wordIndex, startBit, endBit)
			if word != 0 && !visit(wordIndex, word) {
				return
			}

			if wordIndex == 0 {
				break
			}
		}

		return
	}

	for wordIndex := firstWord; wordIndex <= lastWord; wordIndex++ {
		word := bitmap[wordIndex] & tradeBitmapWordMask(wordIndex, startBit, endBit)
		if word != 0 && !visit(wordIndex, word) {
			return
		}
	}
}

func tradeBitmapWordMask(wordIndex, startBit, endBit int) uint64 {
	wordStart := wordIndex * 64
	wordEnd := wordStart + 63
	if wordEnd < startBit || wordStart > endBit {
		return 0
	}

	mask := ^uint64(0)
	if wordStart < startBit {
		mask &= ^uint64(0) << uint(startBit-wordStart)
	}

	if wordEnd > endBit {
		keepBits := endBit - wordStart + 1
		if keepBits <= 0 {
			return 0
		}
		if keepBits < 64 {
			mask &= (uint64(1) << uint(keepBits)) - 1
		}
	}

	return mask
}

func tradeYearDayCount(year int, location *time.Location) int {
	return time.Date(year, time.December, 31, 0, 0, 0, 0, location).YearDay()
}

func tradeDateFromYearDay(year, dayIndex int, location *time.Location) base.Date {
	return base.Date(time.Date(year, time.January, 1+dayIndex, 0, 0, 0, 0, location).Format(base.DateLayout))
}

func (bitmaps YearTradeBitmaps) mustAlignYearKeys() {
	for year, bitmap := range bitmaps {
		if bitmap.Year != year {
			panic(fmt.Sprintf("trade bitmap year mismatch: key=%d bitmap.Year=%d", year, bitmap.Year))
		}
	}
}
