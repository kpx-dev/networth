package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleDynamoDBStream(ctx context.Context, e events.DynamoDBEvent) {
	for _, record := range e.Records {
		fmt.Println("Dyno stream type", record.EventName)

		for name, value := range record.Change.NewImage {
			fmt.Println("value data type is ", value.DataType())
			fmt.Println("New Image name, val", name, value)

			if value.DataType() == events.DataTypeString {
				fmt.Printf("Attribute name: %s, value: %s\n", name, value.String())
			}
		}
	}
}

func main() {
	lambda.Start(handleDynamoDBStream)
}
