package bolt

import "os"
import "errors"
import "github.com/boltdb/bolt"
import qlovastore "github.com/qlova/store"

type Store struct {
	DB *bolt.DB

	err error

	Buckets []string
}

func (store Store) Error() error {
	return store.err
}

func (store Store) Bucket(tx *Tx) *bolt.Bucket {
	var root = tx.Bucket([]byte("root"))

	for _, name := range store.Buckets {
		root = root.Bucket([]byte(name))
		if root == nil {
			return nil
		}
	}

	return root
}

func (store Store) Exists() bool {
	var tx, err = store.DB.Begin(false)
	if err != nil {
		return false
	}
	defer tx.Rollback()

	var root = tx.Bucket([]byte("root"))

	for _, name := range store.Buckets {
		root = root.Bucket([]byte(name))
		if root == nil {
			return false
		}
	}

	return true
}

func (store Store) Goto(path string) qlovastore.Store {
	return Store{
		DB:      store.DB,
		Buckets: append(store.Buckets, path),
	}
}

func (store Store) String() string {
	var result string = "/"
	for i, name := range store.Buckets {
		result += name
		if i < len(store.Buckets)-1 {
			result += "/"
		}
	}
	return result
}

func (store Store) List() qlovastore.Cursor {
	return &Cursor{}
}

func (store Store) Create() error {
	var tx, err = store.DB.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Commit()

	var root = tx.Bucket([]byte("root"))

	for _, name := range store.Buckets {
		root, err = root.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
	}

	return nil
}

func (store Store) Delete() error {
	return errors.New("Unimplemented")
}

func (store Store) Data(name string) qlovastore.Data {

	var data Data
	data.DB = store.DB
	data.Buckets = make([]string, len(store.Buckets))
	data.Key = name

	copy(data.Buckets, store.Buckets)

	return qlovastore.Data{
		Raw: &data,
	}
}

// Open the db file in your current directory.
// It will be created if it doesn't exist.
func Open(name string, mode os.FileMode, options *bolt.Options) Store {
	db, err := bolt.Open(name, mode, options)

	if err == nil {
		var tx, err = db.Begin(true)
		if err == nil {
			_, err = tx.CreateBucketIfNotExists([]byte("root"))
			if err == nil {
				err = tx.Commit()
			}
		}
	}

	return Store{
		DB:  db,
		err: err,
	}
}
