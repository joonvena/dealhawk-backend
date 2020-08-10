package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	ProductID string `json:"productId"`
}

// Response is of type APIGatewayProxyResponse.
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	req := &Request{}

	json.Unmarshal([]byte(request.Body), req)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"productId": {
				S: aws.String(req.ProductID),
			},
		},
		TableName: aws.String(os.Getenv("products_table")),
	}

	_, err := svc.DeleteItem(input)

	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
	}

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Delete failed", Headers: headers, StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: "Product removed", Headers: headers, StatusCode: 200}, nil

}

func main() {
	lambda.Start(Handler)
}
