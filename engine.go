package gokvadmin

import "github.com/gorilla/mux"

type Record struct {
	Key   []byte
	Value []byte
}

type Engine interface {
	GetName() string

	SetDB(db interface{}) error

	Get(key []byte, bucket []byte) ([]byte, error)

	GetByPrefix(prefix []byte, bucket []byte, limit int, offset int) ([]Record, error)

	Set(key []byte, value []byte, bucket []byte) error

	Delete(key []byte, bucket []byte) error
}

type RoutedEngine interface {
	AddRoutes(*mux.Router)
}