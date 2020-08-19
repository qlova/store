package liquidsql

import (
	"context"

	sqle "github.com/liquidata-inc/go-mysql-server"
	"github.com/liquidata-inc/go-mysql-server/memory"
	"github.com/liquidata-inc/go-mysql-server/sql"

	"qlova.store/db"
)

//Driver implements db.Driver
type Driver struct {
	*sqle.Engine
	*sql.Context
}

//Open sets the named db connection to the given sql driver and options.
func open(name db.Connection, driver, options string) error {

	engine := sqle.NewDefault()
	engine.AddDatabase(memory.NewDatabase(""))
	engine.AddDatabase(sql.NewInformationSchemaDatabase(engine.Catalog))

	db.Open(name, Driver{engine, sql.NewContext(context.Background())})

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

	q.WriteString(`INSERT INTO `)
	q.WriteString(insert.Table.Name)
	q.WriteString(` (`)

	for i, column := range insert.Columns {
		column.Name = filterReservedWord(column.Name)

		q.WriteString(column.Name)

		if i < len(insert.Columns)-1 {
			q.WriteByte(',')
		}
	}

	q.WriteString(`) VALUES (`)

	for i := range insert.Values {
		q.WriteByte('?')

		if i < len(insert.Columns)-1 {
			q.WriteByte(',')
		}
	}

	q.WriteString(`);`)

	_, _, err := query(d.Engine, d.Context, q.String(), insert.Values...)

	return q.Error(err)
}

//Link implements db.Driver.Link
func (d Driver) Link(link db.Linker, links ...db.Linker) db.Query {
	var q Query
	q.Driver = d
	q.Table = link.From.Table

	q.WriteString(`INNER JOIN `)
	q.WriteString(link.To.Table.Name)
	q.WriteString(` ON `)
	q.WriteString(link.From.Table.Name)
	q.WriteString(`.`)
	q.WriteString(link.From.Name)
	q.WriteString(`=`)
	q.WriteString(link.To.Table.Name)
	q.WriteString(`.`)
	q.WriteString(link.To.Name)
	q.WriteString(` `)

	return &q
}

//Link implements db.Driver.Link
func (d Driver) Read(interface{}, ...interface{}) (int, error) {
	panic("unimplemented")
}

//Delete implements db.Driver.Delete
func (d Driver) Delete(table db.Table, tables ...db.Table) error {
	_, _, err := query(d.Engine, d.Context, `DROP TABLE `+table.Name+`;`)
	return err
}

//Truncate implements db.Driver.Truncate
func (d Driver) Truncate(table db.Table, tables ...db.Table) error {
	_, _, err := query(d.Engine, d.Context, `TRUNCATE TABLE `+table.Name+`;`)
	return err
}
