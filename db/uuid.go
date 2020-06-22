package db

import (
	"github.com/google/uuid"
)

type UUID struct {
	Column

	id    uuid.UUID
	slice []uuid.UUID
}

func (u UUID) Equals(id uuid.UUID) Condition {
	return Condition{
		Column:   u.Column,
		Operator: Equals,
		Value:    id,
	}
}

func (u *UUID) Scan(value interface{}) error {
	return u.id.Scan(value)
}

func (u UUID) UUID() uuid.UUID {
	return u.id
}

func (u UUID) String() string {
	return u.id.String()
}

func (u UUID) On(other UUID) Linker {
	return Linker{
		From: u.Column,
		Join: Inner,
		To:   other.Column,
	}
}

func (u UUID) Interface() interface{} {
	return u.id
}

func (u *UUID) Pointer() interface{} {
	return &u.id
}

func (u *UUID) Slice(length int) interface{} {
	if len(u.slice) != length {
		u.slice = make([]uuid.UUID, length)
	}
	return u.slice
}

func (u *UUID) Index(index int) {
	if index < len(u.slice) && index >= 0 {
		u.id = u.slice[index]
	}
}

func (u *UUID) Set(val uuid.UUID) {
	u.id = val
}

func (u UUID) SetTo(val uuid.UUID) Setting {
	return Setting{
		Column: u.Column,
		Value:  UUID{u.Column, val, nil},
	}
}

func (u UUID) To(val uuid.UUID) Update {
	return Update{
		Column: u.Column,
		Value:  val,
	}
}
