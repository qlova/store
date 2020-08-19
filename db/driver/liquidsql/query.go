package liquidsql

import (
	"context"
	"database/sql/driver"
	"reflect"
	"strconv"
	"strings"

	"qlova.store/db"

	sqle "github.com/liquidata-inc/go-mysql-server"
	"github.com/liquidata-inc/go-mysql-server/sql"
)

var pid uint64

func query(eng *sqle.Engine, ctx *sql.Context, queryTemplate string, args ...interface{}) (sql.Schema, sql.RowIter, error) {
	var converted []driver.Value

	for _, val := range args {
		c, err := converter{}.ConvertValue(val)
		if err != nil {
			return nil, nil, err
		}
		converted = append(converted, c)
	}

	q, err := prepare(queryTemplate, converted)
	if err != nil {
		return nil, nil, err
	}

	pid++

	return eng.Query(sql.NewContext(context.Background(),
		sql.WithPid(pid),
	), q)
}

func queryRow(eng *sqle.Engine, ctx *sql.Context, queryTemplate string, args ...interface{}) sql.Row {
	_, iter, err := query(eng, ctx, queryTemplate, args...)
	if err != nil {
		return nil
	}
	row, err := iter.Next()
	if err != nil {
		return nil
	}
	iter.Close()
	return row
}

//Query implements db.Query
type Query struct {
	Driver Driver
	Table  db.Table

	strings.Builder
	Values []interface{}

	sortby []db.Column
}

type Error struct {
	Internal error

	Query  string
	Values []interface{}
}

const Debug = true

func (err Error) Error() string {
	if Debug {
		return err.Internal.Error() + " " + err.Query
	}
	return err.Internal.Error()
}

func (q Query) Error(internal error) error {
	if internal == nil {
		return nil
	}

	return Error{
		Internal: internal,

		Query:  q.String(),
		Values: q.Values,
	}
}

//WriteQuery writes a Query to the Query.
func (q *Query) WriteQuery(other *Query) {
	q.WriteString(other.Builder.String())
	q.Values = append(q.Values, other.Values...)
}

//WriteColumn writes a Column to the Query.
func (q *Query) WriteColumn(col db.Column) {
	col.Name = filterReservedWord(col.Name)

	q.WriteString(col.Table.Name)
	q.WriteByte('.')
	q.WriteString(col.Name)
}

//WriteUpdate writes a db.Update to the Query.
func (q *Query) WriteUpdate(update db.Update) {
	update.Column.Name = filterReservedWord(update.Column.Name)

	q.WriteString(update.Column.Name)
	q.WriteByte('=')
	q.WriteString("?")

	q.Values = append(q.Values, update.Value)
}

//WriteCondition writes a db.Condition to the Query.
func (q *Query) WriteCondition(condition db.Condition) {
	if condition.Or != nil {
		q.WriteByte('(')
	}

	q.WriteColumn(condition.Column)

	switch condition.Operator {
	case db.Equals:
		q.WriteByte('=')
	case db.NotEquals:
		q.WriteString("!=")
	case db.Contains:
		q.WriteString(" LIKE ")
	case db.StartsWith:
		q.WriteString(" LIKE ")
	default:
		q.WriteString("[sql: unsupported operator]")
	}

	if condition.Operator == db.Contains {
		condition.Value = "%" + condition.Value.(string) + "%"
	}
	if condition.Operator == db.StartsWith {
		condition.Value = condition.Value.(string) + "%"
	}

	q.WriteString("?")

	q.Values = append(q.Values, condition.Value)

	if condition.Or != nil {
		q.WriteString(` OR `)
		q.WriteCondition(*condition.Or)
		q.WriteByte(')')
	}
}

func (q *Query) String() string {
	return q.Builder.String()
}

//Delete implements db.Query.Delete
func (q *Query) Delete() (int, error) {
	var head Query
	head.WriteString(`DELETE FROM `)
	head.WriteString(q.Table.Name)
	head.WriteString(` `)

	head.WriteQuery(q)

	_, _, err := query(q.Driver.Engine, q.Driver.Context, head.String(), head.Values...)
	if err != nil {
		return 0, head.Error(err)
	}

	return -1, nil
}

func (q *Query) Where(condition db.Condition, conditions ...db.Condition) db.Query {
	q.WriteString("WHERE ")
	q.WriteCondition(condition)

	for _, condition := range conditions {
		q.WriteString(" AND ")
		q.WriteCondition(condition)
	}

	q.WriteByte(' ')

	return q
}

func (q *Query) Link(db.Linker, ...db.Linker) db.Query {
	panic("not implemented")
}

func (q *Query) Slice(index, length int, values ...db.Value) db.Slicer {
	q.WriteString(" LIMIT ")
	q.WriteString(strconv.Itoa(length))
	q.WriteString(" OFFSET ")
	q.WriteString(strconv.Itoa(index))

	return slicer{length, nil, values, *q}
}

