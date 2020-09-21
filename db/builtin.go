package db

import (
	"reflect"
	"sync"
)

type storage struct {
	rtype reflect.Type
	slice reflect.Value
}

var database = make(map[[2]string]*storage)
var mutex sync.RWMutex

func index(database Builtin, table string) [2]string {
	return [2]string{
		string(database),
		table,
	}
}

//Builtin is a builtin database.
type Builtin string

func (b Builtin) connect(v Viewer) {
	Connect(v, b)
}

//Connect connects the given viewer to view this database.
//It then returns the database.
func (b Builtin) Connect(first Viewer, more ...Viewer) Driver {
	b.connect(first)
	for _, viewer := range more {
		b.connect(viewer)
	}
	return b
}

func (b Builtin) sync(table Table) error {
	//Need to create a struct that represents this table.
	var fields = make([]reflect.StructField, table.Columns())

	for i := 0; i < table.Columns(); i++ {
		column := table.Column(i)

		fields[i] = reflect.StructField{
			Name: column.Column(),
			Type: column.Type(),
		}
	}

	database[index(b, table.Table())] = &storage{
		rtype: reflect.StructOf(fields),
		slice: reflect.New(reflect.SliceOf(reflect.StructOf(fields))).Elem(),
	}

	return nil
}

//Sync syncs the Tables with the Database, adding any missing columns.
//If constraints or types do not match up, an error is returned.
func (b Builtin) Sync(table Table, tables ...Table) error {
	if err := b.sync(table); err != nil {
		return err
	}
	for _, table := range tables {
		if err := b.sync(table); err != nil {
			return err
		}
	}
	return nil
}

func (b Builtin) insert(row Row) error {
	mutex.Lock()
	defer mutex.Unlock()

	var in Insertion
	in.Row(row)

	var table = database[index(b, row.Row().Table())]

	if table == nil {
		return ErrTableNotFound
	}

	var structure = reflect.New(table.rtype).Elem()

	for i, column := range in.Columns {

		if in.Uniques[i] {

			//Check if the unique value is taken. If, so reject this insert.
			for i := 0; i < table.slice.Len(); i++ {
				row := table.slice.Index(i)

				if reflect.DeepEqual(row.FieldByName(column).Interface(), in.Values[i]) {
					return ErrDuplicateKey
				}
			}
		}

		structure.FieldByName(column).Set(reflect.ValueOf(in.Values[i]))
	}

	table.slice.Set(reflect.Append(table.slice, structure))

	return nil
}

//Insert inserts the given row into the database.
func (b Builtin) Insert(row Row, rows ...Row) error {
	if row != nil {
		if err := b.insert(row); err != nil {
			return err
		}
	}
	for _, row := range rows {
		if err := b.insert(row); err != nil {
			return err
		}
	}
	return nil
}

//Delete deletes the given tables.
func (b Builtin) Delete(table Table, tables ...Table) error {
	if table != nil {
		delete(database, index(b, table.Table()))
	}
	for _, table := range tables {
		delete(database, index(b, table.Table()))
	}
	return nil
}

func (b Builtin) empty(t Table) error {
	mutex.Lock()
	defer mutex.Unlock()

	var table = database[index(b, t.Table())]

	if table == nil {
		return ErrTableNotFound
	}

	table.slice.Set(reflect.Zero(table.slice.Type()))

	return nil
}

//Empty removes all rows from the given tables so that they are empty.
func (b Builtin) Empty(table Table, tables ...Table) error {
	if table != nil {
		if err := b.empty(table); err != nil {
			return err
		}
	}
	for _, table := range tables {
		if err := b.empty(table); err != nil {
			return err
		}
	}
	return nil
}

//Close closes the connection to the database.
func (Builtin) Close() error {
	return nil
}
