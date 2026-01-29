package migrations

type MigratorConfig struct {
	tableName              string
	migrationNameColumn    string
	migrationNameMaxLength uint
	appliedAtColumn        string
}

func NewMigratorConfig() *MigratorConfig {
	config := MigratorConfig{
		tableName:              "__ormshift_migrations",
		migrationNameColumn:    "name",
		migrationNameMaxLength: 250,
		appliedAtColumn:        "applied_at",
	}
	return &config
}

func (mc *MigratorConfig) WithTableName(tableName string) *MigratorConfig {
	mc.tableName = tableName
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
	return mc.tableName
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
