package tradeday

import "github.com/Knowckx/tradeday/core"

type Error = core.Error
type Date = core.Date
type CalendarID = core.CalendarID
type Calendar = core.Calendar

const (
	DateLayout = core.DateLayout

	CalendarCNStock CalendarID = core.CalendarCNStock
	CalendarUSStock CalendarID = core.CalendarUSStock
)

var ParseDate = core.ParseDate

func New(calendarID CalendarID) (Calendar, error) {
	return core.New(calendarID)
}
