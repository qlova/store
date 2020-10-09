package db

import (
	"github.com/google/uuid"
)

func insert(row Row) error {
	if row.Row().Database() == nil {
		return ErrDisconnectedViewer
	}
	return row.Row().Database().Insert(row)
}

//Insert inserts the given rows into their registered databases.
func Insert(first Row, rows ...Row) error {
	if err := insert(first); err != nil {
		return err
	}
	for _, row := range rows {
		if err := insert(row); err != nil {
			return err
		}
	}
	return nil
}

//Insertion describes an insertion operation into the database.
type Insertion struct {
	Table
	Columns []string
	Uniques []bool
	Values  []interface{}
}

//Row makes the insetion operation insert the given row.
func (insert *Insertion) Row(row Row) error {
	insert.Table = row.Row()

	for i := 0; i < insert.Table.Columns(); i++ {
		col := insert.Table.Column(i)

		insert.Columns = append(insert.Columns, col.Column())
		insert.Uniques = append(insert.Uniques, col.Key())

		id, ok := col.(UUID)
		if !ok {
			var uid *UUID
			uid, ok = col.(*UUID)
			if ok {
				id = *uid
			}
		}

		//Hacky feature that can automatically generate a UUID if it is tagged as a key.
		if ok && id.Value() == uuid.Nil {
			if id.Key() {
				id, err := uuid.NewRandom()
				if err != nil {
					return err
				}

				insert.Values = append(insert.Values, id)
				continue
			}
		}

		insert.Values = append(insert.Values, LookAt(row, col).Interface())
	}

	return nil
}
