package sql

import (
	"errors"
	"fmt"
	"reflect"
)

//Ready reads matching results into the provided interface values.
//returns the number of rows matched.
func (q *Query) Read(into ...interface{}) (int, error) {
	if len(into) == 0 {
		return 0, nil
	}

	if len(into) > 1 || reflect.TypeOf(into).Kind() != reflect.Slice {
		return 0, errors.New("sql.Query.Read: invalid arguments")
	}
	var first = into[0]
	if _, ok := first.(Row); !ok {
		return 0, errors.New("sql.Query.Read: invalid arguments")
	}

	var head Query

	var T = reflect.TypeOf(first).Elem()
	var value = reflect.ValueOf(first).Elem()

	fmt.Fprintf(&head, "SELECT ")
	for i := 1; i < T.NumField(); i++ {
		fmt.Fprintf(&head, `"%v"`, T.Field(i).Name)

		if i < T.NumField()-1 {
			head.WriteByte(',')
		}
	}
	fmt.Fprintf(&head, ` FROM "%v"`+"\n", first.(Row).rowTable())

	head.WriteQuery(q)

	if q.db.DB == nil {
		return 0, NoDatabase
	}

	result, err := q.db.QueryContext(q.db.Context, head.String(), head.Values...)

	var length = value.Len()

	columns, err := result.Columns()
	if err != nil {
		return 0, err
	}

	for i := 0; i < length; i++ {
		if !result.Next() {
			return i, result.Err()
		}

		var index = value.Index(i)
		var fields = make([]interface{}, len(columns))
		for j, column := range columns {
			for k := 0; k < index.NumField(); k++ {
				if T.Field(k).Name == column {
					fields[j] = index.Field(k).Addr().Interface()
				}
			}
		}

		err := result.Scan(fields...)
		if err != nil {
			return i, err
		}
	}

	return length, result.Close()
}
