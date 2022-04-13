package main

import "Investing-API/common/database"

// canAffordTrade checks that there is enough cash in the portfolio to afford the new trade.
func canAffordTrade(openPositions []database.OpenStockPosition) bool {
	var totalCash float64
	for _, position := range openPositions {
		if position.SK == "CASH" {
			totalCash = position.CurrentValue
		}
	}
	return totalCash >= newTradeValue
}

// recalculateCashValue removes the cash from the portfolio that has been used for the trade.
func recalculateCashValue(openPositions []database.OpenStockPosition) []database.OpenStockPosition {
	for index, position := range openPositions {
		if position.SK == "CASH" {
			var newValue = position.PurchaseValue - newTradeValue
			openPositions[index].PurchaseValue = newValue
			openPositions[index].CurrentValue = newValue
		}
	}
	return openPositions
}
