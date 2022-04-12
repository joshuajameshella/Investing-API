package utils

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

// TestCanRun checks that only days following a weekday can run the code.
// This is because the API data has an update delay of 24 hours.
func TestCanRun(t *testing.T) {
	var (
		Monday    = time.Weekday(1) // Shouldn't be able to run - stock market closed on Sundays
		Tuesday   = time.Weekday(2) // Should be able to run - stock market open on Mondays
		Wednesday = time.Weekday(3) // Should be able to run - stock market open on Tuesdays
		Thursday  = time.Weekday(4) // Should be able to run - stock market open on Wednesdays
		Friday    = time.Weekday(5) // Should be able to run - stock market open on Thursdays
		Saturday  = time.Weekday(6) // Should be able to run - stock market open on Fridays
		Sunday    = time.Weekday(0) // Shouldn't be able to run - stock market closed on Saturday
	)

	tests := map[string]struct {
		weekday time.Weekday
		canRun  bool
	}{
		"Run on Monday":    {Monday, false},
		"Run on Tuesday":   {Tuesday, true},
		"Run on Wednesday": {Wednesday, true},
		"Run on Thursday":  {Thursday, true},
		"Run on Friday":    {Friday, true},
		"Run on Saturday":  {Saturday, true},
		"Run on Sunday":    {Sunday, false},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			check := CanRun(testCase.weekday)
			assert.Equal(t, testCase.canRun, check)
		})
	}
}

// TestYesterdaysDate checks that the correct date & format is returned from the GetYesterdaysDate function.
func TestYesterdaysDate(t *testing.T) {
	t1, _ := time.Parse(dateTimeFormat, "2022-02-01")
	t2, _ := time.Parse(dateTimeFormat, "2022-01-01")
	t3, _ := time.Parse(dateTimeFormat, "2020-06-12")

	tests := map[string]struct {
		inputDate    time.Time
		expectedDate string
	}{
		"Previous Month": {t1, "2022-01-31"},
		"Previous Year":  {t2, "2021-12-31"},
		"Previous Date":  {t3, "2020-06-11"},
	}

	for name, testCase := range tests {
		// Check that the returned date is in the correct format
		if !regexp.MustCompile(`^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$`).MatchString(testCase.inputDate.Format(dateTimeFormat)) {
			t.Error("calculated date does not match the required format")
		}
		t.Run(name, func(t *testing.T) {
			calculatedDate := GetYesterdaysDate(testCase.inputDate)
			assert.Equal(t, testCase.expectedDate, calculatedDate)
		})
	}
}
