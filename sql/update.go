package sql

import (
	"errors"
	"fmt"
)

type Update struct {
	Column
	Value interface{}
}

var NoDatabase = errors.New("no database is attached to query")

func (q *Query) Update(updates ...Update) error {
	if len(updates) == 0 {
		return nil
	}

	var table = updates[0].Table

	var head Query
	fmt.Fprintf(&head, `UPDATE "%v" SET `, table)
	for i, update := range updates {
		fmt.Fprintf(&head, "\n\t"+`"%v" = %v`, update.Column.Name, q.value(update.Value))
		if i < len(updates)-1 {
			head.WriteByte(',')
		}
	}
	head.WriteString("\n\t")
	head.WriteQuery(q)

	if q.db.DB == nil {
		return NoDatabase
	}

	_, err := q.db.ExecContext(q.db.Context, head.String(), head.Values...)
	return err
}
