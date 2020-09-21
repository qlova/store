package db

func deleteTable(table Viewer) error {
	if !table.Master() {
		return ErrMasterProtected
	}
	if table.Database() == nil {
		return ErrDisconnectedViewer
	}
	return table.Database().Delete(table)
}

//Delete removes the viewer's table from the database.
//The master viewer must be passed to this function.
func Delete(table Viewer, tables ...Viewer) error {
	if err := deleteTable(table); err != nil {
		return err
	}
	for _, row := range tables {
		if err := deleteTable(row); err != nil {
			return err
		}
	}
	return nil
}

func empty(table Viewer) error {
	if !table.Master() {
		return ErrMasterProtected
	}
	if table.Database() == nil {
		return ErrDisconnectedViewer
	}
	return table.Database().Empty(table)
}

//Empty removes all rows from the given tables so that they are empty.
func Empty(table Viewer, tables ...Viewer) error {
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
