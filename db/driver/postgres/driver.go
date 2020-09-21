package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"qlova.store/db"
)

//Driver is a db driver for Postgres databases.
type driver struct {
	*sql.DB
	error
}

//Connect connects the given viewer to view this database.
//It then returns the database.
func (d driver) Connect(first db.Viewer, more ...db.Viewer) db.Driver {

	connect := func(v db.Viewer) {
		db.Connect(v, d)
	}

	connect(first)
	for _, viewer := range more {
		connect(viewer)
	}
	return d
}

//Insert inserts the given row into the database.
func (d driver) Insert(row db.Row, rows ...db.Row) error {

	insert := func(row db.Row) error {

		var insert db.Insertion
		if err := insert.Row(row); err != nil {
			return err
		}

		var query strings.Builder
		query.WriteString(`INSERT INTO `)
		query.WriteString(row.Row().Table())
		query.WriteString(` (`)

		for i, column := range insert.Columns {
			query.WriteString(cname(column))

			if i < len(insert.Columns)-1 {
				query.WriteByte(',')
			}
		}

		query.WriteString(`) VALUES (`)

		for i := range insert.Values {
			query.WriteByte('$')
			query.WriteString(strconv.Itoa(i + 1))

			if i < len(insert.Columns)-1 {
				query.WriteByte(',')
			}
		}

		query.WriteString(`);`)

		_, err := d.Exec(query.String(), insert.Values...)

		if err != nil {
			return Error{err, query.String()}
		}

		return err
	}

	if err := insert(row); err != nil {
		return err
	}
	for _, row := range rows {
		if err := insert(row); err != nil {
			return err
		}
	}
	return nil
}

//Delete deletes the given tables.
func (d driver) Delete(table db.Table, tables ...db.Table) error {

	delete := func(table db.Table) error {
		_, err := d.Exec(`DROP TABLE ` + table.Table() + `;`)

		return err
	}

	if err := delete(table); err != nil {
		return err
	}
	for _, table := range tables {
		if err := delete(table); err != nil {
			return err
		}
	}
	return nil
}

//Empty removes all rows from the given tables so that they are empty.
func (d driver) Empty(table db.Table, tables ...db.Table) error {

	empty := func(table db.Table) error {
		_, err := d.Exec(`TRUNCATE TABLE ` + table.Table() + `;`)
		return err
	}

	if err := empty(table); err != nil {
		return err
	}
	for _, table := range tables {
		if err := empty(table); err != nil {
			return err
		}
	}
	return nil
}

//Search returns results for the given filter.
func (d driver) Search(filter db.Filter) db.Results {
	var query strings.Builder
	var values []interface{}

	var joined = filter.Link.To != nil || len(filter.Links) > 0

	query.WriteString("FROM " + filter.Table)

	addLink := func(link db.Linker) {
		if link.To != nil {
			query.WriteString(` INNER JOIN `)
			query.WriteString(link.To.Table())
			query.WriteString(` ON `)
			query.WriteString(link.From.Table())
			query.WriteString(`.`)
			query.WriteString(link.From.Column())
			query.WriteString(`=`)
			query.WriteString(link.To.Table())
			query.WriteString(`.`)
			query.WriteString(link.To.Column())
			query.WriteString(` `)
		}
	}

	if joined {
		addLink(filter.Link)
		for _, link := range filter.Links {
			addLink(link)
		}
	}

	var addCondition func(c db.Condition)

	addCondition = func(c db.Condition) {
		if len(c.Cases) > 0 {
			query.WriteByte('(')
			defer func() {
				for _, c := range c.Cases {
					query.WriteString(" OR ")
					addCondition(c)
				}
				query.WriteByte(')')
			}()
		}

		if c.Operator == 0 {
			query.WriteString("TRUE")
			return
		}

		if joined {
			query.WriteString(cname(c.Table))
			query.WriteByte('.')
		}
		query.WriteString(cname(c.Column))
		if c.Operator != db.NoOperator {
			switch c.Operator {
			case db.OpEquals:
				query.WriteByte('=')
			case db.OpNotEquals:
				query.WriteByte('!')
				query.WriteByte('=')
			case db.OpLessThan:
				query.WriteByte('<')
			case db.OpContains:
				query.WriteString(" LIKE ")
				c.Value = fmt.Sprintf("%%%v%%", c.Value)
			case db.OpHasPrefix:
				query.WriteString(" LIKE ")
				c.Value = fmt.Sprintf("%v%%", c.Value)
			default:
				panic("unsupported operator: " + strconv.Itoa(int(c.Operator)))
			}
		}

		query.WriteByte('$')
		query.WriteString(strconv.Itoa(len(values) + 1))
		values = append(values, c.Value)
	}

	if filter.Condition.Operator != db.NoOperator || len(filter.Conditions) > 0 {
		query.WriteString(" WHERE ")
		addCondition(filter.Condition)
		for _, condition := range filter.Conditions {
			query.WriteString(" AND ")
			addCondition(condition)
		}
	}

	return results{
		pq:     d,
		query:  query.String(),
		values: values,

		joined: joined,

		table: filter.Table,

		view: filter.View,

		length: filter.Length,
		offset: filter.Offset,

		columns: filter.Columns,
	}
}

//Close closes the connection to the database.
func (d driver) Close() error {
	return d.DB.Close()
}
