package bolt

import "bytes"
import "github.com/boltdb/bolt"

type Tx struct {
	*bolt.Tx
	Level    int
	ReadOnly bool
}

func (tx *Tx) BeginWrite(db *bolt.DB) (*Tx, error) {
	if tx == nil {
		var tx, err = db.Begin(true)
		return &Tx{tx, 0, false}, err
	} else {

		if tx.ReadOnly {
			err := tx.Rollback()
			if err != nil {
				return nil, err
			}

			tx, err := db.Begin(true)
			return &Tx{tx, 0, false}, err
		}

		tx.Level++
		return tx, nil
	}
}

func (tx *Tx) BeginRead(db *bolt.DB) (*Tx, error) {
	if tx == nil {
		var tx, err = db.Begin(false)
		return &Tx{tx, 0, true}, err
	} else {
		tx.Level++
		return tx, nil
	}
}

func (tx *Tx) Commit() error {
	if tx.Level <= 0 {
		if tx.ReadOnly {
			return tx.Tx.Rollback()
		} else {
			return tx.Tx.Commit()
		}
	} else {
		tx.Level--
		return nil
	}
}

func (tx *Tx) Rollback() error {
	if tx.Level <= 0 {
		if tx.ReadOnly {
			return tx.Tx.Rollback()
		} else {
			return tx.Tx.Commit()
		}
	} else {
		tx.Level--
		return nil
	}
}

type WriterCloser struct {
	*Tx
	Bucket *bolt.Bucket
	Key    string
}

func (wc WriterCloser) Write(data []byte) (int, error) {
	err := wc.Bucket.Put([]byte(wc.Key), append(wc.Bucket.Get([]byte(wc.Key)), data...))
	if err != nil {
		return 0, err
	} else {
		return len(data), nil
	}
}

func (wc WriterCloser) Close() error {
	return wc.Commit()
}

type ReaderCloser struct {
	*Tx
	reader *bytes.Reader
}

func (rc ReaderCloser) Read(buffer []byte) (int, error) {
	return rc.reader.Read(buffer)
}

func (rc ReaderCloser) Close() error {
	return rc.Rollback()
}
