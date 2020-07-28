package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Product struct {
	Name       string `json:"name"`
	ImageURL   string `json:"imageURL"`
	Price      string `json:"price"`
	ProductURL string `json:"productURL,omitempty"`
}

type DB struct {
	UserID   string    `json:"userId"`
	Products []Product `json:"products"`
}

// Response is of type APIGatewayProxyResponse.
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("product-scraper"),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String("fewkfok435ok43pogskg0cml39639"),
			},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	product := DB{}

	if err := dynamodbattribute.UnmarshalMap(result.Item, &product); err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	if len(product.UserID) == 0 {
		return events.APIGatewayProxyResponse{Body: "User not found", StatusCode: 404}, nil
	}

	body, _ := json.Marshal(product)
	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
	}
	return events.APIGatewayProxyResponse{Body: string(body), Headers: headers, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
