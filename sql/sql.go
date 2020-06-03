package sql

import (
	"context"
	"database/sql"
)

const Debug = false

//Open returns a new database connection.
func Open(driver string, args string) (Database, error) {
	db, err := sql.Open(driver, args)

	return Database{db, context.Background()}, err
}

//Database is a sql database name.
type Database struct {
	*sql.DB
	context.Context
}

//CreateDatabase creates and returns a database.
/*func (conn Connection) CreateDatabase(name string) (Database, error) {
	_, err := conn.Exec("CREATE DATABASE " + name + ";")
	return Database{conn, name}, err
}

//EnsureDatabase is shorthand for CreateDatabaseIfNotExists.
func (conn Connection) EnsureDatabase(name string) (Database, error) {
	return conn.CreateDatabaseIfNotExists(name)
}

//CreateDatabaseIfNotExists creates and returns a database if it doesn't exist.
func (conn Connection) CreateDatabaseIfNotExists(name string) (Database, error) {
	_, err := conn.Exec("CREATE DATABASE " + name + ";")
	if err != nil && strings.Contains(err.Error(), "already exists") {
		err = nil
	}
	return Database{conn, name}, err
}*/
