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

			if id, ok := value.(UUID); ok && id.UUID() == uuid.Nil {
				id, err := uuid.NewRandom()
				if err != nil {
					return err
				}

				insert.Values = append(insert.Values, id)
				continue
			}

			insert.Values = append(insert.Values, value.Interface())
		}
	}

	return nil
}
