package sql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/qlova/store/db"
)

//TypeOf returns the SQL type of the given db.Value.
func TypeOf(val db.Value) (string, error) {
	switch val.(type) {
	case db.Int:
		return "int", nil
	case db.String:
		return "text", nil
	case db.UUID:
		return "uuid", nil
	case db.Bool:
		return "boolean", nil
	case db.Time:
		return "timestamp", nil
	case db.Float64:
		return "double precision", nil
	default:
		return "", errors.New("sql: unsupported type: " + reflect.TypeOf(val).String())
	}
}

//DefaultValueOf returns the SQL default value of the given db.Value.
func DefaultValueOf(val db.Value) (string, error) {
	switch val.(type) {
	case db.Int:
		return "0", nil
	case db.String:
		return `''`, nil
	case db.UUID:
		return "'00000000-00000000-00000000-00000000'", nil
	case db.Bool:
		return "false", nil
	case db.Time:
		return "0001-01-01 00:00:00", nil
	case db.Float64:
		return "0", nil
	default:
		return "", errors.New("sql: unsupported type: " + reflect.TypeOf(val).String())
	}
}

//Verify implements db.Driver.Verify
func (d Driver) Verify(schema db.Schema) error {

	var q Query
	q.Driver = d

	q.WriteString(`CREATE TABLE IF NOT EXISTS "`)
	q.WriteString(schema.Table.Name)
	q.WriteString(`" (`)
	for i, val := range schema.Columns {
		col := val.GetColumn()

		t, err := TypeOf(val)
		if err != nil {
			return err
		}

		q.WriteString(`"` + col.Name + `" `)
		q.WriteString(t)
		q.WriteString(" NOT NULL")

		if i < len(schema.Columns)-1 {
			q.WriteByte(',')
		}
	}
	q.WriteByte(')')

	_, err := q.Driver.DB.ExecContext(q.Driver.Context, q.String())
	if err != nil {
		return q.Error(err)
	}

	q.Builder = strings.Builder{}
	fmt.Fprintf(&q, `select * from "%v" LIMIT 1`, schema.Table.Name)

	rows, err := q.Driver.DB.QueryContext(q.Driver.Context, q.String())
	if err != nil {
		return q.Error(err)
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	rows.Close()

	ExistingColumns := make(map[string]struct{}, len(columns))
	for _, column := range columns {
		ExistingColumns[column] = struct{}{}
	}

	//Query existing constraints.
	q.Builder = strings.Builder{}
	fmt.Fprintf(&q, `select
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
	`, schema.Table.Name)

	rows, err = q.Driver.DB.QueryContext(q.Driver.Context, q.String())
	if err != nil {
		return q.Error(err)
	}

	ExistingConstraints := make(map[string]struct{})
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		var name, ctype, check *string
		err := rows.Scan(&name, &ctype, &check)
		if err != nil {
			return err
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
	for _, value := range schema.Columns {
		var target = value.GetColumn()

		if _, ok := ExistingColumns[target.Name]; ok {

			/*for _, constraint := range target.Constraints {
				constraint = strings.ToUpper(strings.TrimSpace(constraint))

				if constraint == "" {
					continue
				}

				if constraint == "NOT NULL" {
					constraint = "IS NOT NULL"
				}

				if _, ok := ExistingConstraints[target.Name+" "+constraint]; !ok {
					return fmt.Errorf("constraint mismatch on column %v, missing constraint %v", target.Name, constraint)
				}
			}*/

			continue
		}

		t, _ := TypeOf(value)
		defaultVal, _ := DefaultValueOf(value)

		q.Builder = strings.Builder{}

		fmt.Fprintf(&q, `ALTER TABLE "%v"`+"\n"+`ADD "%v" %v %v`+"\n"+"DEFAULT %v", schema.Table.Name, target.Name, t, "NOT NULL", defaultVal)

		_, err = q.Driver.DB.ExecContext(q.Driver.Context, q.String())
		if err != nil {
			return q.Error(err)
		}
	}

	return nil

}
