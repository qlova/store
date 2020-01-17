package sql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//CreateTable creates the given table.
func (db Database) CreateTable(table Table) Query {
	return db.createTable(table, "CREATE TABLE")
}

//EnsureTable creates and ensures that the table meets the struct specification.
//Will create any missing columns, does not drop or rename any columns.
func (db Database) EnsureTable(table Table) error {
	if err := db.CreateTableIfNotExists(table).Do().Error(); err != nil {
		return err
	}

	//Blank selection to retrieve columns.
	result := table.Table().Select(star).Where(False).Do()
	if err := result.Error(); err != nil {
		return err
	}

	columns, err := result.Rows.Columns()
	if err != nil {
		return err
	}

	ExistingColumns := make(map[string]struct{}, len(columns))
	for _, column := range columns {
		ExistingColumns[column] = struct{}{}
	}

	//Query existing constraints.
	var query = db.NewQuery()
	fmt.Fprintf(query, `select 
		INFORMATION_SCHEMA.constraint_column_usage.column_name,
		INFORMATION_SCHEMA.TABLE_CONSTRAINTS.constraint_type, 
		INFORMATION_SCHEMA.CHECK_CONSTRAINTS.check_clause
		from INFORMATION_SCHEMA.TABLE_CONSTRAINTS
		full outer join INFORMATION_SCHEMA.CHECK_CONSTRAINTS on                                                                                                          
		INFORMATION_SCHEMA.TABLE_CONSTRAINTS.constraint_name=INFORMATION_SCHEMA.CHECK_CONSTRAINTS.constraint_name
		full outer join INFORMATION_SCHEMA.key_column_usage on                                                                                                                       
		INFORMATION_SCHEMA.TABLE_CONSTRAINTS.constraint_name=INFORMATION_SCHEMA.key_column_usage.constraint_name
		full outer join INFORMATION_SCHEMA.constraint_column_usage on 
		INFORMATION_SCHEMA.TABLE_CONSTRAINTS.constraint_name=INFORMATION_SCHEMA.constraint_column_usage.constraint_name
		where INFORMATION_SCHEMA.TABLE_CONSTRAINTS.table_name='%v';
	`, strings.ToLower(table.Table().name))
	result = query.Do()
	if err := result.Error(); err != nil {
		return err
	}

	ExistingConstraints := make(map[string]struct{})
	for result.Rows.Next() {
		var name, ctype, check *string
		err := result.Rows.Scan(&name, &ctype, &check)
		if err != nil {
			return err
		}

		if name != nil {
			ExistingConstraints[strings.ToLower(*name+" "+*ctype)] = struct{}{}
		} else if check != nil {
			ExistingConstraints[strings.ToLower(*check)] = struct{}{}
		} else {
			return errors.New("unknown constraint: " + *ctype)
		}
	}

	if result.Err() != nil {
		return result.Err()
	}

	fmt.Println(ExistingConstraints)

	TargetColumns := TableInfo{table}.Columns()

	//Ensure columns are in sync.
	for _, target := range TargetColumns {
		if _, ok := ExistingColumns[target.Name]; ok {

			for _, constraint := range target.Constraints {
				constraint = strings.TrimSpace(constraint)

				if constraint == "" {
					continue
				}

				if constraint == "not null" {
					constraint = "is not null"
				}

				if _, ok := ExistingConstraints[target.Name+" "+constraint]; !ok {
					return fmt.Errorf("constraint mismatch on column %v, missing constraint %v", target.Name, constraint)
				}
			}

			continue
		}

		var query = db.NewQuery()
		fmt.Fprintf(query, "ALTER TABLE %v\nADD %v %v %v", table.Table().name, target.Name, target.Datatype, strings.Join(target.Constraints, " "))
		result := query.Do()
		if err := result.Error(); err != nil {
			return err
		}
	}

	return nil
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
		Addr().Interface().(*NewTable)).set(db, name, table)

	var query = db.NewQuery()
	fmt.Fprintf(query, `%v %v (`, header, name)

	var info = TableInfo{table}
	var columns = info.Columns()

	for i, column := range columns {
		fmt.Fprintf(query, "\n\t%v %v %v", column.Name, column.Datatype, strings.Join(column.Constraints, " "))
		if i < len(columns)-1 {
			query.WriteByte(',')
		}
	}

	query.WriteByte('\n')
	query.WriteByte(')')

	return query
}
