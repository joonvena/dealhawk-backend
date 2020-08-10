package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Product struct {
	ProductID  string `json:"productId"`
	Name       string `json:"name"`
	ImageURL   string `json:"imageURL"`
	Price      string `json:"price"`
	OldPrice   string `json:"old_price,omitempty"`
	ProductURL string `json:"productURL"`
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

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("products_table")),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	var product []Product

	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &product); err != nil {
		fmt.Println(err.Error())
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
