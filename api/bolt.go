package main

import (
	"fmt"

	bolt "github.com/coreos/bbolt"
)

const dbBucket = "networth"

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
func (db *BoltClient) Init() {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(dbBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

// GetNetworth get current networth
func (db *BoltClient) GetNetworth() float64 {
	// networth := 0
	fmt.Println("GetNetworth", username)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		err := b.Put([]byte(username), []byte("42"))
		return err
	})

	db.View(func(tx *bolt.Tx) error {
		fmt.Println("inside view...")
		b := tx.Bucket([]byte(dbBucket))
		v := b.Get([]byte(username))
		// networth =
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})

	// db.Update(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("MyBucket"))
	// 	err := b.Put([]byte("answer"), []byte("42"))
	// 	return err
	// })

	return 0
}
