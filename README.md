# ormshift

[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ordershift_ormshift&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ordershift_ormshift)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ordershift_ormshift&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ordershift_ormshift)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=ordershift_ormshift&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=ordershift_ormshift)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ordershift_ormshift&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ordershift_ormshift)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=ordershift_ormshift&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=ordershift_ormshift)

## Background

The ormshift project is a lightweight Go library that provides simple, testable primitives for building and executing SQL across multiple database engines. It was created to make schema migrations, SQL generation, and query execution predictable and easy to test without the complexity of a full ORM. ormshift focuses on clarity, minimal dependencies, and explicit SQL generation so developers retain control over queries while reducing repetitive boilerplate.

Key goals

- Practical: provide readable SQL builders and concise executor abstractions for common tasks like migrations and CRUD.
- Portable: support multiple SQL dialects (PostgreSQL, SQLite, SQL Server) via small, well-scoped driver layers.
- Testable: include utilities and fake drivers to enable reliable unit tests for migration and SQL-generation logic.
- Minimalist: avoid heavy runtime abstractions; favor explicit builders and a clear separation between SQL generation and execution.

Architecture overview

- SQL builders: compose SQL statements in a database-agnostic core with per-driver tweaks for dialect differences.
- Drivers: small adapters that handle connection, parameter formatting, and dialect-specific SQL where required.
- Migration & executor: utilities to run ordered schema changes and execute queries with consistent error handling.
- Tests & fakes: comprehensive tests and fake implementations to validate SQL generation and migration behavior without a live DB.

When to use ormshift

- Projects that need straightforward migrations and query execution without a heavy ORM.
- Tools, CLIs, or services where explicit SQL generation and cross-database support are priorities.
- Codebases that benefit from small, dependency-light libraries and strong unit-testability.

## Install

This package requires Go modules.

```shell
go get github.com/ordershift/ormshift
```

## Usage

Quick examples showing common patterns.

### Connect (SQLite in-memory with SQLite Driver):

```go
db, err := sql.Open(ormshift.DriverSQLite.Name(), ormshift.DriverSQLite.ConnectionString(ormshift.ConnectionParams{InMemory: true}))
if err != nil {
	// handle error
}
defer db.Close()
```

### Postgres and SQL Server Drivers helpers:

```go
ormshift.DriverPostgresql.ConnectionString(ormshift.ConnectionParams{User: "pg", Password: "secret", DBname: "my-db"})
ormshift.DriverSQLServer.ConnectionString(ormshift.ConnectionParams{Host: "host", Port: 1433, Instance: "sqlexpress", User: "sa", Password: "pwd", DBname: "db"})
```

### Create a table with the SQL builder:

```go
tbl, err := ormshift.NewTable("user")
if err != nil {
	// handle error
}
tbl.AddColumn(ormshift.NewColumnParams{Name: "id", Type: ormshift.Varchar, Size: 32, PrimaryKey: true, NotNull: true})
tbl.AddColumn(ormshift.NewColumnParams{Name: "email", Type: ormshift.Varchar, Size: 80, PrimaryKey: true, NotNull: true})
// ... add other columns ...
db.Exec(ormshift.DriverSQLite.SQLBuilder().CreateTable(*tbl))
```

### Insert / Update / Select / Delete examples:

```go
// Insert with values (builder returns SQL and driver-compatible args)
sqlStr, args := ormshift.DriverPostgresql.SQLBuilder().InsertWithValues("product", ormshift.ColumnsValues{"id": 1, "name": "X", "sku": "Y"})
db.Exec(sqlStr, args...)

// Update with values
sqlStr, args = ormshift.DriverSQLite.SQLBuilder().UpdateWithValues("product", []string{"sku", "name"}, []string{"id"}, ormshift.ColumnsValues{"id": 1, "sku": "Y", "name": "X"})
db.Exec(sqlStr, args...)

// Select with values
sqlStr, args = ormshift.DriverSQLServer.SQLBuilder().SelectWithValues("product", []string{"id","sku","name"}, ormshift.ColumnsValues{"category_id": 1, "active": true})
rows, err := db.Query(sqlStr, args...)
if err != nil {
	// handle error
}
rows.Close()

// Delete with values
sqlStr, args = ormshift.DriverSQLite.SQLBuilder().DeleteWithValues("product", ormshift.ColumnsValues{"id": 1})
db.Exec(sqlStr, args...)
```

Notes:

- Builders return SQL and argument lists already adapted to the driver: Postgres uses positional values ($1, $2...), while SQL Server and SQLite return named args (`sql.NamedArg`). Pass the returned args directly to `db.Exec`/`db.Query`.

### Running migrations:

```go
type CreateUserTable struct{}

func (m CreateUserTable) Up(mgr *ormshift.MigrationManager) error {
	tbl, _ := ormshift.NewTable("user")
	if mgr.DBSchema().ExistsTable(tbl.Name()) {
		return nil
	}
	tbl.AddColumn(ormshift.NewColumnParams{Name: "id", Type: ormshift.Varchar, Size: 32, PrimaryKey: true, NotNull: true})
	tbl.AddColumn(ormshift.NewColumnParams{Name: "name", Type: ormshift.Varchar, Size: 50, NotNull: false})
	tbl.AddColumn(ormshift.NewColumnParams{Name: "email", Type: ormshift.Varchar, Size: 120, NotNull: false})
	tbl.AddColumn(ormshift.NewColumnParams{Name: "active", Type: ormshift.Boolean, NotNull: false})
	_, err := mgr.DB().Exec(mgr.SQLBuilder().CreateTable(*tbl))
	return err
}

func (m CreateUserTable) Down(mgr *ormshift.MigrationManager) error {
	tableName, _ := ormshift.NewTableName("user")
	if !mgr.DBSchema().ExistsTable(*tableName) {
		return nil
	}
	_, err := mgr.DB().Exec(mgr.SQLBuilder().DropTable(*tableName))
	return err
}

type AddUpdatedAtColumn struct{}

func (m AddUpdatedAtColumn) Up(mgr *ormshift.MigrationManager) error {
	tableName, _ := ormshift.NewTableName("user")
	col, _ := ormshift.NewColumn(ormshift.NewColumnParams{Name: "updated_at", Type: ormshift.DateTime})
	if mgr.DBSchema().ExistsTableColumn(*tableName, col.Name()) {
		return nil
	}
	_, err := mgr.DB().Exec(mgr.SQLBuilder().AlterTableAddColumn(*tableName, *col))
	return err
}

func (m AddUpdatedAtColumn) Down(mgr *ormshift.MigrationManager) error {
	tableName, _ := ormshift.NewTableName("user")
	colName, _ := ormshift.NewColumnName("updated_at")
	if !mgr.DBSchema().ExistsTableColumn(*tableName, *colName) {
		return nil
	}
	_, err := mgr.DB().Exec(mgr.SQLBuilder().AlterTableDropColumn(*tableName, *colName))
	return err
}

// Run migrations in order
mgr, err := ormshift.Migrate(db, ormshift.DriverSQLite, CreateUserTable{}, AddUpdatedAtColumn{})
if err != nil {
	// handle error
}

```

Each migration must implement `Up()` and `Down()` methods. Use `mgr.DBSchema()` to introspect the database before applying changes, avoiding duplicate columns or tables.

These snippets reflect the patterns exercised by the tests: build tables and columns with `NewTable`/`NewColumn`, generate SQL with the per-driver `SQLBuilder()`, execute with `database/sql`, and run ordered migrations with `Migrate()`/`MigrationManager`.
