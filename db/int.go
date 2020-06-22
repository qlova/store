package db

import "strconv"

type Int struct {
	Column

	int
	slice []int
}

func (i Int) String() string {
	return strconv.Itoa(i.int)
}

func (i Int) Int() int {
	return i.int
}

func (i Int) Interface() interface{} {
	return i.int
}

func (i *Int) Pointer() interface{} {
	return &i.int
}

func (i *Int) Slice(length int) interface{} {
	if len(i.slice) != length {
		i.slice = make([]int, length)
	}
	return i.slice
}

func (i *Int) Index(index int) {
	if index < len(i.slice) && index >= 0 {
		i.int = i.slice[index]
	}
}

func (i *Int) Set(val int) {
	i.int = val
}

func (i Int) SetTo(val int) Setting {
	return Setting{
		Column: i.Column,
		Value:  Int{i.Column, val, nil},
	}
}

func (i Int) Equals(val int) Condition {
	return Condition{
		Column:   i.Column,
		Operator: Equals,
		Value:    val,
	}
}

func (i Int) To(val int) Update {
	return Update{
		Column: i.Column,
		Value:  val,
	}
}
