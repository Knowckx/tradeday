package base

// Error 是本项目唯一公开错误类型。
type Error string

const (
	ErrorUnsupportedCalendar Error = "unsupported_calendar"
	ErrorDateOutOfRange      Error = "date_out_of_range"
	ErrorInvalidDateRange    Error = "invalid_date_range"
	ErrorInvalidOffset       Error = "invalid_offset"
	ErrorInvalidDateFormat   Error = "invalid_date_format"
)

func (e Error) Error() string {
	return string(e)
}

func (e Error) IsUnsupportedCalendar() bool {
	return e == ErrorUnsupportedCalendar
}

func (e Error) IsDateOutOfRange() bool {
	return e == ErrorDateOutOfRange
}

func (e Error) IsInvalidDateRange() bool {
	return e == ErrorInvalidDateRange
}

func (e Error) IsInvalidOffset() bool {
	return e == ErrorInvalidOffset
}

func (e Error) IsInvalidDateFormat() bool {
	return e == ErrorInvalidDateFormat
}

func NewUnsupportedCalendarError() error {
	return ErrorUnsupportedCalendar
}

func NewDateOutOfRangeError() error {
	return ErrorDateOutOfRange
}

func NewInvalidDateRangeError() error {
	return ErrorInvalidDateRange
}

func NewInvalidOffsetError() error {
	return ErrorInvalidOffset
}

func NewInvalidDateFormatError() error {
	return ErrorInvalidDateFormat
}
