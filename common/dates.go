package common

import "time"

type DateHandler struct {
	layout   string
	location *time.Location
}

func NewDateHandler() *DateHandler {
	return &DateHandler{
		layout:   time.DateOnly,
		location: time.Local,
	}
}

func (handler *DateHandler) DateNow() time.Time {
	return time.Now().Local()
}

func (handler *DateHandler) DateParse(dateInSt string) (time.Time, error) {
	return time.ParseInLocation(handler.layout, dateInSt, handler.location)
}

func (handler *DateHandler) DateToString(time time.Time) string {
	return time.Format(handler.layout)
}
