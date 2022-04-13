package types

// NewStockTrade is the data structure of a new stock trade made.
type NewStockTrade struct {
	Symbol   string  `json:"Symbol"`
	Quantity uint    `json:"Quantity"`
	Price    float64 `json:"Price"`
}
