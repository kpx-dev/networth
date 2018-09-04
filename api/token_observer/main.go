package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/networth-app/networth/api/lib"
)

func handleDynamoDBStream(ctx context.Context, e events.DynamoDBEvent) {
	for _, record := range e.Records {
		msg := fmt.Sprintf("Dyno stream type %s, source %s, stream view type %s keys %v", record.EventName, record.EventSource, record.Change.StreamViewType, record.Change.Keys)
		log.Println(msg)
		nwlib.Alert(msg)

		for name, value := range record.Change.NewImage {
			eachMsg := fmt.Sprintf("Each Dyno stream, name %s value %v, data type %v", name, value, value.DataType())
			log.Println(eachMsg)
			nwlib.Alert(eachMsg)

			// if value.DataType() == events.DataTypeString {
			// 	fmt.Printf("Attribute name: %s, value: %s\n", name, value.String())
			// }
		}
	}
}

func main() {
	lambda.Start(handleDynamoDBStream)
}
