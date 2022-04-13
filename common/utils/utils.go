package utils

import (
	"Investing-API/common/database"
	"Investing-API/common/types"
	"math"
	"time"
)

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

// RoundToPrecision takes a float value and rounds-up to the given precision.
func RoundToPrecision(input float64, precision uint) float64 {
	output := math.Pow(10, float64(precision))
	rounded := int(input*output + math.Copysign(0.5, input*output))
	return float64(rounded) / output
}

// RemovePositionFromPortfolio removes the given index from a slice of portfolio positions.
func RemovePositionFromPortfolio(positions []database.OpenStockPosition, index uint) []database.OpenStockPosition {
	if int(index) >= len(positions) {
		return positions
	}
	return append(positions[:index], positions[index+1:]...)
}

// CombinePositions adds the data of an incoming trade to an existing position. (New Average price, total value, shares quantity...)
func CombinePositions(openPosition database.OpenStockPosition, newTrade types.NewStockTrade) database.OpenStockPosition {
	openPosition.Shares = openPosition.Shares + newTrade.Quantity
	openPosition.PurchaseValue = RoundToPrecision(openPosition.PurchaseValue+(newTrade.Price*float64(newTrade.Quantity)), 2)
	openPosition.AveragePrice = RoundToPrecision(openPosition.PurchaseValue/float64(openPosition.Shares), 2)
	return openPosition
}

// CalculatePortfolioRatio takes a list of open stock positions and calculates the ratio each one takes up in the portfolio.
func CalculatePortfolioRatio(records []database.OpenStockPosition) []database.OpenStockPosition {
	var totalPortfolioValue float64
	for _, record := range records {
		totalPortfolioValue += record.PurchaseValue
	}
	for index, record := range records {
		records[index].PortfolioPercentage = RoundToPrecision(record.PurchaseValue/totalPortfolioValue, 4)
	}
	return records
}
