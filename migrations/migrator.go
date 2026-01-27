package migrations

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

type Migrator struct {
	database          *ormshift.Database
	config            MigratorConfig
	migrations        []Migration
	appliedMigrations map[string]bool
}

func NewMigrator(pDatabase *ormshift.Database, pConfig MigratorConfig) (*Migrator, error) {
	if pDatabase == nil {
		return nil, fmt.Errorf("database cannot be nil")
	}

	lAppliedMigrationNames, lError := getAppliedMigrationNames(*pDatabase, pConfig)
	if lError != nil {
		return nil, fmt.Errorf("failed to get applied migration names: %w", lError)
	}
	lAppliedMigrations := make(map[string]bool, len(lAppliedMigrationNames))
	for _, name := range lAppliedMigrationNames {
		lAppliedMigrations[name] = true
	}

	return &Migrator{
		database:          pDatabase,
		config:            pConfig,
		migrations:        []Migration{},
		appliedMigrations: lAppliedMigrations,
	}, nil
}

func (m *Migrator) Add(pMigration Migration) {
	m.migrations = append(m.migrations, pMigration)
}

func (m *Migrator) ApplyAllMigrations() error {
	for _, lMigration := range m.migrations {
		lMigrationName := reflect.TypeOf(lMigration).Name()
		if !m.isApplied(lMigrationName) {
			lError := lMigration.Up(m)
			if lError != nil {
				return fmt.Errorf("failed to apply migration %q: %w", lMigrationName, lError)
			}
			lError = m.recordAppliedMigration(lMigrationName)
			if lError != nil {
				return fmt.Errorf("failed to record applied migration %q: %w", lMigrationName, lError)
			}
			m.appliedMigrations[lMigrationName] = true
		}
	}
	return nil
}

func (m *Migrator) RevertLastAppliedMigration() error {
	for i := len(m.migrations) - 1; i >= 0; i-- {
		lMigration := m.migrations[i]
		lMigrationName := reflect.TypeOf(lMigration).Name()
		if m.isApplied(lMigrationName) {
			lError := lMigration.Down(m)
			if lError != nil {
				return fmt.Errorf("failed to revert migration %q: %w", lMigrationName, lError)
			}
			lError = m.deleteAppliedMigration(lMigrationName)
			if lError != nil {
				return fmt.Errorf("failed to delete applied migration %q: %w", lMigrationName, lError)
			}
			delete(m.appliedMigrations, lMigrationName)
			return nil
		}
	}
	return nil
}

func (m Migrator) Database() *ormshift.Database {
	return m.database
}

func (m Migrator) Migrations() []Migration {
	return m.migrations
}

func (m Migrator) AppliedMigrations() []Migration {

	lMigrations := []Migration{}
	for _, migration := range m.Migrations() {
		name := reflect.TypeOf(migration).Name()
		if m.appliedMigrations[name] {
			lMigrations = append(lMigrations, migration)
		}
	}
	return lMigrations
}

func (m Migrator) isApplied(pMigrationName string) bool {
	_, exists := m.appliedMigrations[pMigrationName]
	return exists
}

func (m Migrator) recordAppliedMigration(pMigrationName string) error {
	q, p := m.database.SQLBuilder().InsertWithValues(
		m.config.tableName,
		ormshift.ColumnsValues{
			m.config.migrationNameColumn: pMigrationName,
			m.config.appliedAtColumn:     time.Now().UTC(),
		},
	)
	_, lError := m.database.SQLExecutor().Exec(q, p...)
	return lError
}

func (m Migrator) deleteAppliedMigration(pMigrationName string) error {
	q, p := m.database.SQLBuilder().DeleteWithValues(
		m.config.tableName,
		ormshift.ColumnsValues{
			m.config.migrationNameColumn: pMigrationName,
		},
	)
	_, lError := m.database.SQLExecutor().Exec(q, p...)
	return lError
}

func getAppliedMigrationNames(pDatabase ormshift.Database, pConfig MigratorConfig) (rMigrationNames []string, rError error) {
	rError = ensureMigrationsTableExists(pDatabase, pConfig)
	if rError != nil {
		return
	}

	q, p := pDatabase.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
		fmt.Sprintf(
			"select %s from %s order by %s",
			pConfig.migrationNameColumn,
			pConfig.tableName,
			pConfig.migrationNameColumn,
		),
	)
	lMigrationsRows, rError := pDatabase.SQLExecutor().Query(q, p...)
	if rError != nil {
		return
	}
	defer func() {
		if err := lMigrationsRows.Close(); err != nil && rError == nil {
			rError = err
		}
	}()
	for lMigrationsRows.Next() {
		var lMigrationName string
		rError = lMigrationsRows.Scan(&lMigrationName)
		if rError != nil {
			break
		}
		rMigrationNames = append(rMigrationNames, lMigrationName)
	}
	return
}

func ensureMigrationsTableExists(pDatabase ormshift.Database, pConfig MigratorConfig) error {
	lMigrationsTable := schema.NewTable(pConfig.tableName)
	if pDatabase.DBSchema().HasTable(lMigrationsTable.Name()) {
		return nil
	}
	lError := lMigrationsTable.AddColumns(
		schema.NewColumnParams{
			Name:       pConfig.migrationNameColumn,
			Type:       schema.Varchar,
			Size:       pConfig.migrationNameMaxLength,
			PrimaryKey: true,
			NotNull:    true,
		},
		schema.NewColumnParams{
			Name:    pConfig.appliedAtColumn,
			Type:    schema.DateTime,
			NotNull: true,
		},
	)
	if lError != nil {
		return lError
	}

	_, lError = pDatabase.SQLExecutor().Exec(pDatabase.SQLBuilder().CreateTable(lMigrationsTable)) // NOSONAR go:S2077 - Dynamic SQL is controlled and sanitized internally
	return lError
}
