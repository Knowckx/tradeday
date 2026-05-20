package tests

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Knowckx/tradeday"
)

func TestCNStockPrevNextTradeDayAgainstTruthTable(t *testing.T) {
	cal, err := tradeday.New(tradeday.CalendarID.CNStock)
	if err != nil {
		t.Fatalf("创建 A 股日历失败: %v", err)
	}

	truthTable := loadCNStockTruthTable(t)

	for day := range iterateDays(t, cnStockTruthTableStart, cnStockTruthTableEnd) {
		wantPrev, ok := findRelativeTradeDayFromTruthTable(truthTable, day, -1)
		gotPrev, err := cal.PrevTradeDay(day)
		if ok {
			if err != nil {
				t.Fatalf("PrevTradeDay(%q) 返回错误: %v", day, err)
			}
			if gotPrev != wantPrev {
				t.Fatalf("PrevTradeDay(%q) = %q, want %q", day, gotPrev, wantPrev)
			}
			offsetPrev, err := cal.OffsetTradeDay(day, -1)
			if err != nil {
				t.Fatalf("OffsetTradeDay(%q, -1) 返回错误: %v", day, err)
			}
			if gotPrev != offsetPrev {
				t.Fatalf("PrevTradeDay(%q) = %q, OffsetTradeDay(%q, -1) = %q", day, gotPrev, day, offsetPrev)
			}
		} else {
			assertErrorIs(t, err, tradeday.Error("date_out_of_range"))
			offsetPrev, err := cal.OffsetTradeDay(day, -1)
			assertErrorIs(t, err, tradeday.Error("date_out_of_range"))
			if offsetPrev != "" {
				t.Fatalf("OffsetTradeDay(%q, -1) = %q, want 空值", day, offsetPrev)
			}
		}

		wantNext, ok := findRelativeTradeDayFromTruthTable(truthTable, day, 1)
		gotNext, err := cal.NextTradeDay(day)
		if ok {
			if err != nil {
				t.Fatalf("NextTradeDay(%q) 返回错误: %v", day, err)
			}
			if gotNext != wantNext {
				t.Fatalf("NextTradeDay(%q) = %q, want %q", day, gotNext, wantNext)
			}
			offsetNext, err := cal.OffsetTradeDay(day, 1)
			if err != nil {
				t.Fatalf("OffsetTradeDay(%q, 1) 返回错误: %v", day, err)
			}
			if gotNext != offsetNext {
				t.Fatalf("NextTradeDay(%q) = %q, OffsetTradeDay(%q, 1) = %q", day, gotNext, day, offsetNext)
			}
		} else {
			assertErrorIs(t, err, tradeday.Error("date_out_of_range"))
			offsetNext, err := cal.OffsetTradeDay(day, 1)
			assertErrorIs(t, err, tradeday.Error("date_out_of_range"))
			if offsetNext != "" {
				t.Fatalf("OffsetTradeDay(%q, 1) = %q, want 空值", day, offsetNext)
			}
		}
	}
}

func findRelativeTradeDayFromTruthTable(
	truthTable map[tradeday.Date]bool,
	day tradeday.Date,
	direction int,
) (tradeday.Date, bool) {
	currentDay, err := time.Parse(dateLayout, string(day))
	if err != nil {
		return "", false
	}

	startDay, err := time.Parse(dateLayout, cnStockTruthTableStart)
	if err != nil {
		return "", false
	}

	endDay, err := time.Parse(dateLayout, cnStockTruthTableEnd)
	if err != nil {
		return "", false
	}

	if direction < 0 {
		for currentDay = currentDay.AddDate(0, 0, -1); !currentDay.Before(startDay); currentDay = currentDay.AddDate(0, 0, -1) {
			candidate := tradeday.Date(currentDay.Format(dateLayout))
			if truthTable[candidate] {
				return candidate, true
			}
		}

		return "", false
	}

	for currentDay = currentDay.AddDate(0, 0, 1); !currentDay.After(endDay); currentDay = currentDay.AddDate(0, 0, 1) {
		candidate := tradeday.Date(currentDay.Format(dateLayout))
		if truthTable[candidate] {
			return candidate, true
		}
	}

	return "", false
}

func TestCNStockPrevNextTradeDayInvalidInput(t *testing.T) {
	cal, err := tradeday.New(tradeday.CalendarID.CNStock)
	if err != nil {
		t.Fatalf("创建 A 股日历失败: %v", err)
	}

	testCases := []tradeday.Date{
		"2024-1-02",
		"2024/01/02",
		"2024-01-02 ",
		"2027-01-01",
	}

	for _, day := range testCases {
		t.Run(fmt.Sprintf("prev_%s", day), func(t *testing.T) {
			_, err := cal.PrevTradeDay(day)
			assertErrorIs(t, err, wantErrorForInvalidInput(day))
		})

		t.Run(fmt.Sprintf("next_%s", day), func(t *testing.T) {
			_, err := cal.NextTradeDay(day)
			assertErrorIs(t, err, wantErrorForInvalidInput(day))
		})
	}
}

func TestCNStockOffsetTradeDayLargeOffset(t *testing.T) {
	cal, err := tradeday.New(tradeday.CalendarID.CNStock)
	if err != nil {
		t.Fatalf("创建 A 股日历失败: %v", err)
	}

	t.Run("positive", func(t *testing.T) {
		_, err := cal.OffsetTradeDay("2024-01-02", 1000)
		assertErrorIs(t, err, tradeday.Error("date_out_of_range"))
	})

	t.Run("negative", func(t *testing.T) {
		_, err := cal.OffsetTradeDay("2015-01-02", -1000)
		assertErrorIs(t, err, tradeday.Error("date_out_of_range"))
	})
}

func wantErrorForInvalidInput(day tradeday.Date) error {
	if day == "2027-01-01" {
		return tradeday.Error("date_out_of_range")
	}

	return tradeday.Error("invalid_date_format")
}

func assertErrorIs(t *testing.T, err error, want error) {
	t.Helper()

	if err == nil {
		t.Fatalf("未返回错误，want %v", want)
	}

	if !errors.Is(err, want) {
		t.Fatalf("错误 = %v, want %v", err, want)
	}
}
