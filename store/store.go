// Package store contains the backend data storage for our apps
package store

// Backend implements the storage functionality for this apps
type Backend interface {
	Put(namespace []byte, key []byte, value []byte) error
	Get(namespace []byte, key []byte) ([]byte, error)
	ForEach(namespace []byte, fun func([]byte, []byte) error) error
	CreateNamespace(namespace []byte) error
}
