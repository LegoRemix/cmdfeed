// Package store contains all of the functions used to create and manage a state store
package store

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

//database name where
const dbName = ".cmdfeeddb"

// NewLocalBackend creates a new local datastore
func NewLocalBackend() (Backend, error) {

	home, err := homedir.Dir()
	if err != nil {
		return nil, errors.Wrap(err, "cannot create database")
	}

}
