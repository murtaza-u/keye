package main

import (
	"fmt"
	"log"

	"go.etcd.io/bbolt"
)

const bucket = "KEYE"

func main() {
	db, err := bbolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	keys := []string{
		"foo",
		"bar",
		"blah",
		"foo1",
		"foo2",
		"foo3",
		"foo/bar/1",
		"foo/bar/2",
		"foo/bar/3",
		"foo/bar/blah",
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		for _, k := range keys {
			err := b.Put([]byte(k), []byte(fmt.Sprintf("val_%s", k)))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
