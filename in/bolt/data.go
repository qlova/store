package bolt

import "bytes"
import "io"
import "errors"

type Data struct {
	Store

	Key string
	Tx  *Tx
}

func (data *Data) Exists() bool {
	var tx, _ = data.Tx.BeginRead(data.DB)

	var exists = data.Bucket(tx) != nil

	tx.Rollback()

	return exists
}

func (data *Data) Create() io.WriteCloser {
	var tx, err = data.Tx.BeginWrite(data.DB)
	if err != nil {
		data.err = err
		return nil
	}

	var bucket = data.Bucket(tx)
	if bucket == nil {
		err := tx.Commit()
		if err != nil {
			data.err = err
			return nil
		}
		data.err = errors.New("Bucket does not exist")
		return nil
	}

	bucket.Put([]byte(data.Key), []byte{})

	return WriterCloser{tx, bucket, data.Key}
}
func (data *Data) Delete() error {
	var tx, err = data.Tx.BeginWrite(data.DB)
	if err != nil {
		data.err = err
		return nil
	}

	var bucket = data.Bucket(tx)
	if bucket == nil {
		err := tx.Commit()
		if err != nil {
			data.err = err
			return nil
		}
		data.err = errors.New("Bucket does not exist")
		return nil
	}

	return bucket.Delete([]byte(data.Key))
}

func (data *Data) Reader() io.ReadCloser {
	var tx, err = data.Tx.BeginRead(data.DB)
	if err != nil {
		data.err = err
		return nil
	}

	var bucket = data.Bucket(tx)
	if bucket == nil {
		err := tx.Rollback()
		if err != nil {
			data.err = err
			return nil
		}
		data.err = errors.New("Bucket does not exist")
		return nil
	}

	return ReaderCloser{tx, bytes.NewReader(bucket.Get([]byte(data.Key)))}
}

func (data *Data) Size() int64 {
	var tx, err = data.Tx.BeginRead(data.DB)
	if err != nil {
		data.err = err
		return -1
	}

	var bucket = data.Bucket(tx)
	if bucket == nil {
		err := tx.Rollback()
		if err != nil {
			data.err = err
			return -1
		}
		data.err = errors.New("Bucket does not exist")
		return -1
	}

	var size = len(bucket.Get([]byte(data.Key)))

	tx.Rollback()

	return int64(size)
}

func (data *Data) String() string {

	var result string

	for _, name := range data.Buckets {
		result += name + "/"
	}

	result += data.Key

	return result
}

func (data *Data) Error() error {
	return data.err
}
