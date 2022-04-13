package database

// OpenStockPosition is the data structure of a portfolio record in DynamoDB.
type OpenStockPosition struct {
	PK                  string  `json:"PK"`
	SK                  string  `json:"SK"`
	PurchaseValue       float64 `json:"PurchaseValue"`
	PortfolioPercentage float64 `json:"PortfolioPercentage"`
	AveragePrice        float64 `json:"AveragePrice"`
	PercentageReturn    float64 `json:"PercentageReturn"`
	Shares              uint    `json:"Shares"`
	CurrentStockPrice   float64 `json:"CurrentStockPrice"`
}
