package migrations

import "github.com/ordershift/ormshift"

type Migration interface {
	Up(pMigrator *Migrator) error
	Down(pMigrator *Migrator) error
}

func Migrate(pDatabase *ormshift.Database, pConfig *MigratorConfig, pMigrations ...Migration) (*Migrator, error) {
	lMigrator, lError := NewMigrator(pDatabase, pConfig)
	if lError != nil {
		return nil, lError
	}
	for _, lMigration := range pMigrations {
		lMigrator.Add(lMigration)
	}
	lError = lMigrator.ApplyAllMigrations()
	if lError != nil {
		return nil, lError
	}
	return lMigrator, nil
}
