package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Knowckx/tradeday"
)

const (
	cnStockCalendarID      = "cn_stock"
	cnStockTruthTableStart = "2015-01-01"
	cnStockTruthTableEnd   = "2026-12-31"
)

func TestCNStockIsTradeDayAgainstTruthTable(t *testing.T) {
	cal := tradeday.CNCalendar()
	truthTable := loadTruthTable(t, "cn_stock_truth_table.json", cnStockCalendarID, cnStockTruthTableStart, cnStockTruthTableEnd)
	totalDays := countDays(t, cnStockTruthTableStart, cnStockTruthTableEnd)
	if len(truthTable) != totalDays {
		t.Fatalf("真值表条数 = %d, want %d", len(truthTable), totalDays)
	}

	mismatchCount := 0
	mismatchSamples := make([]string, 0, 20)
	for day := range iterateDays(t, cnStockTruthTableStart, cnStockTruthTableEnd) {
		want, ok := truthTable[day]
		if !ok {
			t.Fatalf("真值表缺少日期 %s", day)
		}

		got, err := cal.IsTradeDay(day)
		if err != nil {
			t.Fatalf("IsTradeDay(%q) 返回错误: %v", day, err)
		}

		if got == want {
			continue
		}

		mismatchCount++
		if len(mismatchSamples) < 20 {
			mismatchSamples = append(mismatchSamples, fmt.Sprintf("%s got=%v want=%v", day, got, want))
		}
	}

	if mismatchCount > 0 {
		t.Fatalf(
			"逐日比对失败: %d/%d 不一致, 前%d个差异: %s",
			mismatchCount,
			totalDays,
			len(mismatchSamples),
			strings.Join(mismatchSamples, "; "),
		)
	}
}

func TestCNStockIsTradeDayInvalidInput(t *testing.T) {
	cal := tradeday.CNCalendar()
	testCases := []string{
		"2024-1-02",
		"2024/01/02",
		"2024-01-02 ",
		"2027-01-01",
	}

	for _, day := range testCases {
		t.Run(day, func(t *testing.T) {
			_, err := cal.IsTradeDay(day)
			if err == nil {
				t.Fatalf("IsTradeDay(%q) 未返回错误", day)
			}
		})
	}
}
