package tests

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/Knowckx/tradeday"
)

const dateLayout = "2006-01-02"

func loadTruthTable(t *testing.T, fileName, expectedCalendarID, expectedStart, expectedEnd string) map[tradeday.Date]bool {
	t.Helper()

	type truthTableFile struct {
		CalendarID string          `json:"calendar_id"`
		Start      string          `json:"start"`
		End        string          `json:"end"`
		Days       map[string]bool `json:"days"`
	}

	filePath := mustGetTruthTablePath(t, fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("读取真值表失败: %v", err)
	}

	var file truthTableFile
	if err := json.Unmarshal(content, &file); err != nil {
		t.Fatalf("解析真值表失败: %v", err)
	}

	if file.CalendarID != expectedCalendarID {
		t.Fatalf("真值表 calendar_id = %q, want %q", file.CalendarID, expectedCalendarID)
	}

	if file.Start != expectedStart || file.End != expectedEnd {
		t.Fatalf(
			"真值表范围 = [%s, %s], want [%s, %s]",
			file.Start,
			file.End,
			expectedStart,
			expectedEnd,
		)
	}

	truthTable := make(map[tradeday.Date]bool, len(file.Days))
	for day, isTradeDay := range file.Days {
		truthTable[tradeday.Date(day)] = isTradeDay
	}

	return truthTable
}

func mustGetTruthTablePath(t *testing.T, fileName string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("获取当前测试文件路径失败")
	}

	return filepath.Join(filepath.Dir(currentFile), "testdata", fileName)
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

func findRelativeTradeDayFromTruthTable(
	truthTable map[tradeday.Date]bool,
	day tradeday.Date,
	direction int,
	start,
	end string,
) (tradeday.Date, bool) {
	currentDay, err := time.Parse(dateLayout, string(day))
	if err != nil {
		return "", false
	}

	startDay, err := time.Parse(dateLayout, start)
	if err != nil {
		return "", false
	}

	endDay, err := time.Parse(dateLayout, end)
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
