package postgres

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"qlova.store/db"
)

type results struct {
	pq     driver
	query  string
	values []interface{}

	joined bool

	table string

	view db.Table

	offset, length int

	columns []db.Variable
}

//MarshalJSON implements json.Marshaler
func (r results) MarshalJSON() ([]byte, error) {
	var query strings.Builder
	query.WriteString(`SELECT `)

	var ColumnCount int

	if r.columns != nil {

		for i, col := range r.columns {
			if i > 0 {
				query.WriteByte(',')
			}
			if r.joined {
				query.WriteString(col.Table())
				query.WriteByte('.')
			}
			query.WriteString(cname(col.Column()))
		}

		ColumnCount = len(r.columns)

	} else {

		for i := 0; i < r.view.Columns(); i++ {
			column := r.view.Column(i)

			if i > 0 {
				query.WriteByte(',')
			}
			if r.joined {
				query.WriteString(r.view.Table())
				query.WriteByte('.')
			}
			query.WriteString(cname(column.Column()))
		}

		ColumnCount = r.view.Columns()

	}

	query.WriteByte(' ')
	query.WriteString(r.query)

	//fmt.Printf("\n\n"+query.String(), r.values...)

	rows, err := r.pq.Query(query.String(), r.values...)
	if err != nil {
		return nil, Error{err, query.String()}
	}

	results := make([]interface{}, ColumnCount)

	pointers := make([]interface{}, ColumnCount)
	for i := 0; i < ColumnCount; i++ {
		pointers[i] = &results[i]
	}

	var buffer bytes.Buffer
	buffer.WriteString(`[`)

	var index int
	for index = 0; rows.Next(); index++ {
		if index != 0 {
			buffer.WriteByte(',')
		}

		if err := rows.Err(); err != nil {
			return nil, Error{err, query.String()}
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, Error{err, query.String()}
		}

		buffer.WriteString(`{`)
		for i, result := range results {

			if r.columns != nil {

				buffer.WriteString(strconv.Quote(r.columns[i].Column()))
				buffer.WriteByte(':')

				switch r.columns[i].(type) {
				case *db.UUID:
					var id uuid.UUID
					id.Scan(result)
					buffer.WriteString(strconv.Quote(id.String()))
				default:
					encoded, err := json.Marshal(result)
					if err != nil {
						return nil, Error{err, query.String()}
					}
					buffer.Write(encoded)
				}

			} else {

				buffer.WriteString(strconv.Quote(r.view.Column(i).Column()))
				buffer.WriteByte(':')

				switch r.view.Column(i).(type) {
				case *db.UUID:
					var id uuid.UUID
					id.Scan(result)
					buffer.WriteString(strconv.Quote(id.String()))
				default:
					encoded, err := json.Marshal(result)
					if err != nil {
						return nil, Error{err, query.String()}
					}
					buffer.Write(encoded)
				}

			}

			if i < len(results)-1 {
				buffer.WriteByte(',')
			}
		}
		buffer.WriteString(`}`)
	}

	buffer.WriteByte(']')

	return buffer.Bytes(), nil
}

//Update updates the results with the given updates.
//Returns the number of results updated (or -1 if the statistic is unavailable).
func (r results) Update(update db.Update, updates ...db.Update) (int, error) {

	var query strings.Builder
	query.WriteString(`UPDATE `)
	query.WriteString(r.table)
	query.WriteString(` `)

	query.WriteString("SET ")

	addupdate := func(update db.Update) {
		query.WriteString(update.Column)
		query.WriteByte('=')
		query.WriteString("$")
		query.WriteString(strconv.Itoa(len(r.values) + 1))

		r.values = append(r.values, update.Value)
	}

	addupdate(update)
	for _, update := range updates {
		query.WriteString(",")
		addupdate(update)
	}

	query.WriteByte(' ')

	query.WriteString(strings.TrimPrefix(r.query, "FROM "+r.table))

	result, err := r.pq.Exec(query.String(), r.values...)
	if err != nil {
		return 0, Error{err, query.String()}
	}

	number, err := result.RowsAffected()
	if err != nil {
		number = -1
	}

	return int(number), nil
}

