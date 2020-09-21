package postgres

import (
	"errors"
	"fmt"
	"strings"

	"qlova.store/db"
)

//Sync syncs the Tables with the Database, adding any missing columns.
//If constraints or types do not match up, an error is returned.
func (d driver) Sync(table db.Table, tables ...db.Table) error {
	if d.error != nil {
		return d.error
	}

	sync := func(table db.Table) error {
		var query strings.Builder

		fmt.Fprintf(&query, `CREATE TABLE IF NOT EXISTS %v (`, table.Table())
		for i := 0; i < table.Columns(); i++ {
			column := table.Column(i)

			tname, dvalue, err := typeInfo(column.Type())
			if err != nil {
				return err
			}

			if column.Key() {
				tname += " PRIMARY KEY"
			}

			fmt.Fprintf(&query, `%v %v DEFAULT %v`, cname(column.Column()), tname, dvalue)

			if i < table.Columns()-1 {
				query.WriteByte(',')
			}
		}
		query.WriteByte(')')

		_, err := d.Exec(query.String())
		if err != nil {
			return Error{err, query.String()}
		}

		query = strings.Builder{}
		fmt.Fprintf(&query, `select * from %v LIMIT 1`, table.Table())

		rows, err := d.Query(query.String())
		if err != nil {
			return Error{err, query.String()}
		}

		columns, err := rows.Columns()
		if err != nil {
			return Error{err, query.String()}
		}

		rows.Close()

		ExistingColumns := make(map[string]struct{}, len(columns))
		for _, column := range columns {
			ExistingColumns[strings.ToLower(column)] = struct{}{}
		}

		//Query existing constraints.
		query = strings.Builder{}
		fmt.Fprintf(&query, `select
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
	`, table.Table())

		rows, err = d.Query(query.String())
		if err != nil {
			return Error{err, query.String()}
		}

		ExistingConstraints := make(map[string]struct{})
		for rows.Next() {
			if err := rows.Err(); err != nil {
				return Error{err, query.String()}
			}
			var name, ctype, check *string
			err := rows.Scan(&name, &ctype, &check)
			if err != nil {
				return Error{err, query.String()}
			}

			if name != nil {
				ExistingConstraints[*name+" "+*ctype] = struct{}{}
			} else if check != nil {
				ExistingConstraints[*check] = struct{}{}
			} else {
				return errors.New("unknown constraint: " + *ctype)
			}
		}

		//Ensure columns are in sync.
		for i := 0; i < table.Columns(); i++ {
			target := table.Column(i)

			if _, ok := ExistingColumns[strings.ToLower(target.Column())]; ok {
				continue
			}

			tname, dvalue, _ := typeInfo(target.Type())

			if target.Key() {
				tname += " PRIMARY KEY"
			}

			query = strings.Builder{}

			fmt.Fprintf(&query, `ALTER TABLE %v ADD %v %v DEFAULT %v`,
				table.Table(), target.Column(), tname, dvalue)

			_, err = d.Exec(query.String())
			if err != nil {
				return Error{err, query.String()}
			}
		}

		return nil
	}

	if err := sync(table); err != nil {
		return err
	}
	for _, table := range tables {
		if err := sync(table); err != nil {
			return err
		}
	}
	return nil
}
