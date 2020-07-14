package bolt

import (
	"io"
	"path"

	"qlova.store"
)

type keyInfo struct {
	Key      string
	IsBucket bool
}

//List is a list of folders and objects.
type List struct {
	DB

	amount int64
	index  int

	current string
	keys    []keyInfo

	isBucket bool
}

//Name returns the current item's name or blank if the current item is empty.
func (list *List) Name() string {
	return list.current
}

//Data returns the current item as data or nil if the current item is not data.
func (list *List) Data() store.Data {
	if !list.isBucket {
		var key = Key{DB: list.DB}
		key.Key = list.current
		return key
	}
	return nil
}

//Node returns the current item as a node or nil if the current item is not a node.
func (list *List) Node() store.Node {
	if list.isBucket {
		var bucket = Bucket{DB: list.DB}
		bucket.Goto(list.current)
		return bucket
	}
	return nil
}

//Reference is a reference to a child.
func (list *List) Reference() store.Reference {
	return store.Reference{
		Package:  "os",
		Internal: list.current,
		Fallback: path.Join(list.Path(), list.current),
	}
}

//SkipTo skips to the child with name 'name'.
func (list *List) SkipTo(ref store.Reference) error {
	return list.next(&ref.Internal)
}

//Next moves to the next item.
func (list *List) Next() error {
	return list.next(nil)
}

func (list *List) next(ref *string) error {
	if list.keys == nil || list.index >= len(list.keys) {
		list.index = 0

		var tx, err = list.Begin(false)
		if err != nil {
			return err
		}

		var cursor = Bucket{list.DB}.For(tx).Cursor()
		if err != nil {
			return err
		}

		if list.keys == nil {
			key, value := cursor.First()
			if key != nil {
				list.keys = append(list.keys, keyInfo{
					Key:      string(key),
					IsBucket: value == nil,
				})
			}
		} else {
			list.keys = nil
		}

		if list.current != "" {
			cursor.Seek([]byte(list.current))
		}

		if ref != nil {
			cursor.Seek([]byte(*ref))
		}

		var unwind = list.amount
		for i := unwind; i > 0 || i == -1; i-- {
			key, value := cursor.Next()

			if key == nil && value == nil {
				break
			}
			list.keys = append(list.keys, keyInfo{
				Key:      string(key),
				IsBucket: value == nil,
			})
		}

		if err = tx.Rollback(); err != nil {
			return err
		}

	}

	if len(list.keys) == 0 {
		return io.EOF
	}

	list.current = list.keys[list.index].Key
	list.isBucket = list.keys[list.index].IsBucket
	list.index++

	return nil
}
