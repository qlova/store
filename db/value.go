package db

type MutableValue interface {
	Value
	SetColumn(to Column)
	Pointer() interface{}

	Slice(int) interface{}
	Index(int)
}

type Value interface {
	GetColumn() Column
	Interface() interface{}
}

var _ Value = new(Int)
var _ Value = new(Int)
