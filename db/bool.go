package db

type Bool struct {
	Column

	bool
	slice []bool
}

func (b Bool) Bool() bool {
	return b.bool
}

func (b Bool) Equals(val bool) Condition {
	return Condition{
		Column:   b.Column,
		Operator: Equals,
		Value:    val,
	}
}

func (b Bool) String() string {
	if b.bool {
		return "true"
	}
	return "false"
}

func (b Bool) Interface() interface{} {
	return b.bool
}

func (b *Bool) Pointer() interface{} {
	return &b.bool
}

func (b *Bool) Slice(length int) interface{} {
	if len(b.slice) != length {
		b.slice = make([]bool, length)
	}
	return b.slice
}

func (b *Bool) Index(index int) {
	if index < len(b.slice) && index >= 0 {
		b.bool = b.slice[index]
	}
}

func (b Bool) SetTo(val bool) Setting {
	return Setting{
		Column: b.Column,
		Value:  Bool{b.Column, val, nil},
	}
}

func (b Bool) To(val bool) Update {
	return Update{
		Column: b.Column,
		Value:  val,
	}
}
