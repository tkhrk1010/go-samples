package ddb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBClient *dynamodb.DynamoDB

func InitDynamoDBClient() {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localstack:4566"),
		Region:   aws.String("us-east-1"),
	}))
	DynamoDBClient = dynamodb.New(sess)
}
