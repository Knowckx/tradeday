package core

// Error 是本项目唯一公开错误类型。
type Error string

const (
	errorUnsupportedCalendar Error = "unsupported_calendar"
	errorDateOutOfRange      Error = "date_out_of_range"
	errorInvalidDateRange    Error = "invalid_date_range"
	errorInvalidOffset       Error = "invalid_offset"
	errorInvalidDateFormat   Error = "invalid_date_format"
)

func (e Error) Error() string {
	return string(e)
}

func (e Error) IsUnsupportedCalendar() bool {
	return e == errorUnsupportedCalendar
}

func (e Error) IsDateOutOfRange() bool {
	return e == errorDateOutOfRange
}

func (e Error) IsInvalidDateRange() bool {
	return e == errorInvalidDateRange
}

func (e Error) IsInvalidOffset() bool {
	return e == errorInvalidOffset
}

func (e Error) IsInvalidDateFormat() bool {
	return e == errorInvalidDateFormat
}

func newUnsupportedCalendarError() error {
	return errorUnsupportedCalendar
}

func newDateOutOfRangeError() error {
	return errorDateOutOfRange
}

func newInvalidDateRangeError() error {
	return errorInvalidDateRange
}

func newInvalidOffsetError() error {
	return errorInvalidOffset
}

func newInvalidDateFormatError() error {
	return errorInvalidDateFormat
}
