package db

func syncTable(table Table) error {
	if table.Database() == nil {
		return ErrDisconnectedViewer
	}
	return table.Database().Sync(table)
}

//Sync syncs the Tables with the Database, adding any missing columns.
//If constraints or types do not match up, an error is returned.
func Sync(table Table, tables ...Table) error {
	if err := syncTable(table); err != nil {
		return err
	}
	for _, table := range tables {
		if err := syncTable(table); err != nil {
			return err
		}
	}
	return nil
}
