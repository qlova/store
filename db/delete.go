package db

//Delete removes the viewer's table from the database.
//The master viewer must be passed to this function.
func Delete(table Viewer, tables ...Viewer) error {

	delete := func(table Viewer) error {
		if !table.Master() {
			return ErrMasterProtected
		}
		if table.Database() == nil {
			return ErrDisconnectedViewer
		}
		return table.Database().Delete(table)
	}

	if err := delete(table); err != nil {
		return err
	}
	for _, row := range tables {
		if err := delete(row); err != nil {
			return err
		}
	}
	return nil
}

//Empty removes all rows from the given tables so that they are empty.
func Empty(table Viewer, tables ...Viewer) error {

	empty := func(table Viewer) error {
		if !table.Master() {
			return ErrMasterProtected
		}
		if table.Database() == nil {
			return ErrDisconnectedViewer
		}
		return table.Database().Empty(table)
	}

	if err := empty(table); err != nil {
		return err
	}
	for _, row := range tables {
		if err := empty(row); err != nil {
			return err
		}
	}
	return nil
}
