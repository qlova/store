package sql

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/qlova/store/db"
)

//Driver implements db.Driver
type Driver struct {
	*sql.DB
	context.Context
}

//Open sets the named db connection to the given sql driver and options.
func Open(name db.Connection, driver, options string) error {

	d, err := sql.Open(driver, options)
	if err != nil {
		return nil
	}

	db.Open(name, Driver{d, context.Background()})

	return nil
}

func (d Driver) Slice(index, length int, values ...db.Value) db.Slicer {
	panic("not implemented")
}

func (d Driver) Count(v db.Value) (int, error) {
	panic("not implemented")
}

func (d Driver) Average(v db.Value) (float64, error) {
	panic("not implemented")
}

//Insert implements db.Driver.Insert
func (d Driver) Insert(row db.Row, rows ...db.Row) error {

	var insert db.Insertion
	if err := insert.Row(row); err != nil {
		return err
	}

	var q Query

	q.WriteString(`INSERT INTO "`)
	q.WriteString(insert.Table.Name)
	q.WriteString(`" (`)

	for i, column := range insert.Columns {
		q.WriteByte('"')
		q.WriteString(column.Name)
		q.WriteByte('"')

		if i < len(insert.Columns)-1 {
			q.WriteByte(',')
		}
	}

	q.WriteString(`) VALUES (`)

	for i := range insert.Values {
		q.WriteByte('$')
		q.WriteString(strconv.Itoa(i + 1))

		if i < len(insert.Columns)-1 {
			q.WriteByte(',')
		}
	}

	q.WriteString(`);`)

	_, err := d.DB.ExecContext(d.Context, q.String(), insert.Values...)

	return q.Error(err)
}

//Link implements db.Driver.Link
func (d Driver) Link(link db.Linker, links ...db.Linker) db.Query {
	var q Query
	q.Driver = d
	q.Table = link.From.Table

	q.WriteString(`INNER JOIN "`)
	q.WriteString(link.To.Table.Name)
	q.WriteString(`" ON "`)
	q.WriteString(link.From.Table.Name)
	q.WriteString(`"."`)
	q.WriteString(link.From.Name)
	q.WriteString(`"="`)
	q.WriteString(link.To.Table.Name)
	q.WriteString(`"."`)
	q.WriteString(link.To.Name)
	q.WriteString(`" `)

	return &q
}

//Link implements db.Driver.Link
func (d Driver) Read(interface{}, ...interface{}) (int, error) {
	panic("unimplemented")
}

//Delete implements db.Driver.Delete
func (d Driver) Delete(table db.Table, tables ...db.Table) error {
	_, err := d.DB.ExecContext(d.Context, `DROP TABLE "`+table.Name+`";`)
	return err
}

//Truncate implements db.Driver.Truncate
func (d Driver) Truncate(table db.Table, tables ...db.Table) error {
	_, err := d.DB.ExecContext(d.Context, `TRUNCATE TABLE "`+table.Name+`";`)
	return err
}
