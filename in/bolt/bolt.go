package bolt

import (
	"github.com/boltdb/bolt"
	"qlova.store"
)

//DB is the stored bolt.DB state for the types in this package.
type DB struct {
	*bolt.DB

	Buckets []string
	Key     string
}

//Database is a boltdb database.
type Database struct {
	Bucket
}

// Open the db file in your current directory.
// It will be created if it doesn't exist.
func Open(name string) (store.Root, error) {
	db, err := bolt.Open(name, 0700, nil)
	if err != nil {
		return store.Root{}, err
	}

	tx, err := db.Begin(true)
	if err != nil {
		return store.Root{}, err
	}

	_, err = tx.CreateBucketIfNotExists([]byte("root"))
	if err != nil {
		return store.Root{}, err
	}

	err = tx.Commit()
	if err != nil {
		return store.Root{}, err
	}

	//wew messy
	return store.Root{Object: store.Object{
		Node: Database{Bucket{
			DB: DB{
				DB: db,
			},
		}},
	}}, nil
}
