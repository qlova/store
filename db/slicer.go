package db

import (
	"encoding/json"
	"reflect"
)

type Slicer interface {
	json.Marshaler

	Columns(Column, ...Column) Slicer
	Into(Connectable, ...Connectable) error
}

type Slicing struct {
	slices []reflect.Value
}

func Slice(length int, columns []Column, models []Connectable) Slicing {
	return Slicing{}
}

func Slices(model Connectable, length int, columns ...Column) []reflect.Value {
	var RowValue = reflect.ValueOf(model).Elem()

	var result []reflect.Value

	for _, col := range columns {
		value, ok := RowValue.Field(int(col.Field)).Addr().Interface().(MutableValue)
		if ok {
			result = append(result, reflect.ValueOf(value.Slice(length)))
		}
	}

	model.getModel().State.length = length

	return result
}
