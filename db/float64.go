package db

import "fmt"

type Float64 struct {
	Column

	float64
	slice []float64
}

func (f Float64) String() string {
	return fmt.Sprintf("%f", f.float64)
}

func (f Float64) Float64() float64 {
	return f.float64
}

func (f Float64) Interface() interface{} {
	return f.float64
}

func (f *Float64) Pointer() interface{} {
	return &f.float64
}

func (f *Float64) Slice(length int) interface{} {
	if len(f.slice) != length {
		f.slice = make([]float64, length)
	}
	return f.slice
}

func (f *Float64) Index(index int) {
	if index < len(f.slice) && index >= 0 {
		f.float64 = f.slice[index]
	}
}

func (f *Float64) Set(val float64) {
	f.float64 = val
}

func (f *Float64) SetTo(val float64) Setting {
	return Setting{
		Column: f.Column,
		Value:  Float64{f.Column, val, nil},
	}
}

func (f *Float64) Equals(val float64) Condition {
	return Condition{
		Column:   f.Column,
		Operator: Equals,
		Value:    val,
	}
}

func (f Float64) To(val float64) Update {
	return Update{
		Column: f.Column,
		Value:  val,
	}
}
