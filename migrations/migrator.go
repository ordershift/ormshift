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
	config            *MigratorConfig
	migrations        []Migration
	appliedMigrations map[string]bool
}

func NewMigrator(database *ormshift.Database, config *MigratorConfig) (*Migrator, error) {
	if database == nil {
		return nil, fmt.Errorf("database cannot be nil")
	}
	if config == nil {
		return nil, fmt.Errorf("migrator config cannot be nil")
	}

	appliedMigrationNames, err := getAppliedMigrationNames(database, config)
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migration names: %w", err)
	}
	appliedMigrations := make(map[string]bool, len(appliedMigrationNames))
	for _, name := range appliedMigrationNames {
		appliedMigrations[name] = true
	}

	return &Migrator{
		database:          database,
		config:            config,
		migrations:        []Migration{},
		appliedMigrations: appliedMigrations,
	}, nil
}

func (m *Migrator) Add(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrator) ApplyAllMigrations() error {
	for _, migration := range m.migrations {
		migrationName := reflect.TypeOf(migration).Name()
		if !m.isApplied(migrationName) {
			err := migration.Up(m)
			if err != nil {
				return fmt.Errorf("failed to apply migration %q: %w", migrationName, err)
			}
			err = m.recordAppliedMigration(migrationName)
			if err != nil {
				return fmt.Errorf("failed to record applied migration %q: %w", migrationName, err)
			}
			m.appliedMigrations[migrationName] = true
		}
	}
	return nil
}

func (m *Migrator) RevertLastAppliedMigration() error {
	for i := len(m.migrations) - 1; i >= 0; i-- {
		migration := m.migrations[i]
		migrationName := reflect.TypeOf(migration).Name()
		if m.isApplied(migrationName) {
			err := migration.Down(m)
			if err != nil {
				return fmt.Errorf("failed to revert migration %q: %w", migrationName, err)
			}
			err = m.deleteAppliedMigration(migrationName)
			if err != nil {
				return fmt.Errorf("failed to delete applied migration %q: %w", migrationName, err)
			}
			delete(m.appliedMigrations, migrationName)
			return nil
		}
	}
	return nil
}

func (m *Migrator) Database() *ormshift.Database {
	return m.database
}

func (m *Migrator) Migrations() []Migration {
	return m.migrations
}

func (m *Migrator) AppliedMigrations() []Migration {
	migrations := []Migration{}
	for _, migration := range m.Migrations() {
		name := reflect.TypeOf(migration).Name()
		if m.appliedMigrations[name] {
			migrations = append(migrations, migration)
		}
	}
	return migrations
}

func (m *Migrator) isApplied(migrationName string) bool {
	_, exists := m.appliedMigrations[migrationName]
	return exists
}

func (m *Migrator) recordAppliedMigration(migrationName string) error {
	q, p := m.database.SQLBuilder().InsertWithValues(
		m.config.table,
		ormshift.ColumnsValues{
			m.config.migrationNameColumn: migrationName,
			m.config.appliedAtColumn:     time.Now().UTC(),
		},
	)
	_, err := m.database.SQLExecutor().Exec(q, p...)
	return err
}

func (m *Migrator) deleteAppliedMigration(migrationName string) error {
	q, p := m.database.SQLBuilder().DeleteWithValues(
		m.config.table,
		ormshift.ColumnsValues{
			m.config.migrationNameColumn: migrationName,
		},
	)
	_, err := m.database.SQLExecutor().Exec(q, p...)
	return err
}

func getAppliedMigrationNames(database *ormshift.Database, config *MigratorConfig) (migrationNames []string, err error) {
	err = ensureMigrationsTableExists(database, config)
	if err != nil {
		return
	}

	q, p := database.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
		fmt.Sprintf(
			"select %s from %s order by %s",
			database.SQLBuilder().QuoteIdentifier(config.migrationNameColumn),
			database.SQLBuilder().QuoteIdentifier(config.table),
			database.SQLBuilder().QuoteIdentifier(config.migrationNameColumn),
		),
	)
	migrationsRows, err := database.SQLExecutor().Query(q, p...)
	if err != nil {
		return
	}
	defer func() {
		if closeErr := migrationsRows.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	for migrationsRows.Next() {
		var migrationName string
		err = migrationsRows.Scan(&migrationName)
		if err != nil {
			break
		}
		migrationNames = append(migrationNames, migrationName)
	}
	return
}

func ensureMigrationsTableExists(database *ormshift.Database, config *MigratorConfig) error {
	migrationsTable := schema.NewTable(config.TableName())
	if database.DBSchema().HasTable(migrationsTable.Name()) {
		return nil
	}
	err := migrationsTable.AddColumns(
		schema.NewColumnParams{
			Name:       config.MigrationNameColumn(),
			Type:       schema.Varchar,
			Size:       config.MigrationNameMaxLength(),
			PrimaryKey: true,
			NotNull:    true,
		},
		schema.NewColumnParams{
			Name:    config.AppliedAtColumn(),
			Type:    schema.DateTime,
			NotNull: true,
		},
	)
	if err != nil {
		return err
	}

	_, err = database.SQLExecutor().Exec(database.SQLBuilder().CreateTable(migrationsTable))
	return err
}
