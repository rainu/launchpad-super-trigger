package store

import (
	"github.com/boltdb/bolt"
)

type BoltDBStore struct {
	BoltDB *bolt.DB
	Bucket string
	Key    string
}

func (b *BoltDBStore) Set(data []byte) error {
	return b.BoltDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(b.Bucket))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(b.Key), data)
	})
}

func (b *BoltDBStore) Get() ([]byte, error) {
	var result []byte

	err := b.BoltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.Bucket))
		if bucket != nil {
			result = bucket.Get([]byte(b.Key))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
