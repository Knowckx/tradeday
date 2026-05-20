package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Knowckx/tradeday"
)

const (
	cnStockTruthTableStart = "2015-01-01"
	cnStockTruthTableEnd   = "2026-12-31"
	dateLayout             = "2006-01-02"
)

func TestCNStockIsTradeDayAgainstTruthTable(t *testing.T) {
	cal, err := tradeday.New(tradeday.CalendarID.CNStock)
	if err != nil {
		t.Fatalf("创建 A 股日历失败: %v", err)
	}

	truthTable := loadCNStockTruthTable(t)
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
		t.Run(string(day), func(t *testing.T) {
			_, err := cal.IsTradeDay(day)
			if err == nil {
				t.Fatalf("IsTradeDay(%q) 未返回错误", day)
			}
		})
	}
}

func loadCNStockTruthTable(t *testing.T) map[tradeday.Date]bool {
	t.Helper()

	type truthTableFile struct {
		CalendarID string          `json:"calendar_id"`
		Start      string          `json:"start"`
		End        string          `json:"end"`
		Days       map[string]bool `json:"days"`
	}

	filePath := mustGetTruthTablePath(t)
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("读取真值表失败: %v", err)
	}

	var file truthTableFile
	if err := json.Unmarshal(content, &file); err != nil {
		t.Fatalf("解析真值表失败: %v", err)
	}

	if file.CalendarID != tradeday.CalendarID.CNStock {
		t.Fatalf("真值表 calendar_id = %q, want %q", file.CalendarID, tradeday.CalendarID.CNStock)
	}

	if file.Start != cnStockTruthTableStart || file.End != cnStockTruthTableEnd {
		t.Fatalf("真值表范围 = [%s, %s], want [%s, %s]", file.Start, file.End, cnStockTruthTableStart, cnStockTruthTableEnd)
	}

	truthTable := make(map[tradeday.Date]bool, len(file.Days))
	for day, isTradeDay := range file.Days {
		truthTable[tradeday.Date(day)] = isTradeDay
	}

	return truthTable
}

func mustGetTruthTablePath(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("获取当前测试文件路径失败")
	}

	return filepath.Join(filepath.Dir(currentFile), "testdata", "cn_stock_truth_table.json")
}

func countDays(t *testing.T, start, end string) int {
	t.Helper()

	days := 0
	for range iterateDays(t, start, end) {
		days++
	}

	return days
}

func iterateDays(t *testing.T, start, end string) <-chan tradeday.Date {
	t.Helper()

	startDay, err := time.Parse(dateLayout, start)
	if err != nil {
		t.Fatalf("解析开始日期失败: %v", err)
	}

	endDay, err := time.Parse(dateLayout, end)
	if err != nil {
		t.Fatalf("解析结束日期失败: %v", err)
	}

	ch := make(chan tradeday.Date)
	go func() {
		defer close(ch)
		for day := startDay; !day.After(endDay); day = day.AddDate(0, 0, 1) {
			ch <- tradeday.Date(day.Format(dateLayout))
		}
	}()

	return ch
}
