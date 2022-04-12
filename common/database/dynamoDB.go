package database

import (
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
