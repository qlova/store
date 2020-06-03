package sql

import (
	"fmt"
	"reflect"
	"strings"
)

//Table is a reference to a database table.
type Table string

//getTable implements Model
func (t *Table) getTable() Table {
	return *t
}

//rowTable implements Row
func (t Table) rowTable() Table {
	return t
}

type Row interface {
	rowTable() Table
}

func constraints(field reflect.StructField) []string {
	return strings.Split(field.Tag.Get("constraints"), ",")
}

//CreateTable creates a table in the database from the given model.
func (db Database) CreateTable(model Model) error {
	return db.createTable("CREATE TABLE", model)
}

//CreateTableIfNotExists creates a table in the database from the given model if it doesn't already exist.
func (db Database) CreateTableIfNotExists(model Model) error {
	return db.createTable("CREATE TABLE IF NOT EXISTS", model)
}

func (db Database) createTable(mode string, model Model) error {
	var ModelType = reflect.TypeOf(model).Elem()
	var ModelValue = reflect.ValueOf(model).Elem()
	var TableName = model.getTable()

	var q Query
	fmt.Fprintf(&q, mode+` "%v" (`, TableName)

	for i := 0; i < ModelType.NumField(); i++ {
		var field = ModelType.Field(i)

		var zero = ModelValue.Field(i).Interface()

		var column Column
		var t Type

		if getter, ok := zero.(HasColumn); ok {
			column = getter.GetColumn()
		}
		if getter, ok := zero.(HasType); ok {
			t = getter.GetType()
		}

		if column.Name != "" && t != "" {
			fmt.Fprintf(&q, "\n\t"+`"%v" %v %v`, column.Name, t, strings.Join(constraints(field), " "))
			if i < ModelType.NumField()-1 {
				q.WriteByte(',')
			}
		}

	}

	q.WriteByte('\n')
	q.WriteByte(')')

	_, err := db.ExecContext(db.Context, q.String())

	return err
}

//DeleteTable permanantly deletes the data in the model's table.
func (db Database) DeleteTable(model Model) error {
	var q Query
	fmt.Fprintf(&q, `TRUNCATE TABLE "%v"`+"\n", model.getTable())

	_, err := db.ExecContext(db.Context, q.String())
	return err
}

//Insert inserts a struct into this table.
/*func (table *NewTable) Insert(structure interface{}) Query {
	var T = reflect.TypeOf(structure)
	if T.Kind() != reflect.Struct {
		return table.db.QueryError("Table.Insert: must be struct")
	}

	var value = reflect.ValueOf(structure)

	var query = table.db.NewQuery()
	var tail = table.db.NewQuery()

	fmt.Fprintf(query, `INSERT INTO "%v" (`, table.name)
	fmt.Fprint(tail, "\n)\nVALUES (")

	for i := 0; i < T.NumField(); i++ {
		fmt.Fprintf(query, "\n\t"+`"%v"`, T.Field(i).Name)
		if i < T.NumField()-1 {
			query.WriteByte(',')
		}
		fmt.Fprintf(tail, "\n\t%v", query.value(value.Field(i).Interface()))
		if i < T.NumField()-1 {
			tail.WriteByte(',')
		}
	}

	tail.WriteByte('\n')
	tail.WriteByte(')')

	query.Write(tail.Bytes())

	return query
}

//InsertValues inserts a new set of values into this table.
func (table *NewTable) InsertValues(values Values) Query {
	var query = table.db.NewQuery()
	var tail = table.db.NewQuery()

	fmt.Fprintf(query, `INSERT INTO "%v" (`, table.name)
	fmt.Fprint(tail, "\n)\nVALUES (")

	for i, value := range values {
		fmt.Fprintf(query, "\n\t"+`"%v"`, value.Column())
		if i < len(values)-1 {
			query.WriteByte(',')
		}
		fmt.Fprintf(tail, "\n\t%v", value.get(query))
		if i < len(values)-1 {
			tail.WriteByte(',')
		}
	}

	tail.WriteByte('\n')
	tail.WriteByte(')')

	query.Write(tail.Bytes())

	return query
}

//Update updates the filtered records to equal this struct,
func (table *NewTable) Update(structure interface{}) Query {
	var T = reflect.TypeOf(structure)
	if T.Kind() != reflect.Struct {
		return table.db.QueryError("Table.Update: must be struct")
	}

	var value = reflect.ValueOf(structure)

	var query = table.db.NewQuery()

	fmt.Fprintf(query, `UPDATE "%v" SET `, table.name)

	for i := 0; i < T.NumField(); i++ {
		fmt.Fprintf(query, "\n\t"+`"%v" = %v`, T.Field(i).Name, query.value(value.Field(i).Interface()))
		if i < T.NumField()-1 {
			query.WriteByte(',')
		}
	}
	query.WriteByte('\n')

	return query
}

//UpdateValues updates the values of the filtered records.
func (table *NewTable) UpdateValues(values Values) Query {
	var query = table.db.NewQuery()

	fmt.Fprintf(query, `UPDATE "%v" SET `, table.name)

	for i, value := range values {
		fmt.Fprintf(query, "\n\t"+`"%v" = %v`, value.Column(), value.get(query))
		if i < len(values)-1 {
			query.WriteByte(',')
		}
	}

	fmt.Fprintf(query, "\n\t")

	return query
}

func (table *NewTable) selectType(header string, columns ...Column) Query {
	var query = table.db.NewQuery()

	if len(columns) == 0 {
		var T = reflect.TypeOf(table.structure).Elem()
		fmt.Fprintf(query, "%v ", header)
		for i := 1; i < T.NumField(); i++ {
			fmt.Fprintf(query, `"%v"`, T.Field(i).Name)

			if i < T.NumField()-1 {
				query.WriteByte(',')
			}
		}
		fmt.Fprintf(query, ` FROM "%v"`+"\n", table.name)
		return query
	}

	fmt.Fprintf(query, `%v "%v"`, header, columns[0].Name())
	for _, column := range columns[1:] {
		fmt.Fprintf(query, `, "%v"`, column.Name())
	}
	fmt.Fprintf(query, ` FROM "%v"`+"\n", table.name)

	return query
}

//Select selects columns from a table, leave blank for all.
func (table *NewTable) Select(columns ...Column) Query {
	return table.selectType("SELECT", columns...)
}

//SelectDistinct selects unique values.
func (table *NewTable) SelectDistinct(columns ...Column) Query {
	return table.selectType("SELECT DISTINCT", columns...)
}

//Drop drops a table and permanantly deletes it's data.
func (table *NewTable) Drop() Query {
	var query = table.db.NewQuery()
	fmt.Fprintf(query, `DROP TABLE "%v"`+"\n", table.name)
	return query
}

//Truncate permanantly deletes the data in the table.
func (table *NewTable) Truncate() Query {
	var query = table.db.NewQuery()
	fmt.Fprintf(query, `Truncate TABLE "%v"`+"\n", table.name)
	return query
}

//Delete deletes the matching records.
func (table *NewTable) Delete() Query {
	var query = table.db.NewQuery()
	fmt.Fprintf(query, `DELETE FROM "%v"`, table.name)
	return query
}*/
