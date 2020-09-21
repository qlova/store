package db

//Error is a database error.
type Error string

func (err Error) Error() string {
	return string(err)
}

//ErrUnregisteredViewer means that an unregistered viewer was used, register the viewer to resolve this error.
const ErrUnregisteredViewer Error = "unregistered viewer"

//ErrDisconnectedViewer means that an disconnected viewer was used, connect the viewer to a database to resolve this error.
const ErrDisconnectedViewer Error = "disconnected viewer"

//ErrIllegalMaster means that an attempt was made to write to a master viewer, assign the viewer to a variable to resolve this error.
const ErrIllegalMaster Error = "illegal use of master viewer"

//ErrMasterProtected means that an attempt was made to peform an operation on a clone of master.
//The operation was blocked for data-protection purposes. If you meant to perform this operation, pass a master viewer.
const ErrMasterProtected Error = "this operation must be performed on a master viewer"

//ErrTableNotFound means that the operation failed because the viewer's table was not found.
//Check that the name of the table is correct, or create it with Sync()
const ErrTableNotFound Error = "table not found"

//ErrDuplicateKey means that a value tried to be inserted into the database but it failed due to having a duplicate value.
const ErrDuplicateKey Error = "duplicate key"

//ErrNotFound is returned if no row is found when getting from the database.
const ErrNotFound Error = "row not found"
