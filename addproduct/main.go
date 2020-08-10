package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	uuid "github.com/satori/go.uuid"
)

type Product struct {
	ProductID  string `json:"productId"`
	Name       string `json:"name"`
	ImageURL   string `json:"imageURL"`
	Price      string `json:"price"`
	OldPrice   string `json:"old_price,omitempty"`
	ProductURL string `json:"productURL,omitempty"`
}

type DB struct {
	UserID   string    `json:"userId"`
	Products []Product `json:"products"`
}

func checkIfProductExists(product string, products DB) bool {
	match := false
	for _, prod := range products.Products {
		fmt.Println(prod.Name)
		if prod.Name == product {
			match = true
		}
	}
	return match
}

// Response is of type APIGatewayProxyResponse.
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	u2 := uuid.NewV4()

	req := &Product{}
	req.ProductID = u2.String()

	json.Unmarshal([]byte(request.Body), req)

	// Create Aurora Serverless Client
	svc := dynamodb.New(sess)

	prod, err := dynamodbattribute.MarshalMap(req)
	if err != nil {
		log.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	log.Println(req.Name)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("products_table")),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(req.Name),
			},
		},
	},
	)

	if err != nil {
		log.Println(err.Error())
	}

	log.Println(result.Item)

	if result.Item == nil {
		input := &dynamodb.PutItemInput{
			Item:      prod,
			TableName: aws.String(os.Getenv("products_table")),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			log.Println(err.Error())
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
		}
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
	}
	return events.APIGatewayProxyResponse{Body: "Product added", Headers: headers, StatusCode: 201}, nil

}

func main() {
	lambda.Start(Handler)
}
