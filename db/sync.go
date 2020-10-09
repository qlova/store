package db

//Sync syncs the Tables with the Database, adding any missing columns.
//If constraints or types do not match up, an error is returned.
func Sync(table Table, tables ...Table) error {

	sync := func(table Table) error {
		if table.Database() == nil {
			return ErrDisconnectedViewer
		}
		return table.Database().Sync(table)
	}

	if err := sync(table); err != nil {
		return err
	}
	for _, table := range tables {
		if err := sync(table); err != nil {
			return err
		}
	}
	return nil
}
