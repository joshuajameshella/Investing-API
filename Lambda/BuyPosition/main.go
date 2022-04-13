package main

import (
	"Investing-API/Lambda/lambdaHandler"
	"Investing-API/common/database"
	"Investing-API/common/types"
	"Investing-API/common/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var newTradeValue = 0.0
var positionAlreadyExists = false

func main() {
	lambda.Start(Process)
}

func Process(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod != "POST" {
		return lambdaHandler.Response(http.StatusInternalServerError, "Incorrect HTTP method supplied. Need: POST")
	}

	var input = types.NewStockTrade{}
	if unmarshallErr := json.Unmarshal([]byte(request.Body), &input); unmarshallErr != nil {
		log.Printf("Error reading request body into struct: %v\n", unmarshallErr)
		return lambdaHandler.Response(http.StatusInternalServerError, unmarshallErr)
	}

	svc := database.Login()
	openPositions, dbQueryErr := database.GetAllOpenPositions(svc)
	if dbQueryErr != nil {
		log.Printf("Error querying database for open portfolio positions: %v\n", dbQueryErr)
		return lambdaHandler.Response(http.StatusInternalServerError, dbQueryErr)
	}

	// Calculate how much the new trade will cost (will subtract this value from the CASH position).
	newTradeValue = utils.RoundToPrecision(input.Price*float64(input.Quantity), 2)

	// Check that there is enough cash in the portfolio to make the trade.
	if !canAffordTrade(openPositions) {
		log.Printf("Error - not enough cash to enter position")
		return lambdaHandler.Response(http.StatusBadRequest, "not enough cash to enter position!")
	}

	// If a position in the new stock exists, combine the two records.
	for index, position := range openPositions {
		if position.SK == input.Symbol {
			positionAlreadyExists = true
			openPositions[index] = utils.CombinePositions(position, input)
		}
	}

	// If the position doesn't exist, create a new portfolio record.
	if !positionAlreadyExists {
		newPosition := database.OpenStockPosition{
			SK:                  input.Symbol,
			PurchaseValue:       utils.RoundToPrecision(float64(input.Quantity)*input.Price, 2),
			PortfolioPercentage: 0,
			AveragePrice:        utils.RoundToPrecision(input.Price, 2),
			PercentageReturn:    0,
			Shares:              input.Quantity,
			CurrentStockPrice:   utils.RoundToPrecision(input.Price, 2),
		}
		if addRecordErr := database.AddNewPosition(svc, newPosition); addRecordErr != nil {
			log.Printf("Error adding new position into database: %v\n", addRecordErr)
			return lambdaHandler.Response(http.StatusInternalServerError, addRecordErr)
		}
		openPositions = append(openPositions, newPosition)
	}

	// Remove the trade cost from the cash value, and update the position ratio's data.
	updatedRecords := utils.CalculatePortfolioRatio(recalculateCashValue(openPositions))

	// Insert the updated records into the DynamoDB table.
	for _, position := range updatedRecords {
		if updateErr := database.UpdateOpenPosition(svc, position); updateErr != nil {
			log.Printf("Error updating position %v in database: %v\n", position.SK, updateErr)
			return lambdaHandler.Response(http.StatusInternalServerError, updateErr)
		}
	}

	log.Println("Successfully added new stock position!")
	return lambdaHandler.Response(http.StatusOK, "Successfully added new stock position!")
}
