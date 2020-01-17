package sql

import (
	"reflect"
	"strings"
)

//TableInfo is a wrapper around the Table interface that is a useful helper for getting
//metadata about the table.
type TableInfo struct {
	Table
}

//Columns returns a slice of column names that the table has.
func (info TableInfo) Columns() (result []ColumnInfo) {
	var T = reflect.TypeOf(info.Table).Elem()

	for i := 0; i < T.NumField(); i++ {
		var field = T.Field(i)
		if field.Name == "NewTable" {
			continue
		}

		var value = reflect.ValueOf(info.Table).Elem().Field(i)
		if _, ok := value.Interface().(Type); ok {
			value.FieldByName("NewType").Set(reflect.ValueOf(NewType{
				string: field.Name,
			}))
		}

		var constraints = field.Tag.Get("constraints")

		result = append(result,
			ColumnInfo{
				strings.ToLower(field.Name),
				reflect.Zero(field.Type).Interface().(Type).String(),
				strings.Split(constraints, ","),
			})
	}

	return
}

//ColumnInfo provides information about a given column.
type ColumnInfo struct {
	Name, Datatype string
	Constraints    []string
}
