package api

import (
	"fmt"
	// "log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/tkhrk1010/go-samples/scatter-gather/fruits-trading-demo/infra/ddb"
)

type ApplePurchasing struct {
	ID          string    `json:"id"`
	Count       int       `json:"count"`
	Price       float32   `json:"price"`
	PurchasedAt time.Time `json:"purchasedAt"`
}

func HandleApplePurchasingPost(c *gin.Context) {
	fmt.Println("HandleApplePurchasingPost")

	ddb.InitDynamoDBClient()

	r := gin.Default()

	r.POST("/add/purchasing/apple", func(c *gin.Context) {
		// FIXME: postしても、ここまで到達しない
		fmt.Println("POST /add/purchasing/apple")

		var applePurchasing ApplePurchasing
		if err := c.ShouldBindJSON(&applePurchasing); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := insertApplePurchasing(&applePurchasing); err != nil {
			fmt.Printf("Error inserting ApplePurchasing: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ApplePurchasing inserted successfully"})
	})

}

func insertApplePurchasing(applePurchasing *ApplePurchasing) error {
	item, err := dynamodbattribute.MarshalMap(applePurchasing)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("applePurchasing"),
	}

	_, err = ddb.DynamoDBClient.PutItem(input)
	return err
}
