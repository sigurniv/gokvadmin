package gokvadmin

import (
	"github.com/boltdb/bolt"
	"bytes"
	"errors"
)

type EngineBoltDB struct {
	db     *bolt.DB
	bucket []byte
}

func init() {
	RegisterEngine(engineBoltDB, &EngineBoltDB{})
}

func (e *EngineBoltDB) SetDB(db interface{}) error {
	edb, ok := db.(*bolt.DB)
	if !ok {
		return errors.New("db must be an instance of *bolt.DB")
	}

	e.db = edb
	return nil
}

func (e EngineBoltDB) GetName() string {
	return "boltdb"
}

func (e EngineBoltDB) Get(key []byte, bucket []byte) ([]byte, error) {
	var value []byte

	err := e.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}

		value = b.Get(key)
		return nil
	})

	return value, err
}

func (e EngineBoltDB) GetByPrefix(prefix []byte, bucket []byte, limit int, offset int) ([]Record, error) {
	var records []Record

	err := e.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}

		c := b.Cursor()
		outer : for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix) && limit > 0; k, v = c.Next() {
			offset--;
			for offset >= 0 {
				continue outer
			}

			records = append(records, Record{Key: k, Value:v})
			limit--
			continue outer

		}

		return nil
	})

	return records, err
}

func (e EngineBoltDB) Set(key []byte, value []byte, bucket []byte) error {
	var err error

	e.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(bucket)
		return nil
	})

	if err != nil {
		return err
	}

	return e.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		err := b.Put(key, value)
		return err
	})
}

func (e EngineBoltDB) Delete(key []byte, bucket []byte) error {
	var err error

	e.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}

		err = b.Delete(key)
		return nil
	})

	return err
}



