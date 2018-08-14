package main

import (
	bolt "github.com/coreos/bbolt"
)

// BoltClient bolt client struct
type BoltClient struct {
	*bolt.DB
}

// NewBoltClient new bolt client
func NewBoltClient() *BoltClient {
	db, _ := bolt.Open("bolt.db", 0600, nil)
	defer db.Close()

	return &BoltClient{db}
}

// Init create required buckets
// func (db *BoltClient) Init() {
// 	db.Update(func(tx *bolt.Tx) error {
// 		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
// 		if err != nil {
// 			return fmt.Errorf("create bucket: %s", err)
// 		}
// 		return nil
// 	})
// }

// GetNetworth get current networth
func (c *BoltClient) GetNetworth() string {
	return "0"
}
