package lambdaHandler

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// Response builds the AWS Lambda response to return to the user following each API request.
func Response(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	resp.StatusCode = status

	stringBody, stringErr := json.Marshal(body)
	resp.Body = string(stringBody)

	return &resp, stringErr
}
