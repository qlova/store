package db

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

type Driver interface {
	Verify(Schema) error

	Insert(Row, ...Row) error
	Delete(Table, ...Table) error
	Truncate(Table, ...Table) error

	Querier
}

var connections = make(map[Connection]Driver)
var mutex sync.Mutex

func Open(connection Connection, driver Driver) {
	mutex.Lock()
	defer mutex.Unlock()

	connections[connection] = driver
}

type Schema struct {
	Table
	Columns []Value
}

func register(model Connectable) error {
	var ModelType = reflect.TypeOf(model).Elem()
	var ModelValue = reflect.ValueOf(model).Elem()

	var table Table

	var EmbeddedModelValue reflect.Value

	for i := 0; i < ModelValue.NumField(); i++ {
		var value = ModelValue.Field(i)
		var field = ModelType.Field(i)

		if field.Type == reflect.TypeOf(Model{}) {
			if name, ok := field.Tag.Lookup("db"); ok {
				EmbeddedModelValue = value
				table.Name = name
				break
			}
			return errors.New("sql.RegisterModel: model must have a db tag `db:\"table name\"` ")
		}
	}

	var columns []Value

	for i := 0; i < ModelValue.NumField(); i++ {
		var value = ModelValue.Field(i)
		var field = ModelType.Field(i)

		if setter, ok := value.Addr().Interface().(MutableValue); ok {
			var col Column
			var name = field.Name
			var key bool

			//The name can be overriden in the tag.
			//Further tags include key and unique
			if tag, ok := field.Tag.Lookup("db"); ok {
				args := strings.Split(tag, ",")
				if args[0] != "" {
					name = args[0]
				}
				if len(args) > 0 {
					if args[1] == "key" {
						key = true
					}
				}
			}

			col = Column{
				Name:  name,
				Table: table,
				Field: int16(i),
				Key:   key,
			}
			setter.SetColumn(col)
			columns = append(columns, value.Interface().(Value))
		}
	}

	if table.Name == "" {
		return errors.New("table name not detected")
	}

	EmbeddedModelValue.Set(reflect.ValueOf(Model{
		Table: table,
		State: state{
			parent: ModelType,
			root:   ModelValue,
		},
	}))

	return connections[model.getModel().Connection].Verify(Schema{
		Table:   model.GetTable(),
		Columns: columns,
	})
}

func Register(models ...Connectable) error {
	for _, model := range models {
		if err := register(model); err != nil {
			return err
		}
	}
	return nil
}

//Truncate removes all records from the given table.
func Truncate(table Table) error {
	return connections[table.Connection].Truncate(table)
}

//Insert inserts the given rows into the database.
func Insert(row Row) error {
	table := row.GetTable()
	return connections[table.Connection].Insert(row)
}

func Link(linker Linker, linkers ...Linker) Query {
	return connections[linker.From.Connection].Link(linker, linkers...)
}
