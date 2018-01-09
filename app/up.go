package app

func (a Apricot) RunUp() error {
	// create schema table if needed
	err := a.DatabaseManager.Connect()
	if err != nil {
		return err
	}

	defer a.DatabaseManager.Close()

	if a.DatabaseManager.SchemaTableMissing() {
		if err := a.DatabaseManager.CreateSchemaTable(); err != nil {
			return err
		}
	}

	// find migrations
	// - `up` or `always`

	// sort migrations by version

	// filter migrations using schema table
	// - if version is newer or if the migration is `always`

	// apply migrations
	// - add to schema table - pending/unknown
	// - exec migration (think about ddl auto commit)
	// - update table row on success OR delete table row on failure

	return nil
}
