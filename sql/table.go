package sql

import (
	"fmt"
	"reflect"
)

//Table is a sql Table.
type Table interface {
	Table() *NewTable
}

//NewTable is embedded in structs to create sql table definitions.
type NewTable struct {
	db   Database
	name string

	structure Table
}

//Table defines that New Table is table.
func (table *NewTable) Table() *NewTable {
	return table
}

func (table *NewTable) set(db Database, name string, structure Table) {
	table.db = db
	table.name = name
	table.structure = structure
}

//Insert inserts a struct into this table.
func (table *NewTable) Insert(structure interface{}) Query {
	var T = reflect.TypeOf(structure)
	if T.Kind() != reflect.Struct {
		return table.db.QueryError("Table.Insert: must be struct")
	}

	var value = reflect.ValueOf(structure)

	var query = table.db.NewQuery()
	var tail = table.db.NewQuery()

	fmt.Fprintf(query, `INSERT INTO %v (`, table.name)
	fmt.Fprint(tail, "\n)\nVALUES (")

	for i := 0; i < T.NumField(); i++ {
		fmt.Fprintf(query, "\n\t %v", T.Field(i).Name)
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

	fmt.Fprintf(query, `INSERT INTO %v (`, table.name)
	fmt.Fprint(tail, "\n)\nVALUES (")

	for i, value := range values {
		fmt.Fprintf(query, "\n\t %v", value.Column())
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

	fmt.Fprintf(query, `UPDATE %v SET `, table.name)

	for i := 0; i < T.NumField(); i++ {
		fmt.Fprintf(query, "\n\t %v = %v", T.Field(i).Name, query.value(value.Field(i).Interface()))
		if i < T.NumField()-1 {
			query.WriteByte(',')
		}
	}

	return query
}

//UpdateValues updates the values of the filtered records.
func (table *NewTable) UpdateValues(values Values) Query {
	var query = table.db.NewQuery()

	fmt.Fprintf(query, `UPDATE %v SET `, table.name)

	for i, value := range values {
		fmt.Fprintf(query, "\n\t %v = %v", value.Column(), value.get(query))
		if i < len(values)-1 {
			query.WriteByte(',')
		}
	}

	return query
}

func (table *NewTable) selectType(header string, columns ...Column) Query {
	var query = table.db.NewQuery()

	if len(columns) == 0 {
		var T = reflect.TypeOf(table.structure).Elem()
		fmt.Fprintf(query, "%v ", header)
		for i := 1; i < T.NumField(); i++ {
			fmt.Fprintf(query, "%v", T.Field(i).Name)
			if i < T.NumField()-1 {
				query.WriteByte(',')
			}
		}
		fmt.Fprintf(query, " FROM %v\n", table.name)
		return query
	}

	fmt.Fprintf(query, "%v %v", header, columns[0].Name())
	for _, column := range columns[1:] {
		fmt.Fprintf(query, ", %v", column.Name())
	}
	fmt.Fprintf(query, " FROM %v\n", table.name)

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
	fmt.Fprintf(query, "DROP TABLE %v\n", table.name)
	return query
}

//Truncate permanantly deletes the data in the table.
func (table *NewTable) Truncate() Query {
	var query = table.db.NewQuery()
	fmt.Fprintf(query, "Truncate TABLE %v\n", table.name)
	return query
}

//Delete deletes the matching records.
func (table *NewTable) Delete() Query {
	var query = table.db.NewQuery()
	fmt.Fprintf(query, `DELETE FROM %v`, table.name)
	return query
}
