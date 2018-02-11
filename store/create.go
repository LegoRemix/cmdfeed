// Package store contains all of the functions used to create and manage a state store
package store

import (
	"path"

	"github.com/coreos/bbolt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

//database name where
const dbName = ".cmdfeeddb"

type impl struct {
	db *bolt.DB
}

// NewLocalBackend creates a new local datastore
func NewLocalBackend() (Backend, error) {

	home, err := homedir.Dir()
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database - nonexistent home dir")
	}

	dbPath := path.Join(home, dbName)
	db, err := bolt.Open(dbPath, 660, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database - cannot open db file")
	}
	return &impl{
		db: db,
	}, nil
}
