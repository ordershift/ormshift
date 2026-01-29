package migrations

import "github.com/ordershift/ormshift"

type Migration interface {
	Up(migrator *Migrator) error
	Down(migrator *Migrator) error
}

func Migrate(database *ormshift.Database, config *MigratorConfig, migrations ...Migration) (*Migrator, error) {
	migrator, err := NewMigrator(database, config)
	if err != nil {
		return nil, err
	}
	for _, migration := range migrations {
		migrator.Add(migration)
	}
	err = migrator.ApplyAllMigrations()
	if err != nil {
		return nil, err
	}
	return migrator, nil
}
