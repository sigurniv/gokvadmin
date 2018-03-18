package gokvadmin

import (
	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
)

type EngineBadger struct {
	db     *badger.DB
	bucket []byte
}

func init() {
	RegisterEngine(engineBadger, &EngineBadger{})
}


func (e *EngineBadger) SetDB(db interface{}) error {
	edb, ok := db.(*badger.DB)
	if !ok {
		return errors.New("db must be an instance of *badger.DB")
	}

	e.db = edb
	return nil
}

func (e EngineBadger) GetName() string {
	return "badger"
}

func (e EngineBadger) Get(key []byte, bucket []byte) ([]byte, error) {
	var value []byte

	err := e.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		value, err = item.Value()
		return err
	})

	return value, err
}

func (e EngineBadger) GetByPrefix(prefix []byte, bucket []byte, limit int, offset int) ([]Record, error) {
	var records []Record

	err := e.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		outer : for it.Seek(prefix); it.ValidForPrefix(prefix) && limit > 0; it.Next() {
			offset--;
			for offset >= 0 {
				continue outer
			}

			item := it.Item()
			k := item.Key()
			v, err := item.Value()
			if err == nil {
				records = append(records, Record{Key: k, Value:v})
			}
			limit--
			continue outer

		}

		return nil
	})

	return records, err
}

func (e EngineBadger) Set(key []byte, value []byte, bucket []byte) error {
	return e.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
}

func (e EngineBadger) Delete(key []byte, bucket []byte) error {
	return e.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})

}



