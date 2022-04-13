package database

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Login creates a new DynamoDB client we can use to interact with the database.
func Login() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("ACCESS_KEY"),
				os.Getenv("SECRET_KEY"),
				"",
			),
			Region: aws.String(os.Getenv("REGION")),
		},
	}))

	return dynamodb.New(sess)
}

// GetAllOpenPositions queries the database for all active portfolio positions.
func GetAllOpenPositions(svc *dynamodb.DynamoDB) ([]OpenStockPosition, error) {
	var openPositions []OpenStockPosition

	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("PORTFOLIO"),
		KeyConditions: map[string]*dynamodb.Condition{
			"PK": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("OPEN-POSITION"),
					},
				},
			},
		},
	}

	result, queryErr := svc.Query(queryInput)
	if queryErr != nil {
		log.Printf("Error querying DynamoDB: %v\n", queryErr)
		return openPositions, queryErr
	}

	if unmarshallErr := dynamodbattribute.UnmarshalListOfMaps(result.Items, &openPositions); unmarshallErr != nil {
		log.Printf("Error unmarshalling DynamoDB response: %v\n", unmarshallErr)
		return openPositions, unmarshallErr
	}

	return openPositions, nil
}

// AddNewPosition creates a new open portfolio position in the DynamoDB table.
func AddNewPosition(svc *dynamodb.DynamoDB, record OpenStockPosition) error {
	record.PK = "OPEN-POSITION"

	dbRecord, marshallErr := dynamodbattribute.MarshalMap(record)
	if marshallErr != nil {
		log.Printf("Error marshalling record: %v\n", marshallErr)
		return marshallErr
	}

	input := &dynamodb.PutItemInput{
		Item:      dbRecord,
		TableName: aws.String("PORTFOLIO"),
	}

	if _, putItemErr := svc.PutItem(input); putItemErr != nil {
		log.Printf("Error inserting record: %v\n", putItemErr)
		return putItemErr
	}

	return nil
}

// UpdateOpenPosition updates a portfolio record in the DynamoDB table.
func UpdateOpenPosition(svc *dynamodb.DynamoDB, record OpenStockPosition) error {
	queryInput := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":purchaseValue": {
				N: aws.String(fmt.Sprintf("%v", record.PurchaseValue)),
			},
			":portfolioPercentage": {
				N: aws.String(fmt.Sprintf("%v", record.PortfolioPercentage)),
			},
			":averagePrice": {
				N: aws.String(fmt.Sprintf("%v", record.AveragePrice)),
			},
			":percentageReturn": {
				N: aws.String(fmt.Sprintf("%v", record.PercentageReturn)),
			},
			":shares": {
				N: aws.String(fmt.Sprintf("%v", record.Shares)),
			},
			":currentValue": {
				N: aws.String(fmt.Sprintf("%v", record.CurrentValue)),
			},
		},
		TableName: aws.String("PORTFOLIO"),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(record.PK),
			},
			"SK": {
				S: aws.String(record.SK),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set " +
			"PurchaseValue = :purchaseValue, " +
			"PortfolioPercentage = :portfolioPercentage, " +
			"AveragePrice = :averagePrice, " +
			"PercentageReturn = :percentageReturn, " +
			"Shares = :shares, " +
			"CurrentValue = :currentValue",
		),
	}

	if _, err := svc.UpdateItem(queryInput); err != nil {
		log.Printf("Got error calling UpdateItem: %s", err)
		return err
	}

	return nil
}
