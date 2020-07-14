package bolt

import (
	"github.com/boltdb/bolt"
	"qlova.store"
)

//Bucket is a bolt bucket.
type Bucket struct {
	DB
}

//For returns the bolt.Bucket for a given Tx
func (bucket Bucket) For(tx *bolt.Tx) *bolt.Bucket {
	var root = tx.Bucket([]byte("root"))

	for _, name := range bucket.Buckets {
		root = root.Bucket([]byte(name))
		if root == nil {
			return nil
		}
	}

	return root
}

//Available returns if the bucket is available for modifications.
func (bucket Bucket) Available() bool {
	var tx, err = bucket.Begin(false)
	if err != nil {
		return false
	}

	var exists = bucket.For(tx) != nil

	if err := tx.Rollback(); err != nil {
		return false
	}

	return exists
}

//Children returns a list of buckets and keys.
func (bucket Bucket) Children(amount ...int) store.Children {
	var list List
	list.DB = bucket.DB
	if len(amount) > 0 {
		var number = int64(amount[0])
		list.amount = number
	} else {
		list.amount = -1
	}
	return &list
}

//Data returns the key with the given name as a store.Data.
func (bucket Bucket) Data(name string) store.Data {

	var key = Key{bucket.DB}
	key.Buckets = make([]string, len(bucket.Buckets))
	key.Key = name

	copy(key.Buckets, bucket.Buckets)

	return key
}

//Goto navigates to the relative bucket by path.
func (bucket Bucket) Goto(location string) store.Node {
	bucket.Buckets = append(bucket.Buckets, location)
	return bucket
}

//Name returns the name of this folder.
func (bucket Bucket) Name() string {
	return bucket.Buckets[len(bucket.Buckets)-1]
}

//Parent returns the parent bucket of this bucket.
func (bucket Bucket) Parent() store.Node {
	bucket.Buckets = bucket.Buckets[:len(bucket.Buckets)-1]
	return bucket
}

//Path returns a path representing the bucket's absolute location.
func (bucket Bucket) Path() string {
	var result = "/"
	for i, name := range bucket.Buckets {
		result += name
		if i < len(bucket.Buckets)-1 {
			result += "/"
		}
	}
	return result
}

//Create creates the current bucket.
func (bucket Bucket) Create() error {
	var tx, err = bucket.Begin(true)
	if err != nil {
		return err
	}

	var root = tx.Bucket([]byte("root"))

	for _, name := range bucket.Buckets {
		root, err = root.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

//Delete deletes the bucket.
func (bucket Bucket) Delete() error {
	var tx, err = bucket.Begin(true)
	if err != nil {
		return err
	}

	if err := bucket.Parent().(Bucket).For(tx).DeleteBucket([]byte(bucket.Name())); err != nil {
		return err
	}

	return tx.Commit()
}
