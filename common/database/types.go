package database

// OpenStockPosition is the data structure of a portfolio record in DynamoDB.
type OpenStockPosition struct {
	PK              string  `json:"-"` // Record type is not needed after DB query
	SK              string  `json:"SK"`
	BuyDate         int64   `json:"BuyDate"`
	AverageBuyPrice float64 `json:"AverageBuyPrice"`
	NumOfShares     uint    `json:"NumOfShares"`
}
