package db

import (
	"reflect"
)

type table struct {
	name     string
	columns  []Column
	database Driver
}

var _ Viewable = Int64{}
var _ Variable = &Int64{}
var _ value = &Int64{}

//Table is a database table definition.
type Table interface {
	Table() string

	Columns() int
	Column(i int) Column

	Database() Driver
}

//Database returns the Driver of this update.
func (u Update) Database() Driver {
	return u.driver
}

//Column is the definition of a column.
type Column interface {
	Column() string
	Database() Driver
	Type() reflect.Type
	Key() bool
	Field() int
}

//Row can return its table definition.
type Row interface {
	Row() Table
}

//View is a type that implements Viewer & Row.
//Embed this in your models and tag it with the name of the table you would like to view.
//ie `db:"users"`
type View struct {
	table

	driver Driver

	index, length int

	vtype reflect.Type

	//The master viewer is readonly. Any attempt to write to master is considered to be an illegal datarace.
	master Viewer
}

//SetDriver implements Viewer.
func (v *View) SetDriver(database Driver) {
	v.driver = database
}

//Table implements Table.
func (v View) Table() string {
	return v.table.name
}

//SetWindow sets the length of elements to view.
func (v *View) SetWindow(length int) {
	v.length = length
}

//Column implements Table.
func (v View) Column(i int) Column {
	return v.columns[i]
}

//Columns implements Table.
func (v View) Columns() int {
	return len(v.columns)
}

//Database implements Table.
func (v View) Database() Driver {
	return v.driver
}

//Row implements Row.
func (v View) Row() Table {
	return v
}

//Setup returns true if this viewer is setup and ready to use.
func (v View) Setup() bool {
	return v.table.name != ""
}

//Master implements Viewer.
func (v *View) Master() bool {
	if v.master == nil {
		return false
	}
	return v.master.address() == v
}

//Viewer is any type that is viewing a table of a specific database.
//Embed a View type to implement Viewer.
type Viewer interface {
	Table

	SetDriver(database Driver)

	Setup() bool
	Master() bool

	address() *View
}

type viewable interface {
	view() View
}

func (v View) view() View {
	return v
}

func (v *View) address() *View {
	return v
}

type value interface {
	Variable

	setprivate(
		table, column string,
		field int,
		key bool,
		driver Driver,
		view Table,
	)
}

//Variable is a variable value in the database.
type Variable interface {
	Viewable

	Pointer() interface{}

	Make(int) interface{}

	Index(int) bool
	Slice(int) interface{}

	Master() bool
}

//Viewable is a value in the database that can be viewed.
type Viewable interface {
	Column

	Interface() interface{}
	Table() string
}
