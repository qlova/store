package db

import "time"

type Time struct {
	Column

	time  time.Time
	slice []time.Time
}

func (t Time) String() string {
	return t.time.String()
}

func (t Time) Time() time.Time {
	return t.time
}

func (t Time) Interface() interface{} {
	return t.time.In(time.UTC)
}

func (t *Time) Pointer() interface{} {
	return &t.time
}

func (t *Time) Slice(length int) interface{} {
	if len(t.slice) != length {
		t.slice = make([]time.Time, length)
	}
	return t.slice
}

func (t *Time) Index(index int) {
	if index < len(t.slice) && index >= 0 {
		t.time = t.slice[index]
	}
}

func (t Time) SetTo(val time.Time) Setting {
	return Setting{
		Column: t.Column,
		Value:  Time{t.Column, val, nil},
	}
}

func (t Time) To(val time.Time) Update {
	return Update{
		Column: t.Column,
		Value:  val.In(time.UTC),
	}
}
