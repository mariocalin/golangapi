package date_test

import (
	"library-api/kit/date"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocalDateHandler(t *testing.T) {
	dateHandler := date.NewLocalHandler()

	t.Run("Obtain date now", func(t *testing.T) {
		now := dateHandler.DateNow()
		assert.NotNil(t, now)
	})

	t.Run("Parse date with DateOnly layout", func(t *testing.T) {
		dateSt := "2024-01-20"
		parsed, err := dateHandler.DateParse(dateSt)
		assert.NotNil(t, parsed)
		assert.Nil(t, err)
		assert.Equal(t, 2024, parsed.Year())
		assert.Equal(t, 20, parsed.Day())
		assert.Equal(t, 1, int(parsed.Month()))
	})

	t.Run("Parse date with wrong layout results in error", func(t *testing.T) {
		dateSt := "2024-01-20 20:23:00"
		_, err := dateHandler.DateParse(dateSt)
		assert.NotNil(t, err)
	})

	t.Run("Format date", func(t *testing.T) {
		date := time.Date(2024, 5, 15, 0, 0, 0, 0, time.Local)
		assert.Equal(t, "2024-05-15", dateHandler.DateToString(date))
	})
}
