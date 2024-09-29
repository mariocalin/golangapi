package date

import "time"

//go:generate mockery --name=Handler --inpackage=true --filename=handler_repository_mock.go --with-expecter=true
type Handler interface {
	DateNow() time.Time
	DateParse(dateInSt string) (time.Time, error)
	DateToString(time time.Time) string
}
