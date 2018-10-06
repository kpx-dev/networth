package main

import (
	"testing"

	_ "github.com/networth-app/networth/api/lib/dotenv"
	"github.com/stretchr/testify/assert"
)

func TestTransactions(t *testing.T) {
	err := syncTransactions(username, invalidToken)
	assert.Equal(t, err != nil, true)
}
