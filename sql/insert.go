package sql

import (
	"fmt"
	"reflect"
)

//Insert inserts the given row into the database table specified by its model.
func (db Database) Insert(row Row) error {
	if row.rowTable() == "" {
		fmt.Println("metadata")
		row = register(row)
	}

	var RowType = reflect.TypeOf(row)
	var RowValue = reflect.ValueOf(row)

	var head Query
	var tail Query

	fmt.Fprintf(&head, `INSERT INTO "%v" (`, row.rowTable())
	fmt.Fprint(&tail, "\n)\nVALUES (")

	for i := 0; i < RowType.NumField(); i++ {
		var val, ok = RowValue.Field(i).Interface().(Value)
		if ok {
			fmt.Fprintf(&head, "\n\t"+`"%v"`, val.GetColumn().Name)
			if i < RowType.NumField()-1 {
				head.WriteByte(',')
			}
			tail.WriteString("\n\t")
			fmt.Println(val.Interface())
			tail.WriteString(tail.value(val.Interface()))
			if i < RowType.NumField()-1 {
				tail.WriteByte(',')
			}
		}
	}

	tail.WriteByte('\n')
	tail.WriteByte(')')

	head.WriteQuery(&tail)

	_, err := db.ExecContext(db.Context, head.String(), head.Values...)

	fmt.Println(head.String())

	return err
}
