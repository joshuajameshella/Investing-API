package utils

import (
	"Investing-API/common/database"
	"Investing-API/common/types"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

// TestCombinePositions ensures that adding to a position returns the correct updated value.
func TestCombinePositions(t *testing.T) {
	tests := map[string]struct {
		openPosition             database.OpenStockPosition
		newTrade                 types.NewStockTrade
		expectedCombinedPosition database.OpenStockPosition
	}{
		"General Check": {
			database.OpenStockPosition{
				PK:                  "OPEN-POSITION",
				SK:                  "AAPL",
				PurchaseValue:       150.00,
				PortfolioPercentage: 1.0000,
				AveragePrice:        150.00,
				PercentageReturn:    0.1000,
				Shares:              1,
				CurrentValue:        150.00,
			},
			types.NewStockTrade{
				Symbol:   "AAPL",
				Quantity: 1,
				Price:    200.00,
			},
			database.OpenStockPosition{
				PK:                  "OPEN-POSITION",
				SK:                  "AAPL",
				PurchaseValue:       350.00,
				PortfolioPercentage: 1.0000,
				AveragePrice:        175.00,
				PercentageReturn:    0.1000,
				Shares:              2,
				CurrentValue:        150.00,
			},
		},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			combinedPosition := CombinePositions(testCase.openPosition, testCase.newTrade)
			assert.Equal(t, testCase.expectedCombinedPosition, combinedPosition)
		})
	}
}

// TestRoundToPrecision checks that the float-rounding function performs as expected.
func TestRoundToPrecision(t *testing.T) {
	tests := map[string]struct {
		inputFloat     float64
		inputPrecision uint
		expectedOutput float64
	}{
		"General Check":    {100.01233, 2, 100.01},
		"No Precision":     {100, 0, 100},
		"High Precision":   {100, 5, 100.00000},
		"Round Up Decimal": {52.389123, 2, 52.39},
		"Round Up":         {52.9999, 1, 53.0},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			roundedValue := RoundToPrecision(testCase.inputFloat, testCase.inputPrecision)
			assert.Equal(t, testCase.expectedOutput, roundedValue)
		})
	}
}

// TestCalculatePortfolioRatio takes example positions and checks the calculated portfolio sizes match the expected values.
func TestCalculatePortfolioRatio(t *testing.T) {
	tests := map[string]struct {
		openPositions          []database.OpenStockPosition
		expectedPositionRatios []database.OpenStockPosition
	}{
		"Test 1": {
			[]database.OpenStockPosition{
				{SK: "CASH", PurchaseValue: 1000},
				{SK: "AAPL", PurchaseValue: 100},
				{SK: "TSLA", PurchaseValue: 100},
			},
			[]database.OpenStockPosition{
				{SK: "CASH", PurchaseValue: 1000, PortfolioPercentage: 0.8333},
				{SK: "AAPL", PurchaseValue: 100, PortfolioPercentage: 0.0833},
				{SK: "TSLA", PurchaseValue: 100, PortfolioPercentage: 0.0833},
			},
		},
		"Test 2": {
			[]database.OpenStockPosition{
				{SK: "CASH", PurchaseValue: 8612311.44},
				{SK: "AAPL", PurchaseValue: 234424.40},
				{SK: "TSLA", PurchaseValue: 1023.3},
			},
			[]database.OpenStockPosition{
				{SK: "CASH", PurchaseValue: 8612311.44, PortfolioPercentage: 0.9734},
				{SK: "AAPL", PurchaseValue: 234424.40, PortfolioPercentage: 0.0265},
				{SK: "TSLA", PurchaseValue: 1023.3, PortfolioPercentage: 0.0001},
			},
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			calculatedRatios := CalculatePortfolioRatio(testCase.openPositions)
			assert.Equal(t, testCase.expectedPositionRatios, calculatedRatios)
		})
	}
}
