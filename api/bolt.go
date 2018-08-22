package main

import (
	"fmt"
	"log"
	"strconv"

	bolt "github.com/coreos/bbolt"
)

const dbBucket = "networth"

// BoltClient bolt client struct
type BoltClient struct {
	*bolt.DB
}

// NewBoltClient new bolt client
func NewBoltClient() *BoltClient {
	path := "/tmp/" + dbBucket + ".db"
	db, err := bolt.Open(path, 0600, nil)

	if err != nil {
		panic(err)
	}

	client := &BoltClient{db}
	client.init()

	return client
}

// init create required buckets
func (db *BoltClient) init() {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(dbBucket))
		if err != nil {
			return fmt.Errorf("Error creating bucket: %s", err)
		}
		return nil
	})
}

// Set key value
func (db *BoltClient) Set(key string, value string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		err := b.Put([]byte(username+":"+key), []byte(value))

		return err
	})

	if err != nil {
		log.Println(err)
	}

	return err
}

// Get value
func (db *BoltClient) Get(key string) string {
	var payload string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		v := b.Get([]byte(username + ":" + key))
		payload = string(v)

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	return payload
}

// SetNetworth set current networth
func (db *BoltClient) SetNetworth(networth float64) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		err := b.Put([]byte(username), []byte(strconv.FormatFloat(networth, 'f', 6, 64)))

		return err
	})

	if err != nil {
		log.Println(err)
	}

	return err
}

// GetAccessToken return the access token
func (db *BoltClient) GetAccessToken() string {
	return db.Get(username + ":access_token")
}

// GetNetworth get current networth
func (db *BoltClient) GetNetworth() float64 {
	networth := 0.0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		v := b.Get([]byte(username))
		nw, err := strconv.ParseFloat(string(v), 64)
		networth = nw

		return err
	})

	if err != nil {
		log.Println(err)
	}

	return networth
}
