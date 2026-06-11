package base

// Calendar provides trading-day queries for a specific market.
// All date inputs must use the format "2006-01-02", for example "2020-01-01".
type Calendar interface {
	// IsTradeDay reports whether day is a trading day.
	// day must use the format "2006-01-02", for example "2020-01-01".
	IsTradeDay(day string) (bool, error)

	// PrevTradeDay returns the trading day immediately before day.
	// day must use the format "2006-01-02", for example "2020-01-01".
	PrevTradeDay(day string) (string, error)

	// NextTradeDay returns the trading day immediately after day.
	// day must use the format "2006-01-02", for example "2020-01-01".
	NextTradeDay(day string) (string, error)

	// OffsetTradeDay returns the trading day offset from day by offset trading days.
	// day must use the format "2006-01-02", for example "2020-01-01".
	OffsetTradeDay(day string, offset int) (string, error)

	// ListTradeDays returns all trading days in the inclusive range [start, end].
	// start and end must use the format "2006-01-02", for example "2020-01-01".
	ListTradeDays(start, end string) ([]string, error)
}
