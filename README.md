# ormshift

[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ordershift_ormshift&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ordershift_ormshift)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ordershift_ormshift&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ordershift_ormshift)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=ordershift_ormshift&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=ordershift_ormshift)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ordershift_ormshift&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ordershift_ormshift)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=ordershift_ormshift&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=ordershift_ormshift)

## Table of contents

- [ormshift](#ormshift)
	- [Table of contents](#table-of-contents)
	- [About ormshift](#about-ormshift)
		- [Goals](#goals)
		- [Supported dialects](#supported-dialects)
		- [Architecture overview](#architecture-overview)
		- [When to use ormshift](#when-to-use-ormshift)
	- [Install](#install)
	- [Getting started](#getting-started)
		- [Connecting to a database](#connecting-to-a-database)
		- [Create a table with the SQLBuilder](#create-a-table-with-the-sqlbuilder)
		- [CRUD with the SQLBuilder and SQLExecutor](#crud-with-the-sqlbuilder-and-sqlexecutor)
			- [INSERT](#insert)
			- [UPDATE](#update)
			- [SELECT](#select)
			- [DELETE](#delete)
		- [Migrations](#migrations)
			- [Example migrations](#example-migrations)
			- [Applying and reverting migrations](#applying-and-reverting-migrations)
			- [Migrator configurations](#migrator-configurations)
	- [License](#license)

## About ormshift

ormshift is a simple, transparent Go ORM focused on database-agnostic SQL generation, migrations, and schema introspection.


### Goals

- Be practical: provide readable SQL builders and concise executor abstractions for common features like migrations and CRUDs
- Be portable: support multiple SQL dialects via scoped, easily replaceable layers
- Be minimalist: avoid heavy runtime abstractions; favor explicit builders and a clear separation between SQL generation and execution

### Supported dialects

| Dialect | `database/sql` driver |
|--|--|
| SQLite | `modernc.org/sqlite` |
| PostgreSQL | `github.com/lib/pq` |
| SQL Server | `github.com/microsoft/go-mssqldb` |

### Architecture overview

- **Migrations:** utilities to run ordered schema changes with consistent error handling
- **SQL builders:** compose SQL statements in a database-agnostic core with per-driver tweaks for dialect differences
- **Drivers:** small adapters that handle connection strings formatting, and dialect-specific SQL where required
- **Executor:** utility to execute SQL commands with consistent error handling

### When to use ormshift

- When your project needs straightforward migrations and query execution
- When database-agnosticism is a priority
- Or, when you just believe it's the right choice :smile:

## Install

This package requires Go modules.

```shell
go get github.com/ordershift/ormshift
```

## Getting started

### Connecting to a database

Using SQLite in-memory:

```go
import "github.com/ordershift/ormshift/dialects/sqlite"

db, err := ormshift.OpenDatabase(sqlite.Driver(), ormshift.ConnectionParams{InMemory: true})
if err != nil {
	// handle error
}
defer db.Close()
```

Using PostgreSQL:

```go
db, err := ormshift.OpenDatabase(postgresql.Driver(), ormshift.ConnectionParams{
	Host: "localhost",
	Port: 5432,
	User: "pg",
	Password: "pwd",
	Database: "db",
})
if err != nil {
	// handle error
}
defer db.Close()
```

Using SQL Server:

```go
db, err := ormshift.OpenDatabase(sqlserver.Driver(), ormshift.ConnectionParams{
	Instance: "sqlexpress",
	User: "sa",
	Password: "pwd",
	Database: "db",
})
if err != nil {
	// handle error
}
defer db.Close()
```

### Create a table with the SQLBuilder

```go
table := schema.NewTable("users")

err := table.AddColumns(
	schema.NewColumnParams{Name: "id", Type: schema.Integer, PrimaryKey: true, AutoIncrement: true},
	schema.NewColumnParams{Name: "name", Type: schema.Varchar, Size: 50, NotNull: false},
)

if err != nil {
	// handle error
}

db.SQLExecutor().Exec(db.SQLBuilder().CreateTable(table))
```

### CRUD with the SQLBuilder and SQLExecutor

The `SQLBuilder` returns the SQL command text and driver-compatible arguments, which can be passed to the `SQLExecutor` functions.

**Example:** PostgreSQL uses positional values ($1, $2...), while SQLite returns named args (`sql.NamedArg`).

#### INSERT

```go
query, args := db.SQLBuilder().InsertWithValues(
	"product",
	ormshift.ColumnsValues{"id": 1, "name": "X", "sku": "Y"},
)
db.SQLExecutor().Exec(query, args...)
```

#### UPDATE

```go
query, args = db.SQLBuilder().UpdateWithValues(
	"product",
	[]string{"sku", "name"}, []string{"id"},
	ormshift.ColumnsValues{"id": 1, "sku": "Y", "name": "X"},
)
db.SQLExecutor().Exec(query, args...)
```

#### SELECT

```go
query, args = db.SQLBuilder().SelectWithValues(
	"product",
	[]string{"id","sku","name"},
	ormshift.ColumnsValues{"category_id": 1, "active": true},
)
rows, err := db.SQLExecutor().Query(query, args...)
if err != nil {
	// handle error
}
rows.Close()
```

#### DELETE

```go
query, args = db.SQLBuilder().DeleteWithValues(
	"product",
	ormshift.ColumnsValues{"id": 1},
)
db.SQLExecutor().Exec(query, args...)
```

### Migrations

Every migration must implement the `Up()` and `Down()` functions. To ensure idempotency, use the schema from `db.DBSchema()` to introspect the database before applying changes, avoiding errors such as duplicate columns or tables.

The applied migrations history is recorded/persisted in your database in the `__ormshift_migrations` table (this name is configurable, check [Migrator configurations](#migrator-configurations)).

To ensure migrations stay organized in your project, it's recommended to follow some naming convention that includes a sequential number, like `M0001CreateUserTable` (then `M0002`, `M0003`, ...), or part of the timestamp when it was created, such as `M202601241957CreateUserTable` (cut on the minute part).

#### Example migrations

**Migration 1:** `m0001_create_user_table.go`

```go
// M0001CreateUserTable creates the "user" table
type M0001CreateUserTable struct{}

func (m M0001CreateUserTable) Up(migrator *migrations.Migrator) error {
	db := migrator.Database()
	table := schema.NewTable("user")

	if db.DBSchema().HasTable(table.Name()) {
		// if the table already exists, nothing to do here
		return nil
	}

	err := table.AddColumns(
		schema.NewColumnParams{Name: "id", Type: schema.Integer, PrimaryKey: true, AutoIncrement: true},
		schema.NewColumnParams{Name: "name", Type: schema.Varchar, Size: 50, NotNull: false},
		schema.NewColumnParams{Name: "email", Type: schema.Varchar, Size: 120, NotNull: false},
		schema.NewColumnParams{Name: "is_active", Type: schema.Boolean, NotNull: false},
	)
	if err != nil {
		return err
	}

	_, err = db.SQLExecutor().Exec(db.SQLBuilder().CreateTable(table))
	return err
}

func (m M0001CreateUserTable) Down(migrator *migrations.Migrator) error {
	db := migrator.Database()
	tableName := "user"
	if !db.DBSchema().HasTable(tableName) {
		// if the table already doesn't exist, nothing to do here
		return nil
	}
	_, err := db.SQLExecutor().Exec(db.SQLBuilder().DropTable(tableName))
	return err
}
```

**Migration 2:** `m0001_add_updated_at_column.go`

```go
// M0002 Adds the "updated_at" column to the "user" table
type M0002AddUpdatedAtColumn struct{}

func (m M0002AddUpdatedAtColumn) Up(migrator *migrations.Migrator) error {
	db := migrator.Database()
	tableName := "user"
	col := schema.NewColumn(schema.NewColumnParams{Name: "updated_at", Type: schema.DateTime})

	if db.DBSchema().HasColumn(tableName, col.Name()) {
		// if the column already exists, nothing to do here
		return nil
	}
	
	_, err := db.SQLExecutor().Exec(db.SQLBuilder().AlterTableAddColumn(tableName, col))
	return err
}

func (m M0002AddUpdatedAtColumn) Down(migrator *migrations.Migrator) error {
	db := migrator.Database()
	tableName := "user"
	colName := "updated_at"

	if !db.DBSchema().HasColumn(tableName, colName) {
		// if the column already doesn't exist, nothing to do here
		return nil
	}
	_, err := db.SQLExecutor().Exec(db.SQLBuilder().AlterTableDropColumn(tableName, colName))
	return err
}
```

#### Applying and reverting migrations

Applying migrations in order:

```go
migrator, err := migrations.Migrate(
	db,
	migrations.NewMigratorConfig(),
	// keep appending migrations in order here:
	M0001CreateUserTable{},
	M0002AddUpdatedAtColumn{},
)
if err != nil {
	// handle error
}
```

Reverting the last applied migration:

```go
err := migrator.RevertLastAppliedMigration()
if err != nil {
	// handle error
}
```

#### Migrator configurations

When applying migrations with the `migrations.Migrate()` function, a `MigratorConfig` is expected in the arguments, and it can be quickly resolved by `migrations.NewMigratorConfig()`, which uses the default configurations:

- Table name: `__ormshift_migrations`
- Name column name: `name`
- Name column maximum length: `250`
- Timestamp column name: `applied_at`

To change any of these defaults, use the provided functions to change the values, for example:

```go
config := migrations.NewMigratorConfig().
	WithTableName("custom_migrations").
	WithColumnNames("migration_name", "applied_on").
	WithMigrationNameMaxLength(500)

migrator, err := migrations.Migrate(
	db,
	config,
	M0001CreateUserTable{},
	M0002AddUpdatedAtColumn{},
)
if err != nil {
	// handle error
}
```

## License

[Apache-2.0](./LICENSE)
