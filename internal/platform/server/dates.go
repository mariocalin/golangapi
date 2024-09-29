package server

import "time"

func parseDateYearMonthDay(dateInSt string) (time.Time, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateInSt)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
