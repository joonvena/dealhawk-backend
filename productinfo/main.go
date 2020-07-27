package main

import (
	"encoding/json"
	"strings"

	"github.com/gocolly/colly"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Product struct {
	Name       string `json:"name"`
	ImageURL   string `json:"imageURL"`
	Price      string `json:"price"`
	ProductURL string `json:"productURL"`
}

type URL struct {
	URL string `json:"itemURL"`
}

// Response is of type APIGatewayProxyResponse.
type Response events.APIGatewayProxyResponse

func getPrice(c *colly.Collector, url string) Product {
	item := Product{}
	c.OnHTML("#main", func(e *colly.HTMLElement) {
		item.Name = e.ChildText(".product-header-title")
		item.ImageURL = e.ChildAttr(".product-image", "src")
		price := strings.Join(strings.Fields(e.ChildText(".price-tag-price__euros")), "")
		item.Price = price
	})
	c.Visit(url)
	return item
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	req := &URL{}
	json.Unmarshal([]byte(request.Body), req)

	if len(req.URL) < 10 {
		return events.APIGatewayProxyResponse{Body: "URL was not provided or too short.", StatusCode: 500}, nil
	}

	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0"))

	resp := getPrice(c, req.URL)

	if len(resp.Name) > 0 {
		body, _ := json.Marshal(resp)

		return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
	} else {
		return events.APIGatewayProxyResponse{Body: "Product not found", StatusCode: 404}, nil
	}

}

func main() {
	lambda.Start(Handler)
}
