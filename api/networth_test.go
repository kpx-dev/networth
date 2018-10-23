package main

import (
	"testing"

	_ "github.com/networth-app/networth/dotenv"
	"github.com/networth-app/networth/lib"
	"github.com/stretchr/testify/assert"
)

var (
	db           = nwlib.NewDynamoDBClient()
	testUsername = "c1fa7e12-529e-4b63-8c64-855ba23690ff"
)

func TestGroupByDaily(t *testing.T) {
	startDate := "2018-10-20"
	endDate := "2018-10-21"
	networth, err := db.GetNetworthByDateRange(testUsername, startDate, endDate)

	daily := groupBy(networth, "daily")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(daily) >= 8, true)
}

func TestGroupByWeeklyMonthly(t *testing.T) {
	startDate := "2018-10-01"
	endDate := "2018-11-01"
	networth, err := db.GetNetworthByDateRange(testUsername, startDate, endDate)

	weekly := groupBy(networth, "weekly")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(weekly) >= 4, true)
}
