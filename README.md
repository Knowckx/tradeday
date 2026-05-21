# tradeday

<div align="center">

## 🌐 Choose your language / 选择语言

| 🌐 Language | 选择语言 |
|--------------|----------|
| [English](#english) | [中文](#中文) |

</div>

## 中文

`tradeday` 是一个可离线查询的 Go 交易日历库，用于判断指定日期是否为交易日。  
它内置交易日数据，适合量化、交易系统、回测、行情采集和定时任务调度等场景。

当前仓库已支持两个市场，支持范围分别为：

- `CNStock`：`2015-01-01` 到 `2026-12-31`
- `USStock`：`2015-01-01` 到 `2026-12-31`

项目底层使用`交易日位图`来保存每年的交易日数据，具有极佳的查询效率
- 单个交易日查询实现`O(1)`效率
- 范围交易日查询实现`O(n)`效率。

## 特性

- 支持按市场创建交易日历
- 支持判断某一天是否为交易日
- 支持前一个交易日、后一个交易日、按偏移取交易日
- 支持闭区间交易日列表查询
- 交易日判断基于本地位图数据
- 支持周末、法定节假日和特殊休市日的真实交易日结果
- 人性化的 API 设计，适合直接嵌入 Go 项目
- 无需网络请求即可使用

## 安装

```bash
go get github.com/Knowckx/tradeday
```

## 快速开始

`Date` 是一个以 `string` 为底层类型的日期类型，因此可以直接传入形如 `"2024-10-08"` 的字符串字面量。

```go
package main

import (
	"fmt"
	"log"

	"github.com/Knowckx/tradeday"
)

func main() {
	cal, err := tradeday.New(tradeday.CalendarID.CNStock)
	if err != nil {
		log.Fatal(err)
	}

	ok, err := cal.IsTradeDay("2024-10-08")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ok)
}
```

## 核心概念

- `Date`
  - 对外统一使用的日期类型，固定格式为 `2006-01-02`
- `Calendar`
  - 某一类市场的交易日历实例
- `CalendarID`
  - 用于选择具体市场，目前包含 `CNStock` 和 `USStock`
- `交易日位图`
  - 底层实现，每个年份对应一份位图数据，bit=1 表示交易日，bit=0 表示非交易日

## API 用法

### 创建日历

```go
cal, err := tradeday.New(tradeday.CalendarID.USStock)
if err != nil {
	return
}
```

### 判断交易日

```go
ok, err := cal.IsTradeDay("2024-10-08")
if err != nil {
	return
}
_ = ok
```

### 前后交易日

```go
prev, err := cal.PrevTradeDay("2024-10-08")
next, err := cal.NextTradeDay("2024-10-08")
_ = prev
_ = next
_ = err
```

### 偏移交易日

```go
// 返回"2024-10-08"之后的第 5 个有效交易日
day, err := cal.OffsetTradeDay("2024-10-08", 5)
if err != nil {
	return
}
_ = day
```

### 交易日区间

```go
days, err := cal.ListTradeDays("2024-10-01", "2024-10-08")
if err != nil {
	return
}
_ = days
```

### 错误判断

```go
import (
	"errors"

	"github.com/Knowckx/tradeday"
)

if err != nil {
	switch {
	case errors.Is(err, tradeday.ErrorInvalidDateFormat):
	case errors.Is(err, tradeday.ErrorDateOutOfRange):
	case errors.Is(err, tradeday.ErrorInvalidDateRange):
	case errors.Is(err, tradeday.ErrorInvalidOffset):
	case errors.Is(err, tradeday.ErrorUnsupportedCalendar):
	}
}
```

## 测试

```bash
go test ./...
```

## 限制

- 每年的12月末，当交易所发布了下一年的交易日历时，作者也会及时更新。此时需要开发者手动更新包依赖。

- 本项目不判断停牌、临时停牌或盘中交易状态
- 本项目不替代交易所官方公告，假如遇到特殊事件，请以交易所官方公告为准。

## 路线图

- 增加更多市场支持，比如香港交易所。
- 扩展更多年份数据
- 增加自动生成日历数据的工具
- 增加命令行工具

## 贡献

欢迎提交 issue、测试用例和交易日历数据修订。修改日历数据时请注明来源，避免引入不可追溯的差异。

## License

MIT License. See [LICENSE](./LICENSE).

## English

`tradeday` is an offline-capable Go trading-calendar library used to check whether a given date is a trading day.  
It ships with built-in trading-day data and is suitable for quant systems, trading applications, backtesting, market-data collection, and scheduled jobs.

The repository currently supports two markets, with the following ranges:

- `CNStock`: `2015-01-01` through `2026-12-31`
- `USStock`: `2015-01-01` through `2026-12-31`

The project uses trading-day bitmaps to store each year's data and delivers excellent query efficiency
- Single-day trading queries are O(1)
- Range trading queries are O(n).

## Features

- Support creating calendars by market
- Support checking whether a given day is a trading day
- Support getting the previous trading day, the next trading day, and an offset trading day
- Support listing trading days in an inclusive date range
- Trading-day checks are based on local bitmap data
- Support real trading-day results for weekends, exchange holidays, and special closures
- Human-friendly API design that is easy to embed in Go projects
- Usable without any network requests

## Installation

```bash
go get github.com/Knowckx/tradeday
```

## Quick Start

`Date` is a string-backed date type, so you can pass a string literal such as `"2024-10-08"` directly.

```go
package main

import (
	"fmt"
	"log"

	"github.com/Knowckx/tradeday"
)

func main() {
	cal, err := tradeday.New(tradeday.CalendarID.CNStock)
	if err != nil {
		log.Fatal(err)
	}

	ok, err := cal.IsTradeDay("2024-10-08")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ok)
}
```

## Core Concepts

- `Date`
  - The public date type, with a fixed format of `2006-01-02`
- `Calendar`
  - A calendar instance for a specific market
- `CalendarID`
  - Used to choose a specific market, currently including `CNStock` and `USStock`
- `Trading-day bitmap`
  - The underlying implementation, where each year corresponds to a bitmap dataset; bit=1 means trading day and bit=0 means non-trading day

## API Usage

### Create a calendar

```go
cal, err := tradeday.New(tradeday.CalendarID.USStock)
if err != nil {
	return
}
```

### Check whether a trading day

```go
ok, err := cal.IsTradeDay("2024-10-08")
if err != nil {
	return
}
_ = ok
```

### Previous / next trading day

```go
prev, err := cal.PrevTradeDay("2024-10-08")
next, err := cal.NextTradeDay("2024-10-08")
_ = prev
_ = next
_ = err
```

### Offset by trading days

```go
// Return the 5th valid trading day after "2024-10-08"
day, err := cal.OffsetTradeDay("2024-10-08", 5)
if err != nil {
	return
}
_ = day
```

### Trading-day range

```go
days, err := cal.ListTradeDays("2024-10-01", "2024-10-08")
if err != nil {
	return
}
_ = days
```

### Error handling

```go
import (
	"errors"

	"github.com/Knowckx/tradeday"
)

if err != nil {
	switch {
	case errors.Is(err, tradeday.ErrorInvalidDateFormat):
	case errors.Is(err, tradeday.ErrorDateOutOfRange):
	case errors.Is(err, tradeday.ErrorInvalidDateRange):
	case errors.Is(err, tradeday.ErrorInvalidOffset):
	case errors.Is(err, tradeday.ErrorUnsupportedCalendar):
	}
}
```

## Testing

```bash
go test ./...
```

## Limitations

- At the end of each year, when the exchange publishes the next year's trading calendar, the author will update it in time. Developers then need to update the package dependency manually.

- This project does not determine suspension, temporary suspension, or intraday trading status
- This project does not replace official exchange announcements; if special events occur, please follow the exchange's official notices

## Roadmap

- Add more market support, such as the Hong Kong exchange.
- Extend the year coverage
- Add tooling for generating calendar data
- Add a command-line tool

## Contributing

Issues, test cases, and calendar-data updates are welcome. When updating calendar data, please include the source so changes remain traceable.

## License

MIT License. See [LICENSE](./LICENSE).
