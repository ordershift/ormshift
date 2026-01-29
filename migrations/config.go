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

func (mc *MigratorConfig) WithTableName(pTableName string) *MigratorConfig {
	mc.tableName = pTableName
	return mc
}
func (mc *MigratorConfig) WithColumnNames(pMigrationNameColumn, pAppliedAtColumn string) *MigratorConfig {
	mc.migrationNameColumn = pMigrationNameColumn
	mc.appliedAtColumn = pAppliedAtColumn
	return mc
}
func (mc *MigratorConfig) WithMigrationNameMaxLength(pMaxLength uint) *MigratorConfig {
	mc.migrationNameMaxLength = pMaxLength
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
