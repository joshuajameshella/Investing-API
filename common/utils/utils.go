package utils

import "time"

var dateTimeFormat = "2006-01-02"

// CanRun checks whether the code can run on the given day.
// API Data is only updated at midnight on each weekday, so the code should run on each following day.
func CanRun(dayOfWeek time.Weekday) bool {
	return 2 <= int(dayOfWeek) && int(dayOfWeek) <= 6
}

// GetYesterdaysDate returns the formatted date of the previous day.
func GetYesterdaysDate(today time.Time) string {
	return today.AddDate(0, 0, -1).Format(dateTimeFormat)
}
