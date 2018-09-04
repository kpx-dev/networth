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
		if record.Change.StreamViewType != "NEW_IMAGE" {
			log.Println("Not a NEW_IMAGE stream view type")
			return
		}

		switch record.EventName {
		case "INSERT":
			for key := range record.Change.Keys {
				msg := "insert key is " + key
				log.Println(msg)
				nwlib.Alert(msg)
			}

			for name, value := range record.Change.NewImage {
				eachMsg := fmt.Sprintf("Each insert Dyno stream, name %s value %v, data type %v", name, value, value.DataType())
				log.Println(eachMsg)
				nwlib.Alert(eachMsg)

				// if value.DataType() == events.DataTypeString {
				// 	fmt.Printf("Attribute name: %s, value: %s\n", name, value.String())
				// }
			}

			break
		case "REMOVE":
			for key := range record.Change.Keys {
				msg := "remove key is " + key
				log.Println(msg)
				nwlib.Alert(msg)
			}

			for name, value := range record.Change.NewImage {
				eachMsg := fmt.Sprintf("Each remove Dyno stream, name %s value %v, data type %v", name, value, value.DataType())
				log.Println(eachMsg)
				nwlib.Alert(eachMsg)
			}
			break
		}
	}
}

func main() {
	lambda.Start(handleDynamoDBStream)
}
