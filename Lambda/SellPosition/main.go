package main

import (
	"Investing-API/Lambda/lambdaHandler"
	"Investing-API/common/database"
	"Investing-API/common/types"
	"Investing-API/common/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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

	// Look for the specified trade in the portfolio.
	queryPosition, positionIndex, exists := getPositionOfInterest(openPositions, input.Symbol)
	if !exists {
		var errMsg = fmt.Sprintf("Cannot find %v in the portfolio", input.Symbol)
		log.Println(errMsg)
		return lambdaHandler.Response(http.StatusInternalServerError, errMsg)
	}

	// Check that the user isn't requesting to sell more shares than they own.
	if input.Quantity > queryPosition.Shares {
		var errMsg = fmt.Sprintf("Cannot sell more shares than you own. You have %v shares in your account", queryPosition.Shares)
		log.Println(errMsg)
		return lambdaHandler.Response(http.StatusBadRequest, errMsg)
	}

	// Calculate the value of the trade
	sellPrice := utils.RoundToPrecision(input.Price*float64(input.Quantity), 2)

	// If the user is selling all their shares, delete the record. Otherwise, update the record.
	if input.Quantity == queryPosition.Shares {
		if deleteErr := database.DeleteOpenPosition(svc, queryPosition); deleteErr != nil {
			log.Printf("Error removing position from portfolio: %v\n", deleteErr)
			return lambdaHandler.Response(http.StatusInternalServerError, deleteErr)
		}
		openPositions = utils.RemovePositionFromPortfolio(openPositions, positionIndex)
	} else {
		openPositions[positionIndex].PurchaseValue = queryPosition.PurchaseValue - sellPrice
		openPositions[positionIndex].Shares = queryPosition.Shares - input.Quantity
	}

	// Update each position's ratio's data.
	updatedRecords := utils.CalculatePortfolioRatio(openPositions)

	// Update each portfolio position in the database.
	for index, position := range updatedRecords {
		if position.SK == "CASH" {
			openPositions[index].PurchaseValue = position.PurchaseValue + sellPrice
		}
		if updateErr := database.UpdateOpenPosition(svc, openPositions[index]); updateErr != nil {
			log.Printf("Error updating %v record: %v\n", position.SK, updateErr)
			return lambdaHandler.Response(http.StatusInternalServerError, updateErr)
		}
	}

	log.Println("Successfully sold stock position!")
	return lambdaHandler.Response(http.StatusOK, "Successfully sold stock position!")
}
