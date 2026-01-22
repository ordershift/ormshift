package migrations

import (
	"database/sql"
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/ordershift/ormshift"
	"github.com/ordershift/ormshift/schema"
)

type Migration interface {
	Up(pMigrator *Migrator) error
	Down(pMigrator *Migrator) error
}

type MigratorConfig struct {
	tableName              string
	migrationNameColumn    string
	appliedAtColumn        string
	maxMigrationNameLength uint
}

func NewMigratorConfig(pOptions ...func(*MigratorConfig)) MigratorConfig {
	lConfig := MigratorConfig{
		tableName:              "__ormshift_migrations",
		migrationNameColumn:    "name",
		appliedAtColumn:        "applied_at",
		maxMigrationNameLength: 250,
	}
	for _, o := range pOptions {
		o(&lConfig)
	}
	return lConfig
}

func WithTableName(pTableName string) func(*MigratorConfig) {
	return func(mc *MigratorConfig) {
		mc.tableName = pTableName
	}
}

func WithColumnNames(pMigrationNameColumn, pAppliedAtColumn string) func(*MigratorConfig) {
	return func(mc *MigratorConfig) {
		mc.migrationNameColumn = pMigrationNameColumn
		mc.appliedAtColumn = pAppliedAtColumn
	}
}

func WithMaxMigrationNameLength(pMaxLength uint) func(*MigratorConfig) {
	return func(mc *MigratorConfig) {
		mc.maxMigrationNameLength = pMaxLength
	}
}

func Migrate(pDatabase ormshift.Database, pConfig MigratorConfig, pMigrations ...Migration) (*Migrator, error) {
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

type Migrator struct {
	db         *sql.DB
	sqlBuilder ormshift.SQLBuilder
	dbSchema   schema.DBSchema
	config     MigratorConfig
	// TODO: Can't this be a dictionary / set for faster lookups?
	migrations            []Migration
	appliedMigrationNames []string
}

func NewMigrator(pDatabase ormshift.Database, pConfig MigratorConfig) (*Migrator, error) {
	lError := pDatabase.Validate()
	if lError != nil {
		return nil, fmt.Errorf("invalid database: %w", lError)
	}
	lAppliedMigrationNames, lError := getAppliedMigrationNames(pDatabase, pConfig)
	if lError != nil {
		return nil, lError
	}
	return &Migrator{
		db:                    pDatabase.DB(),
		sqlBuilder:            pDatabase.SQLBuilder(),
		dbSchema:              pDatabase.DBSchema(),
		config:                pConfig,
		migrations:            []Migration{},
		appliedMigrationNames: lAppliedMigrationNames,
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
				return fmt.Errorf("on migration up %q: %w", lMigrationName, lError)
			}
			lError = m.recordAppliedMigration(lMigrationName)
			if lError != nil {
				return fmt.Errorf("on recording migration %q: %w", lMigrationName, lError)
			}
			m.appliedMigrationNames = append(m.appliedMigrationNames, lMigrationName)
		}
	}
	return nil
}

func (m *Migrator) RevertLatestMigration() error {
	if len(m.appliedMigrationNames) > 0 {
		lLatestMigrationIndex := len(m.appliedMigrationNames) - 1
		lLatestMigrationName := m.appliedMigrationNames[lLatestMigrationIndex]
		for _, lMigration := range m.migrations {
			lMigrationName := reflect.TypeOf(lMigration).Name()
			if lMigrationName == lLatestMigrationName {
				lError := lMigration.Down(m)
				if lError != nil {
					return fmt.Errorf("on migration down %q: %w", lMigrationName, lError)
				}
				lError = m.deleteAppliedMigration(lMigrationName)
				if lError != nil {
					return fmt.Errorf("on deleting applied migration %q: %w", lMigrationName, lError)
				}
				m.appliedMigrationNames = m.appliedMigrationNames[:lLatestMigrationIndex]
				break
			}
		}
	}
	return nil
}

func (m Migrator) DB() *sql.DB {
	return m.db
}

func (m Migrator) SQLBuilder() ormshift.SQLBuilder {
	return m.sqlBuilder
}

func (m Migrator) DBSchema() schema.DBSchema {
	return m.dbSchema
}

func (m Migrator) AppliedMigrationNames() []string {
	return m.appliedMigrationNames
}

func (m Migrator) isApplied(pMigrationName string) bool {
	return slices.Contains(m.appliedMigrationNames, pMigrationName)
}

func (m Migrator) recordAppliedMigration(pMigrationName string) error {
	q, p := m.sqlBuilder.InsertWithValues(
		m.config.tableName,
		ormshift.ColumnsValues{
			m.config.migrationNameColumn: pMigrationName,
			m.config.appliedAtColumn:     time.Now().UTC(),
		},
	)
	_, lError := m.db.Exec(q, p...)
	return lError
}

func (m Migrator) deleteAppliedMigration(pMigrationName string) error {
	q, p := m.sqlBuilder.DeleteWithValues(
		m.config.tableName,
		ormshift.ColumnsValues{
			m.config.migrationNameColumn: pMigrationName,
		},
	)
	_, lError := m.db.Exec(q, p...)
	return lError
}

func getAppliedMigrationNames(pDatabase ormshift.Database, pConfig MigratorConfig) ([]string, error) {
	var lAppliedMigrationNames []string

	lError := ensureMigrationsTableExists(pDatabase, pConfig)
	if lError != nil {
		return nil, lError
	}

	q, p := pDatabase.SQLBuilder().InteroperateSQLCommandWithNamedArgs(
		fmt.Sprintf(
			"select %s from %s order by %s",
			pConfig.migrationNameColumn,
			pConfig.tableName,
			pConfig.migrationNameColumn,
		),
	)
	lMigrationsRows, lError := pDatabase.DB().Query(q, p...)
	if lError != nil {
		return nil, lError
	}
	defer func() {
		if err := lMigrationsRows.Close(); err != nil && lError == nil {
			lError = err
		}
	}()
	for lMigrationsRows.Next() {
		var lMigrationName string
		lError = lMigrationsRows.Scan(&lMigrationName)
		if lError != nil {
			break
		}
		lAppliedMigrationNames = append(lAppliedMigrationNames, lMigrationName)
	}
	return lAppliedMigrationNames, lError
}

func ensureMigrationsTableExists(pDatabase ormshift.Database, pConfig MigratorConfig) error {
	lMigrationsTable, lError := schema.NewTable(pConfig.tableName)
	if lError != nil {
		return lError
	}
	if !pDatabase.DBSchema().ExistsTable(lMigrationsTable.Name()) {
		columns := []schema.NewColumnParams{
			{
				Name:       pConfig.migrationNameColumn,
				Type:       schema.Varchar,
				Size:       pConfig.maxMigrationNameLength,
				PrimaryKey: true,
				NotNull:    true,
			},
			{
				Name:    pConfig.appliedAtColumn,
				Type:    schema.DateTime,
				NotNull: true,
			},
		}

		for _, col := range columns {
			if err := lMigrationsTable.AddColumn(col); err != nil {
				return err
			}
		}

		_, lError = pDatabase.DB().Exec(pDatabase.SQLBuilder().CreateTable(*lMigrationsTable)) // NOSONAR go:S2077 - Dynamic SQL is controlled and sanitized internally
	}
	return lError
}
