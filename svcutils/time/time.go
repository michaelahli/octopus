package time

import "time"

const (
	FormatDDMMYYY  = "02-01-2006"
	FormatYYYYMMDD = "2006-01-02"
)

func TruncateToStartOfDay(t time.Time) time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

func TruncateToEndOfDay(t time.Time) time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, loc)
}
