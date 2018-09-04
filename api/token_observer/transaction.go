package main

import (
	"fmt"
)

func transactions(username string, accessTokens []string) {
	for _, token := range accessTokens {
		fmt.Println("handling trans ", token)
	}
}
