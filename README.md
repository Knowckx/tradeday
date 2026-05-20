# tradeday

- 每次新窗口进入本项目，请按以下顺序读取文档：
1. `_dev_docs/1.ai-init.md`

## 用法

```go
cal, err := tradeday.New(tradeday.CalendarID.CNStock)
if err != nil {
	return
}

ok, err := cal.IsTradeDay("2024-10-08")
_ = ok
_ = err
```
