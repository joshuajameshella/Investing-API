package API

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// GetSymbolDatePrice looks up the price of a symbol on a specific date. The date should be in the format YYYY-MM-DD
func GetSymbolDatePrice(symbol, date string) (float64, error) {
	var price float64

	// Check that the date matches the expected format of YYYY-MM-DD
	if !checkDateFormat(date) {
		dateErr := fmt.Sprintf("Incorrect date format. expecting YYYY-MM-DD, but got: \t %v \n", date)
		return price, errors.New(dateErr)
	}

	queryURL := buildURL(symbol)
	response, requestErr := http.Get(queryURL)
	if requestErr != nil {
		log.Printf("Error while quierying URL: %v\n", requestErr)
		return price, requestErr
	}

	// Check for non-successful response codes from the API
	if response.StatusCode != 200 {
		log.Printf("Unexpected StatusCode returned from query: %v\n", response.StatusCode)
	}

	responseData, responseErr := ioutil.ReadAll(response.Body)
	if responseErr != nil {
		log.Printf("Error while reading API response body: %v\n", responseErr)
		return price, responseErr
	}

	priceMap, parseErr := parseData(responseData)
	if parseErr != nil {
		log.Printf("Error while structuring price data: %v\n", responseErr)
		return price, responseErr
	}

	data, exists := priceMap[date]
	if !exists {
		log.Printf("Data for the follwoing data does not exist: %v\n", date)
		return price, responseErr
	}

	return data, nil
}
