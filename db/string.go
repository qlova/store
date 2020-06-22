package db

import "encoding/json"

type String struct {
	Column

	string
	slice []string
}

func (s String) Interface() interface{} {
	return s.string
}

func (s String) String() string {
	return s.string
}

func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.string)
}

func (s *String) Pointer() interface{} {
	return &s.string
}

func (s *String) Slice(length int) interface{} {
	if len(s.slice) != length {
		s.slice = make([]string, length)
	}
	return s.slice
}

func (s *String) Index(i int) {
	if i < len(s.slice) && i >= 0 {
		s.string = s.slice[i]
	}
}

func (s *String) Set(val string) {
	s.string = val
}

func (s String) To(val string) Update {
	return Update{
		Column: s.Column,
		Value:  val,
	}
}

func (s String) SetTo(val string) Setting {
	return Setting{
		Column: s.Column,
		Value:  String{s.Column, val, nil},
	}
}

func (s String) Equals(val string) Condition {
	return Condition{
		Column:   s.Column,
		Operator: Equals,
		Value:    val,
	}
}

func (s String) IsNot(val string) Condition {
	return Condition{
		Column:   s.Column,
		Operator: NotEquals,
		Value:    val,
	}
}

func (s String) Exists() Condition {
	return s.IsNot("")
}

func (s String) Contains(val string) Condition {
	return Condition{
		Column:   s.Column,
		Operator: Contains,
		Value:    val,
	}
}
