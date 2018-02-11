// Package store - put.go implements how we store an item
package store

import (
	"github.com/coreos/bbolt"
)

// Close shuts down acces to the underlying database
func (backend *impl) Close() error {
	return backend.db.Close()
}

// Put adds a (k,v) pair into a given namespace
func (backend *impl) Put(namespace []byte, key []byte, value []byte) error {
	return backend.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(namespace)
		err := b.Put(key, value)
		return err
	})
}

// Get returns the data associated with a given key in a namespace
func (backend *impl) Get(namespace []byte, key []byte) ([]byte, error) {
	var result []byte
	err := backend.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(namespace)
		result = b.Get(key)
		return nil
	})

	return result, err
}

// ForEach lets you iterate a function over each k-v pair in a bucket
func (backend *impl) ForEach(namespace []byte, fun func([]byte, []byte) error) error {
	return backend.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(namespace)
		return b.ForEach(fun)
	})
}

// CreateNamespace creates a namespace in our database
func (backend *impl) CreateNamespace(namespace []byte) error {
	return backend.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(namespace)
		return err
	})
}
