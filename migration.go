package ormshift

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

const (
	migrations_table_name  = "_db_migrations_"
	migrations_column_name = "migration_name"
)

type Migration interface {
	Up(pMigrationManager *MigrationManager) error
	Down(pMigrationManager *MigrationManager) error
}

func Migrate(pDB *sql.DB, pDriverDB DriverDB, pMigations ...Migration) (*MigrationManager, error) {
	lMigrationManager, lError := NewMigrationManager(pDB, pDriverDB)
	if lError != nil {
		return nil, lError
	}
	for _, lMigration := range pMigations {
		lMigrationManager.Add(lMigration)
	}
	lError = lMigrationManager.UpAll()
	if lError != nil {
		return nil, lError
	}
	return lMigrationManager, nil
}

type MigrationManager struct {
	db                  *sql.DB
	sqlBuilder          SQLBuilder
	dbSchema            *DBSchema
	migrations          []Migration
	upedMigrationsNames []string
}

func NewMigrationManager(pDB *sql.DB, pDriverDB DriverDB) (*MigrationManager, error) {
	if pDB == nil {
		return nil, errors.New("sql.DB cannot be nil")
	}
	if !pDriverDB.IsValid() {
		return nil, errors.New("driver db should be valid")
	}
	lUpedMigrationsNames, lError := upedMigrationsNames(pDB, pDriverDB.SQLBuilder())
	if lError != nil {
		return nil, lError
	}
	lDBSchema, lError := NewDBSchema(pDB)
	if lError != nil {
		return nil, lError
	}
	return &MigrationManager{
		db:                  pDB,
		sqlBuilder:          pDriverDB.SQLBuilder(),
		dbSchema:            lDBSchema,
		migrations:          []Migration{},
		upedMigrationsNames: lUpedMigrationsNames,
	}, nil
}

func (mm *MigrationManager) Add(pMigration Migration) {
	mm.migrations = append(mm.migrations, pMigration)
}

func (mm *MigrationManager) UpAll() error {
	for _, lMigration := range mm.migrations {
		lMigrationName := reflect.TypeOf(lMigration).Name()
		if !mm.uped(lMigrationName) {
			lError := lMigration.Up(mm)
			if lError != nil {
				return fmt.Errorf("on up migration %q: %w", lMigrationName, lError)
			}
			lError = mm.insertUpedMigration(lMigrationName)
			if lError != nil {
				return fmt.Errorf("on insert uped migration %q: %w", lMigrationName, lError)
			}
			mm.upedMigrationsNames = append(mm.upedMigrationsNames, lMigrationName)
		}
	}
	return nil
}

func (mm *MigrationManager) DownLast() error {
	if len(mm.upedMigrationsNames) > 0 {
		lLastMigrationIndex := len(mm.upedMigrationsNames) - 1
		lLastMigrationName := mm.upedMigrationsNames[lLastMigrationIndex]
		for _, lMigration := range mm.migrations {
			lMigrationName := reflect.TypeOf(lMigration).Name()
			if strings.EqualFold(lMigrationName, lLastMigrationName) {
				lError := lMigration.Down(mm)
				if lError != nil {
					return fmt.Errorf("on down migration %q: %w", lMigrationName, lError)
				}
				lError = mm.deleteDownedMigration(lMigrationName)
				if lError != nil {
					return fmt.Errorf("on delete downed migration %q: %w", lMigrationName, lError)
				}
				mm.upedMigrationsNames = mm.upedMigrationsNames[:lLastMigrationIndex]
				break
			}
		}
	}
	return nil
}

func (mm MigrationManager) DB() *sql.DB {
	return mm.db
}

func (mm MigrationManager) SQLBuilder() SQLBuilder {
	return mm.sqlBuilder
}

func (mm MigrationManager) DBSchema() *DBSchema {
	return mm.dbSchema
}

func (mm MigrationManager) UpedMigrationsNames() []string {
	return mm.upedMigrationsNames
}

func (mm MigrationManager) uped(pMigrationName string) bool {
	return slices.Contains(mm.upedMigrationsNames, pMigrationName)
}

func (mm MigrationManager) insertUpedMigration(pMigrationName string) error {
	q, p := mm.sqlBuilder.SelectWithValues(
		migrations_table_name,
		[]string{migrations_column_name},
		ColumnsValues{migrations_column_name: pMigrationName},
	)
	lMigrationsRows, lError := mm.db.Query(q, p...)
	if lError != nil {
		return lError
	}
	defer lMigrationsRows.Close()
	if !lMigrationsRows.Next() {
		q, p := mm.sqlBuilder.InsertWithValues(
			migrations_table_name,
			ColumnsValues{migrations_column_name: pMigrationName},
		)
		_, lError = mm.db.Exec(q, p...)
		if lError != nil {
			return lError
		}
	}
	return nil
}

func (mm MigrationManager) deleteDownedMigration(pNomeDaMigracao string) error {
	q, p := mm.sqlBuilder.DeleteWithValues(
		migrations_table_name,
		ColumnsValues{migrations_column_name: pNomeDaMigracao},
	)
	_, lError := mm.db.Exec(q, p...)
	if lError != nil {
		return lError
	}
	return nil
}

func upedMigrationsNames(pDB *sql.DB, pSQLBuilder SQLBuilder) ([]string, error) {
	lUpedMigrationsNames := []string{}
	lMigrationsTable, lError := NewTable(migrations_table_name)
	if lError != nil {
		return nil, lError
	}
	lDBSchema, lError := NewDBSchema(pDB)
	if lError != nil {
		return nil, lError
	}
	if !lDBSchema.ExistsTable(lMigrationsTable.Name()) {
		lMigrationsTable.AddColumn(NewColumnParams{
			Name:       migrations_column_name,
			Type:       Varchar,
			Size:       250,
			PrimaryKey: true,
			NotNull:    true,
		})
		_, lError = pDB.Exec(pSQLBuilder.CreateTable(*lMigrationsTable)) //NOSONAR go:S2077
		if lError != nil {
			return nil, lError
		}
	}
	q, p := pSQLBuilder.InteroperateSQLCommandWithNamedArgs(
		fmt.Sprintf(
			"select %s from %s order by %s",
			migrations_column_name,
			migrations_table_name,
			migrations_column_name,
		),
	)
	lMigrationsRows, lError := pDB.Query(q, p...)
	if lError != nil {
		return nil, lError
	}
	defer lMigrationsRows.Close()
	for lMigrationsRows.Next() {
		var lMigrationName string
		lError = lMigrationsRows.Scan(&lMigrationName)
		if lError != nil {
			break
		}
		lUpedMigrationsNames = append(lUpedMigrationsNames, lMigrationName)
	}
	return lUpedMigrationsNames, nil
}
