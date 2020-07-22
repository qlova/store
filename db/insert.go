package db

import (
	"reflect"

	"github.com/google/uuid"
)

type Insertion struct {
	Table
	Columns []Column
	Values  []interface{}
}

func (insert *Insertion) Row(row Row) error {
	var RowType = reflect.TypeOf(row)
	var RowValue = reflect.ValueOf(row)

	insert.Table = row.GetTable()

	for i := 0; i < RowType.NumField(); i++ {
		value, ok := RowValue.Field(i).Interface().(Value)
		if ok {
			insert.Columns = append(insert.Columns, value.GetColumn())

			//Hacky feature that can automatically generate a UUID if it is tagged as a key.
			if id, ok := value.(UUID); ok && id.UUID() == uuid.Nil {
				if id.Key {
					id, err := uuid.NewRandom()
					if err != nil {
						return err
					}

					insert.Values = append(insert.Values, id)
					continue
				}
			}

			insert.Values = append(insert.Values, value.Interface())
		}
	}

	return nil
}
