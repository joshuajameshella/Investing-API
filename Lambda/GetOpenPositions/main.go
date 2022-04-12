package main

import (
	"Investing-API/Lambda/lambdaHandler"
	"Investing-API/database"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Process)
}

func Process(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("Incoming request from: %v\n", request.RequestContext.Identity.SourceIP)

	svc := database.Login()

	openPositions, dbQueryErr := database.GetAllOpenPositions(svc)
	if dbQueryErr != nil {
		log.Printf("Error querying database for open portfolio positions: %v\n", dbQueryErr)
		return lambdaHandler.Response(http.StatusInternalServerError, dbQueryErr)
	}

	return lambdaHandler.Response(http.StatusOK, openPositions)
}
