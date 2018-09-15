package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {

	config := &aws.Config{
		Region:   aws.String("eu-west-1"),
		Endpoint: aws.String("http://localhost:8000"),
	}

	sess := session.Must(session.NewSession(config))
	svc := dynamodb.New(sess)

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("key"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("key"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String("LockTable"),
	}

	result, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
}