//Get implements db.Query.Get
func (q *Query) Get(value db.MutableValue, more ...db.MutableValue) error {
	var head Query
	head.WriteString(`SELECT `)
	head.WriteColumn(value.GetColumn())
	for _, val := range more {
		head.WriteByte(',')
		head.WriteColumn(val.GetColumn())
	}
	head.WriteString(` FROM `)
	head.WriteString(q.Table.Name)

	head.WriteByte(' ')
	head.WriteQuery(q)

	head.WriteString(`LIMIT 1`)

	row := queryRow(q.Driver.Engine, q.Driver.Context, head.String(), head.Values...)

	var pointers = make([]interface{}, len(more)+1)
	pointers[0] = value.Pointer()
	for i, val := range more {
		pointers[i+1] = val.Pointer()
	}

	return head.Error(scan(row, pointers...))
}

//Count implements db.Query.Count
func (q *Query) Count(v db.Value) (int, error) {
	var head Query
	head.WriteString(`SELECT `)
	head.WriteString(`COUNT(`)
	head.WriteColumn(v.GetColumn())
	head.WriteByte(')')

	head.WriteString(` FROM `)
	head.WriteString(q.Table.Name)
	head.WriteString(` `)

	head.WriteByte(' ')
	head.WriteQuery(q)

	row := queryRow(q.Driver.Engine, q.Driver.Context, head.String(), head.Values...)

	var count int

	err := scan(row, &count)

	return count, q.Error(err)
}

//Average implements db.Query.Average
func (q *Query) Average(v db.Value) (float64, error) {
	var head Query
	head.WriteString(`SELECT `)
	head.WriteString(`AVG(`)
	head.WriteColumn(v.GetColumn())
	head.WriteByte(')')

	head.WriteString(` FROM `)
	head.WriteString(q.Table.Name)
	head.WriteString(` `)

	head.WriteByte(' ')
	head.WriteQuery(q)

	row := queryRow(q.Driver.Engine, q.Driver.Context, head.String(), head.Values...)

	var count float64

	err := scan(row, &count)
	return count, q.Error(err)
}

//Read implements db.Query.Read
func (q *Query) Read(value db.Connectable) error {
	var head Query
	head.WriteString(`SELECT `)

	var columns = db.Columns(value)
	head.WriteColumn(columns[0])
	for _, col := range columns[1:] {
		head.WriteByte(',')
		head.WriteColumn(col)
	}

	head.WriteString(` FROM `)
	head.WriteString(q.Table.Name)
	head.WriteString(` `)

	head.WriteByte(' ')
	head.WriteQuery(q)

	head.WriteString(`LIMIT 1`)

	row := queryRow(q.Driver.Engine, q.Driver.Context, head.String(), head.Values...)

	var reflected = reflect.ValueOf(value).Elem()

	var pointers = make([]interface{}, len(columns))
	for i, col := range columns {
		pointers[i] = reflected.Field(int(col.Field)).Addr().Interface().(db.MutableValue).Pointer()
	}

	return q.Error(scan(row, pointers...))
}

//Update implements db.Driver.Update
func (q *Query) Update(update db.Update, updates ...db.Update) (int, error) {
	var head Query
	head.WriteString(`UPDATE `)
	head.WriteString(q.Table.Name)
	head.WriteString(` `)

	head.WriteString("SET ")
	head.WriteUpdate(update)

	for _, update := range updates {
		head.WriteString(",")
		head.WriteUpdate(update)
	}

	head.WriteByte(' ')

	head.WriteQuery(q)

	_, _, err := query(q.Driver.Engine, q.Driver.Context, head.String(), head.Values...)
	if err != nil {
		return 0, head.Error(err)
	}

	return -1, nil
}

//SortBy implements db.Driver.SortBy
func (d Driver) SortBy(column db.Column, columns ...db.Column) db.Query {
	var q Query
	q.Driver = d
	q.Table = column.Table

	q.SortBy(column, columns...)

	return &q
}

func (q *Query) SortBy(column db.Column, columns ...db.Column) db.Query {
	q.WriteString(" ORDER BY ")
	q.WriteColumn(column)
	if column.SortMode == db.Decreasing {
		q.WriteString(" DESC")
	}

	for _, val := range columns {
		q.WriteByte(',')
		q.WriteColumn(val.GetColumn())
		if val.GetColumn().SortMode == db.Decreasing {
			q.WriteString(" DESC")
		}
	}

	return q
}

//Update implements db.Driver.Update
func (d Driver) Update(db.Update, ...db.Update) (int, error) {
	panic("unimplemented")
}

//Where implements db.Driver.Where
func (d Driver) Where(condition db.Condition, conditions ...db.Condition) db.Query {
	var q Query
	q.Driver = d
	q.Table = condition.Column.Table

	return q.Where(condition, conditions...)
}
