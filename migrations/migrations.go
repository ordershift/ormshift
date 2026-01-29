package migrations

import "github.com/ordershift/ormshift"

type Migration interface {
	Up(migrator *Migrator) error
	Down(migrator *Migrator) error
}

func Migrate(database *ormshift.Database, config *MigratorConfig, migrations ...Migration) (*Migrator, error) {
	lMigrator, lError := NewMigrator(database, config)
	if lError != nil {
		return nil, lError
	}
	for _, lMigration := range migrations {
		lMigrator.Add(lMigration)
	}
	lError = lMigrator.ApplyAllMigrations()
	if lError != nil {
		return nil, lError
	}
	return lMigrator, nil
}
