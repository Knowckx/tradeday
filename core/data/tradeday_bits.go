package data

// YearTradeBitmap 表示某一年的交易日真值位图。
// 6 个 uint64 共 384 bit，足够覆盖一年 366 天。
type YearTradeBitmap struct {
	Year int
	Bits [6]uint64
}
