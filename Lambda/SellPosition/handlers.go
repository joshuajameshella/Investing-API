package main

import "Investing-API/common/database"

// getPositionOfInterest looks in a slice of portfolio positions for a specific symbol.
func getPositionOfInterest(openPositions []database.OpenStockPosition, symbol string) (database.OpenStockPosition, uint, bool) {
	for index, position := range openPositions {
		if position.SK == symbol {
			return position, uint(index), true
		}
	}
	return database.OpenStockPosition{}, 0, false
}
