package date

import "time"

type LocalHandler struct {
	layout   string
	location *time.Location
}

func NewLocalHandler() LocalHandler {
	return LocalHandler{
		layout:   time.DateOnly,
		location: time.Local,
	}
}

func (handler LocalHandler) DateNow() time.Time {
	return time.Now().Local()
}

func (handler LocalHandler) DateParse(dateInSt string) (time.Time, error) {
	return time.ParseInLocation(handler.layout, dateInSt, handler.location)
}

func (handler LocalHandler) DateToString(time time.Time) string {
	return time.Format(handler.layout)
}
