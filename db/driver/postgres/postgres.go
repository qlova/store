package postgres

import (
	"github.com/lib/pq"
	"github.com/qlova/store/db"
	"github.com/qlova/store/db/driver/sql"
)

var _ = pq.Driver{}

//Open sets the given connection to be backed by a postgres connection with the given options.
func Open(connection db.Connection, options string) error {
	return sql.Open(connection, "postgres", options)
}
