package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"qlova.store/db"
)

var _ = pq.Driver{}

//Open sets the given connection to be backed by a postgres connection with the given options.
func Open(connection string) db.Driver {
	d, err := sql.Open("postgres", connection)

	return driver{d, err}
}

//Error wraps an error and a query.
type Error struct {
	error
	Query string
}

//Error returns an error string.
func (err Error) Error() string {
	return fmt.Sprintf("%v: %v", err.error.Error(), err.Query)
}

//Unwrap returns the internal error of the Error.
func (err Error) Unwrap() error {
	return err.error
}

func typeInfo(rtype reflect.Type) (tname string, tvalue string, err error) {
	var zero = reflect.Zero(rtype).Interface()

	switch zero.(type) {
	case int32:
		return "integer", `0`, nil
	case int64:
		return "bigint", `0`, nil

	case float32:
		return "real", `0`, nil
	case float64:
		return "double precision", `0`, nil

	case string:
		return "text", `''`, nil
	case []byte:
		return "bytea", `NULL`, nil

	case bool:
		return "boolean", `false`, nil

	case time.Time:
		return "timestamp", `'0001-01-01 00:00:00'`, nil
	case uuid.UUID:
		return "uuid", `'00000000-0000-0000-0000-000000000000'`, nil

	default:
		return "", "", errors.New("unsupported postgres db data type: " + rtype.String())
	}
}

func cname(name string) string {
	name = strings.ToLower(name)
	switch name {
	case
		"all", "analyse", "analyze", "and", "any", "array", "as", "asc", "asymetric",
		"authorization", "between", "binary", "both", "case", "cast", "check",
		"collate", "column", "constraint", "create", "cross", "current_date", "current_role",
		"current_time", "current_timestamp", "current_user", "default", "deferrable", "desc",
		"distinct", "do", "else", "end", "except", "false", "for", "foreign", "freeze", "from",
		"full", "grant", "group", "having", "ilike", "in", "initially", "inner", "intersect",
		"into", "is", "isnull", "join", "leading", "left", "like", "limit", "localtime", "localtimestamp",
		"natural", "new", "not", "notnull", "null", "off", "offset", "old", "on", "only", "or", "order",
		"outer", "overlaps", "placing", "primary", "references", "right", "select", "session_user",
		"similar", "some", "symmetric", "table", "then", "to", "trailing", "true", "union", "unique",
		"user", "using", "verbose", "when", "where":

		return strconv.Quote(name)
	default:
		return name
	}
}
