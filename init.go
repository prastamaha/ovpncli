package main

import (
	"fmt"

	"go.etcd.io/bbolt"
)

func InitDatabase(database string) *bbolt.DB {

	db, err := bbolt.Open(database, 0755, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Profiles"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte("Connections"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return db
}
