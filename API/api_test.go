package API

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDateFormat checks that only valid dates can be sent to the API price query function.
func TestDateFormat(t *testing.T) {
	tests := map[string]struct {
		date      string
		regexPass bool
	}{
		"Correct Date":        {"2022-01-12", true},
		"Incorrect Date":      {"2022-01-8", false},
		"Out Of Bounds Month": {"2022-13-12", false},
		"Out Of Bounds Day":   {"2022-11-52", false},
		"Future Date":         {time.Now().AddDate(0, 0, 1).Format("2006-01-02"), false},
		"Today's Date":        {time.Now().Format("2006-01-02"), false},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			check := checkDateFormat(testCase.date)
			assert.Equal(t, testCase.regexPass, check)
		})
	}
}