//Delete deletes all the results from the database.
func (r results) Delete() (int, error) {
	var query strings.Builder
	query.WriteString(`DELETE `)

	query.WriteString(r.query)

	result, err := r.pq.Exec(query.String(), r.values...)
	if err != nil {
		return 0, Error{err, query.String()}
	}

	number, err := result.RowsAffected()
	if err != nil {
		number = -1
	}

	return int(number), nil
}

//Get gets the matching columns of the results.
func (r results) Get(variable db.Variable, variables ...db.Variable) (int, error) {
	var query strings.Builder
	query.WriteString(`SELECT `)

	if r.joined {
		query.WriteString(variable.Table())
		query.WriteByte('.')
	}
	query.WriteString(cname(variable.Column()))

	for _, variable := range variables {
		query.WriteByte(',')
		if r.joined {
			query.WriteString(variable.Table())
			query.WriteByte('.')
		}
		query.WriteString(cname(variable.Column()))
	}

	query.WriteByte(' ')
	query.WriteString(r.query)

	query.WriteString(` LIMIT `)
	query.WriteString(strconv.Itoa(r.length))

	//fmt.Printf("\n\n"+query.String()+"\n", r.values...)

	if r.length == 1 {
		row := r.pq.QueryRow(query.String(), r.values...)

		var pointers = make([]interface{}, len(variables)+1)

		pointers[0] = variable.Pointer()

		for i, variable := range variables {
			pointers[i+1] = variable.Pointer()
		}

		if err := row.Scan(pointers...); err != nil {

			if err == sql.ErrNoRows {
				return 0, db.ErrNotFound
			}
			return 0, Error{err, query.String()}
		}

		return 1, nil
	}

	rows, err := r.pq.Query(query.String(), r.values...)
	if err != nil {
		return 0, Error{err, query.String()}
	}

	var pointers = make([]interface{}, len(variables)+1)

	variable.Make(r.length)
	for _, variable := range variables {
		variable.Make(r.length)
	}

	var index int
	for index = 0; rows.Next(); index++ {
		if err := rows.Err(); err != nil {
			return 0, Error{err, query.String()}
		}

		pointers[0] = variable.Slice(index)
		for i, variable := range variables {
			pointers[i+1] = variable.Slice(index)
		}

		if err := rows.Scan(pointers...); err != nil {
			if err == sql.ErrNoRows {
				return 0, db.ErrNotFound
			}
			return 0, Error{err, query.String()}
		}

	}

	return index + 1, nil
}

//Count returns the number of results.
func (r results) Count(value db.Viewable) (int, error) {
	var query strings.Builder
	query.WriteString(`SELECT `)
	query.WriteString(`COUNT(*)`)
	query.WriteByte(' ')

	query.WriteString(r.query)

	row := r.pq.QueryRow(query.String(), r.values...)

	var count int

	err := row.Scan(&count)

	if err != nil {
		return count, Error{err, query.String()}
	}

	return count, err
}

//Sum returns the sum amount of the value in the given column of all results.
func (r results) Sum(value db.Variable) error {
	var query strings.Builder
	query.WriteString(`SELECT `)
	query.WriteString(`SUM(`)
	if r.joined {
		query.WriteString(value.Table())
		query.WriteByte('.')
	}
	query.WriteString(value.Column())
	query.WriteByte(')')
	query.WriteByte(' ')

	query.WriteString(r.query)

	row := r.pq.QueryRow(query.String(), r.values...)

	err := row.Scan(value.Pointer())

	if err != nil {
		if strings.Contains(err.Error(), "converting NULL to") {
			fmt.Sscan("0", value.Pointer())
			return nil
		}
		return Error{err, query.String()}
	}

	return nil
}

//Average returns the average value in the given column for all results.
func (r results) Average(value db.Viewable) (float64, error) {
	var query strings.Builder
	query.WriteString(`SELECT `)
	query.WriteString(`AVG(`)
	if r.joined {
		query.WriteString(value.Table())
		query.WriteByte('.')
	}
	query.WriteString(value.Column())
	query.WriteByte(')')
	query.WriteByte(' ')

	query.WriteString(r.query)

	row := r.pq.QueryRow(query.String(), r.values...)

	var avg *float64

	err := row.Scan(&avg)

	if avg == nil {
		return math.NaN(), Error{err, query.String()}
	}

	if err != nil {
		return 0, Error{err, query.String()}
	}

	return *avg, err
}
