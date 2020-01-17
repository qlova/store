package sql

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

//Result is a query result.
type Result struct {
	q Query
	*sql.Rows
	error error
}

func (result Result) Error() error {
	return result.error
}

//Query returns the query used for this result.
func (result Result) Query() Query {
	return result.q
}

//Scan is like database/sql.Rows.Scan except it automatically calls database/sql.Rows.Next
func (result Result) Scan(values ...interface{}) error {
	result.Rows.Next()
	return result.Rows.Scan(values...)
}

//Read reads the result into the specified slice. Returns the number of rows read.
func (result Result) Read(slice interface{}) (int, error) {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return 0, errors.New("sql.Result.Read: not a slice")
	}

	if result.Error() != nil {
		return 0, result.Error()
	}

	var value = reflect.ValueOf(slice)
	var length = value.Len()

	var columns, err = result.Columns()
	if err != nil {
		return 0, err
	}

	//TODO type check.
	var T = reflect.TypeOf(slice).Elem()

	for i := 0; i < length; i++ {
		if !result.Next() {
			return i, result.Err()
		}

		var index = value.Index(i)
		var fields = make([]interface{}, len(columns))
		for j, column := range columns {
			for k := 0; k < index.NumField(); k++ {
				if strings.ToLower(T.Field(k).Name) == column {
					fields[j] = index.Field(k).Addr().Interface()
				}
			}
		}

		err := result.Scan(fields...)
		if err != nil {
			return i, err
		}
	}

	return length, nil
}
