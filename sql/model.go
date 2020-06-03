package sql

import (
	"reflect"
)

//Model is any struct with an embedded Table.
type Model interface {
	getTable() Table
}

func register(row Row) Row {
	var model = reflect.New(reflect.TypeOf(row))
	RegisterModel(model.Interface().(Model))
	return model.Elem().Interface().(Row)
}

//RegisterModel registers the given model.
//You must pass a pointer to a model variable otherwise this function will panic.
//It is safe to call this multiple times on the same model.
//panics if the model does not have a name tag on its table field.
//RegisterModel should be called before the model is used in order to increase performance.
func RegisterModel(model Model) {
	var ModelType = reflect.TypeOf(model).Elem()
	var ModelValue = reflect.ValueOf(model).Elem()

	var TableName string

	for i := 0; i < ModelValue.NumField(); i++ {
		var value = ModelValue.Field(i)
		var field = ModelType.Field(i)

		if field.Type == reflect.TypeOf(Table("")) {
			if name, ok := field.Tag.Lookup("name"); ok {
				value.Set(reflect.ValueOf(Table(name)))
				TableName = name
				break
			} else {
				panic("sql.RegisterModel: model must have a name tag `name:\"name\"` ")
			}
		}
	}

	for i := 0; i < ModelValue.NumField(); i++ {
		var value = ModelValue.Field(i)
		var field = ModelType.Field(i)

		if setter, ok := value.Addr().Interface().(settableColumn); ok {
			if name, ok := field.Tag.Lookup("name"); ok {
				setter.setColumn(Column{Name: name, Table: Table(TableName)})
			} else {
				setter.setColumn(Column{Name: field.Name, Table: Table(TableName)})
			}
		}
	}
}

//UpdateModel creates, ensures and updates the database to match the Go model.
//this operation will not delete any column, only add missing ones.
//any inconistencies that cannot be resolved will result in an error.
func (db Database) UpdateModel(model Model) error {
	RegisterModel(model)

	return db.CreateTableIfNotExists(model)
}
