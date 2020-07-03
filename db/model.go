package db

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
)

type Connectable interface {
	Row

	getModel() *Model
	Slice(int, int)
}

type Row interface {
	GetTable() Table
}

type Setting struct {
	Column
	Value interface{}
}

func (setting Setting) applyTo(value reflect.Value) {
	value.Field(int(setting.Field)).Set(reflect.ValueOf(setting.Value))
}

type state struct {
	parent reflect.Type
	root   reflect.Value

	key Column

	index, length int
}

func (state) String() string {
	return "Model"
}

type Model struct {
	Table
	State state
}

func (model Model) JSONEncoder() func(interface{}) ([]byte, error) {
	return func(i interface{}) ([]byte, error) {
		var c = i.(Connectable)

		var RowType = reflect.TypeOf(c).Elem()
		var RowValue = reflect.ValueOf(c).Elem()

		var buffer bytes.Buffer

		buffer.WriteByte('[')

		if model.State.length > 0 {
			var length = len(model.Range())
			for i := range model.Range() {
				Next(c)

				buffer.WriteByte('{')

				var first = true
				for i := 0; i < RowType.NumField(); i++ {
					value, ok := RowValue.Field(i).Interface().(Value)
					if ok {
						var val = value.Interface()
						if val == reflect.Zero(reflect.TypeOf(val)).Interface() {
							continue
						}

						if !first {
							buffer.WriteByte(',')
						}

						buffer.WriteString(strconv.Quote(value.GetColumn().Name))
						buffer.WriteByte(':')

						js, err := json.Marshal(val)
						if err != nil {
							return nil, err
						}

						buffer.Write(js)

						first = false
					}
				}

				buffer.WriteByte('}')

				if i < length-1 {
					buffer.WriteByte(',')
				}
			}
		}

		buffer.WriteByte(']')

		return buffer.Bytes(), nil
	}
}

type Range = []struct{}

func (r *Model) Range() Range {
	r.State.index = 0
	return make(Range, r.State.length)
}

func (r *Model) Slice(index, length int) {
	r.State.index = index
	r.State.length = length
}

func (r Model) GetTable() Table {
	return r.Table
}

func (model *Model) getModel() *Model {
	return model
}

func (model Model) With(setting Setting, settings ...Setting) Row {
	var value = reflect.New(model.State.parent)
	value.Elem().Set(model.State.root)

	setting.applyTo(value.Elem())
	for _, setting := range settings {
		setting.applyTo(value.Elem())
	}

	var i = value.Interface()
	*(i.(Connectable).getModel()) = model

	return value.Elem().Interface().(Row)
}

func Next(mut Connectable) {
	var RowType = reflect.TypeOf(mut)
	var RowValue = reflect.ValueOf(mut).Elem()

	var model = mut.getModel()

	for i := 0; i < RowType.Elem().NumField(); i++ {
		value, ok := RowValue.Field(i).Addr().Interface().(MutableValue)
		if ok {
			value.Index(model.State.index)
		}
	}

	model.State.index++
}

func Columns(row Row) []Column {
	var RowType = reflect.TypeOf(row)
	var RowValue = reflect.ValueOf(row).Elem()

	var result []Column

	for i := 0; i < RowType.Elem().NumField(); i++ {
		value, ok := RowValue.Field(i).Interface().(Value)
		if ok {
			result = append(result, value.GetColumn())
		}
	}

	return result
}
