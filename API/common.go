package API

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// buildURL constructs the API query URL for fetching the given symbol's price data
func buildURL(symbol string) string {
	return fmt.Sprintf(
		"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%v&outputsize=compact&apikey=%v",
		symbol, os.Getenv("API_KEY"),
	)
}

// checkDateFormat ensures that the date is in the format YYYY-MM-DD, and is less than the current date.
func checkDateFormat(date string) bool {
	if !regexp.MustCompile(`^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$`).MatchString(date) {
		return false
	}
	formattedCurrentDate := strings.Replace(time.Now().Format("2006-01-02"), "-", "", -1)
	formattedQueryDate := strings.Replace(date, "-", "", -1)
	if formattedQueryDate >= formattedCurrentDate {
		return false
	}
	return true
}

// parseData reads the API response body into a date: price lookup map : [date] => closing-price
func parseData(data []byte) (map[string]float64, error) {
	var stockData = make(map[string]float64)
	var apiResponse QueryResponse
	if err := json.Unmarshal(data, &apiResponse); err != nil {
		return stockData, err
	}

	// For each closing price record returned from search, format to float, and add into a lookup-map
	for key, value := range apiResponse.TimeSeries {
		var priceData TimeSeries
		if byteData, err := json.Marshal(value); err != nil {
			continue
		} else if err = json.Unmarshal(byteData, &priceData); err != nil {
			continue
		}
		formattedPrice, formattingErr := strconv.ParseFloat(priceData.Close, 64)
		if formattingErr != nil {
			continue
		}
		stockData[key] = formattedPrice
	}

	return stockData, nil
}
