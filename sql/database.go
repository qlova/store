package sql

import (
	"fmt"
	"reflect"
)

//CreateTable creates the given table.
func (db Database) CreateTable(table Table) Query {
	return db.createTable(table, "CREATE TABLE")
}

//EnsureTable is shorthand for CreateTableIfNotExists.
func (db Database) EnsureTable(table Table) Query {
	return db.CreateTableIfNotExists(table)
}

//CreateTableIfNotExists creates and returns a database if it doesn't exist.
func (db Database) CreateTableIfNotExists(table Table) Query {
	return db.createTable(table, "CREATE TABLE IF NOT EXISTS")
}

func (db Database) createTable(table Table, header string) Query {
	var T = reflect.TypeOf(table).Elem()
	var field, ok = T.FieldByName("NewTable")
	if !ok {
		return db.QueryError("sql: sql.NewTable must be the embedded within the table type")
	}

	var name = field.Tag.Get("name")
	if name == "" {
		return db.QueryError("sql: sql.NewTable must have a nametag `name:\"name\"` ")
	}

	(reflect.ValueOf(table).Elem().FieldByName("NewTable").
		Addr().Interface().(*NewTable)).set(db, name)

	var query = db.NewQuery()
	fmt.Fprintf(query, `%v %v (`, header, name)

	for i := 0; i < T.NumField(); i++ {
		var field = T.Field(i)
		if field.Name == "NewTable" {
			continue
		}

		var t = "string"
		switch field.Type {
		case reflect.TypeOf(Int{}):
			t = "int"
		case reflect.TypeOf(String{}):
			t = "varchar(255)"
		}

		var name = field.Name

		var value = reflect.ValueOf(table).Elem().Field(i)
		if _, ok := value.Interface().(Type); ok {
			value.FieldByName("NewType").Set(reflect.ValueOf(NewType{
				string: name,
			}))
		}

		fmt.Fprintf(query, "\n\t%v %v", field.Name, t)
		if i < T.NumField()-1 {
			query.WriteByte(',')
		}
	}

	query.WriteByte('\n')
	query.WriteByte(')')

	return query
}
