package sql

import (
	"github.com/qlova/store/db"
)

//Slicer implements db.Slicer
type slicer struct {
	length int

	columns []db.Column

	Query
}

func (slice slicer) Columns(col db.Column, cols ...db.Column) db.Slicer {
	slice.columns = append(cols, col)
	return slice
}

//Into implements db.Slicer.Into
func (slice slicer) Into(model db.Connectable) error {
	var q = slice.Query

	var columns = slice.columns

	if columns == nil {
		columns = db.Columns(model)
	}

	var head Query
	head.WriteString(`SELECT `)
	head.WriteColumn(columns[0])
	for _, col := range columns[1:] {
		head.WriteByte(',')
		head.WriteColumn(col)
	}
	head.WriteString(` FROM "`)
	head.WriteString(q.Table.Name)
	head.WriteString(`" `)

	head.WriteByte(' ')
	head.WriteQuery(&q)

	rows, err := q.Driver.DB.QueryContext(q.Driver.Context, head.String(), head.Values...)
	if err != nil {
		return head.Error(err)
	}

	var slices = db.Slices(model, slice.length, columns...)
	var pointers = make([]interface{}, len(slices))

	var index int
	for index = 0; rows.Next(); index++ {
		if err := rows.Err(); err != nil {
			return head.Error(err)
		}

		for i := range pointers {
			pointers[i] = slices[i].Index(index).Addr().Interface()
		}

		if err := rows.Scan(pointers...); err != nil {
			return head.Error(err)
		}
	}

	model.Slice(0, index)

	return nil
}
