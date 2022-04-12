package API

// QueryResponse is the container response that is returned from the Stock-Price query.
type QueryResponse struct {
	MetaData   MetaData               `json:"Meta Data"`
	TimeSeries map[string]interface{} `json:"Time Series (Daily)"`
}

// MetaData contains the top-level info of the stock symbol being queried.
type MetaData struct {
	Info          string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

// TimeSeries is the stock-price structure of each day returned from the API query.
type TimeSeries struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}
