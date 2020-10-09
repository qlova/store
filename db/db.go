//go:generate go2 tool go2go translate types.go2
//go:generate sed -i "/type STARTMOCK int/,/type ENDMOCK int/d" types.go
//go:generate sed -i "s|//line|//|g" types.go

//Package db provides an abstract database interface for Go.
package db

import (
	"encoding/json"
	"reflect"
)

//Update describes a modification to make to a row in the database.
type Update struct {
	driver        Driver
	Column, Table string
	Value         interface{}

	Then *Update
}

func (u Update) And(other Update) Update {
	u.Then = &other
	return u
}

//Driver is a database driver.
type Driver interface {
	//Connect connects the given viewers to view this database.
	//It then returns the database.
	Connect(Viewer, ...Viewer) Driver

	//Sync syncs the Tables with the Database, adding any missing columns.
	//If constraints or types do not match up, an error is returned.
	Sync(Table, ...Table) error

	//Insert inserts the given row into the database.
	Insert(Row, ...Row) error

	//Delete deletes the given tables.
	Delete(Table, ...Table) error

	//Empty removes all rows from the given tables so that they are empty.
	Empty(Table, ...Table) error

	//Search returns results for the given filter.
	Search(Filter) Results

	//Close closes the connection to the database.
	Close() error
}

//Results of a searchfilter.
type Results interface {
	json.Marshaler

	//Update updates the results with the given updates.
	//Returns the number of results updated (or -1 if the statistic is unavailable).
	Update(Update, ...Update) (int, error)

	//Delete deletes all the results from the database.
	Delete() (int, error)

	//Get gets the matching columns of the results.
	Get(Variable, ...Variable) (int, error)

	//Count returns the number of results.
	Count(Viewable) (int, error)

	//Sum returns the sum amount of the value in the given column of all results.
	Sum(Variable) error

	//Average returns the average value in the given column for all results.
	Average(Viewable) (float64, error)
}

//Open opens a database based on the provided optional arguments.
//The first argument is a database name and subsequent arguments are passed to it.
//If no arguments are provided, a builtin database is used.
func Open(args ...string) Driver {
	if len(args) == 0 {
		return Builtin("")
	}
	return nil
}

type Iterator struct {
	Viewer
	Index int
}

//Next loads viewer's next result into the viewer.
//Returns false when there are no more results.
func (r *Iterator) Next() bool {
	var viewer = r.Viewer

	var RowType = reflect.TypeOf(viewer)
	var RowValue = reflect.ValueOf(viewer).Elem()

	var last = true

	for i := 0; i < RowType.Elem().NumField(); i++ {
		value, ok := RowValue.Field(i).Addr().Interface().(Variable)
		if ok {
			if value.Index(r.Index) {
				last = false
			}
		}
	}

	r.Index++

	return !last
}

//Range returns a new iterator.
func Range(viewer Viewer) *Iterator {
	return &Iterator{
		Viewer: viewer,
		Index:  0,
	}
}
