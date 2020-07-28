package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func AddItem(svc *dynamodb.DynamoDB, name string, imageURL string, price string, productURL string) {

	av, err := dynamodbattribute.MarshalMap(Product{
		Name:       name,
		ImageURL:   imageURL,
		Price:      price,
		ProductURL: productURL,
	})

	fmt.Println(av)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	avx := &dynamodb.AttributeValue{
		M: av,
	}

	var r []*dynamodb.AttributeValue
	r = append(r, avx)

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String("fewkfok435ok43pogskg0cml39639"),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":gid": {
				L: r,
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		TableName:        aws.String("product-scraper"),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("SET products = list_append(if_not_exists(products, :empty_list), :gid)"),
	}
	_, err = svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
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
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	req := &Product{}
	json.Unmarshal([]byte(request.Body), req)

	test := DB{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &test)

	exists := checkIfProductExists(req.Name, test)

	if !exists {
		fmt.Println("Item does not exists. Moving to adding item...")
		AddItem(svc, req.Name, req.ImageURL, req.Price, req.ProductURL)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record %v", err))
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
	}
	return events.APIGatewayProxyResponse{Body: "hello", Headers: headers, StatusCode: 200}, nil

}

func main() {
	lambda.Start(Handler)
}
