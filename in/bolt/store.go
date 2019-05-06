package bolt

import "os"
import "github.com/boltdb/bolt"
import qlovastore "github.com/qlova/store"

type Store struct {
	DB *bolt.DB
	
	err error
	
	Buckets []string
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
		DB: db,
		err: err,
	}
}
