package tests

import (
	"fmt"
	"testing"

	"github.com/Knowckx/tradeday"
)

func TestUSStockPrevNextTradeDayAgainstTruthTable(t *testing.T) {
	cal, err := tradeday.New(tradeday.CalendarID.USStock)
	if err != nil {
		t.Fatalf("创建美股日历失败: %v", err)
	}

	truthTable := loadTruthTable(t, "us_stock_truth_table.json", tradeday.CalendarID.USStock, usStockTruthTableStart, usStockTruthTableEnd)

	for day := range iterateDays(t, usStockTruthTableStart, usStockTruthTableEnd) {
		wantPrev, ok := findRelativeTradeDayFromTruthTable(truthTable, day, -1, usStockTruthTableStart, usStockTruthTableEnd)
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

		wantNext, ok := findRelativeTradeDayFromTruthTable(truthTable, day, 1, usStockTruthTableStart, usStockTruthTableEnd)
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

func TestUSStockPrevNextTradeDayInvalidInput(t *testing.T) {
	cal, err := tradeday.New(tradeday.CalendarID.USStock)
	if err != nil {
		t.Fatalf("创建美股日历失败: %v", err)
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

func TestUSStockOffsetTradeDayLargeOffset(t *testing.T) {
	cal, err := tradeday.New(tradeday.CalendarID.USStock)
	if err != nil {
		t.Fatalf("创建美股日历失败: %v", err)
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

func TestUSStockOffsetTradeDayZero(t *testing.T) {
	cal, err := tradeday.New(tradeday.CalendarID.USStock)
	if err != nil {
		t.Fatalf("创建美股日历失败: %v", err)
	}

	t.Run("trade_day", func(t *testing.T) {
		got, err := cal.OffsetTradeDay("2024-10-08", 0)
		if err != nil {
			t.Fatalf("OffsetTradeDay(%q, 0) 返回错误: %v", "2024-10-08", err)
		}
		if got != "2024-10-08" {
			t.Fatalf("OffsetTradeDay(%q, 0) = %q, want %q", "2024-10-08", got, "2024-10-08")
		}
	})

	t.Run("non_trade_day", func(t *testing.T) {
		got, err := cal.OffsetTradeDay("2024-07-04", 0)
		assertErrorIs(t, err, tradeday.Error("invalid_offset"))
		if got != "" {
			t.Fatalf("OffsetTradeDay(%q, 0) = %q, want 空值", "2024-07-04", got)
		}
	})
}

func TestUSStockListTradeDays(t *testing.T) {
	cal, err := tradeday.New(tradeday.CalendarID.USStock)
	if err != nil {
		t.Fatalf("创建美股日历失败: %v", err)
	}

	t.Run("single_trade_day", func(t *testing.T) {
		got, err := cal.ListTradeDays("2024-07-04", "2024-07-05")
		if err != nil {
			t.Fatalf("ListTradeDays 返回错误: %v", err)
		}

		want := []tradeday.Date{"2024-07-05"}
		if len(got) != len(want) {
			t.Fatalf("ListTradeDays 长度 = %d, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("ListTradeDays[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("empty_range", func(t *testing.T) {
		got, err := cal.ListTradeDays("2024-07-04", "2024-07-04")
		if err != nil {
			t.Fatalf("ListTradeDays 返回错误: %v", err)
		}

		if len(got) != 0 {
			t.Fatalf("ListTradeDays = %v, want 空列表", got)
		}
	})
}
