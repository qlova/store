package bolt

import (
	"bytes"
	"io"
	"path"
)

//Key is a bolt Key.
type Key struct {
	DB
}

//Available returns if the file is available for reading.
func (key Key) Available() bool {
	var tx, err = key.Begin(false)
	if err != nil {
		return false
	}

	var bucket = Bucket{DB: key.DB}.For(tx)

	var exists = bucket.Get([]byte(key.Key)) != nil

	if err := tx.Rollback(); err != nil {
		return false
	}

	return exists
}

//CopyTo copies a key to the specified io.Writer
func (key Key) CopyTo(writer io.Writer) error {
	var tx, err = key.Begin(false)
	if err != nil {
		return err
	}

	var bucket = Bucket{DB: key.DB}.For(tx)

	var data = bucket.Get([]byte(key.Key))
	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return tx.Rollback()
}

//From creates and writes a key from the specified reader.
func (key Key) From(reader io.Reader) error {
	var tx, err = key.Begin(true)
	if err != nil {
		return err
	}

	var bucket = Bucket{DB: key.DB}.For(tx)

	var buffer bytes.Buffer

	_, err = io.Copy(&buffer, reader)
	if err != nil {
		return err
	}

	err = bucket.Put([]byte(key.Key), buffer.Bytes())
	if err != nil {
		return err
	}

	return tx.Commit()
}

//Path returns a path representing the objects absolute location.
func (key Key) Path() string {
	return path.Join(Bucket{key.DB}.Path(), key.Key)
}

//Delete the object.
func (key Key) Delete() error {
	var tx, err = key.Begin(true)
	if err != nil {
		return err
	}

	var bucket = Bucket{DB: key.DB}.For(tx)

	err = bucket.Delete([]byte(key.Key))
	if err != nil {
		return err
	}

	return tx.Commit()
}

//Size returns the size of the object.
func (key Key) Size() int64 {
	var tx, err = key.Begin(false)
	if err != nil {
		return -1
	}

	var bucket = Bucket{DB: key.DB}.For(tx)

	data := bucket.Get([]byte(key.Key))
	var length = len(data)

	if err := tx.Rollback(); err != nil {
		return -1
	}

	return int64(length)
}
