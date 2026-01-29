package migrations

type MigratorConfig struct {
	table                  string
	migrationNameColumn    string
	migrationNameMaxLength uint
	appliedAtColumn        string
}

func NewMigratorConfig() *MigratorConfig {
	config := MigratorConfig{
		table:                  "__ormshift_migrations",
		migrationNameColumn:    "name",
		migrationNameMaxLength: 250,
		appliedAtColumn:        "applied_at",
	}
	return &config
}

func (mc *MigratorConfig) WithTableName(table string) *MigratorConfig {
	mc.table = table
	return mc
}
func (mc *MigratorConfig) WithColumnNames(migrationNameColumn, appliedAtColumn string) *MigratorConfig {
	mc.migrationNameColumn = migrationNameColumn
	mc.appliedAtColumn = appliedAtColumn
	return mc
}
func (mc *MigratorConfig) WithMigrationNameMaxLength(maxLength uint) *MigratorConfig {
	mc.migrationNameMaxLength = maxLength
	return mc
}

func (mc *MigratorConfig) TableName() string {
	return mc.table
}
func (mc *MigratorConfig) MigrationNameColumn() string {
	return mc.migrationNameColumn
}
func (mc *MigratorConfig) MigrationNameMaxLength() uint {
	return mc.migrationNameMaxLength
}
func (mc *MigratorConfig) AppliedAtColumn() string {
	return mc.appliedAtColumn
}
