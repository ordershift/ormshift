package migrations

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

func (mc MigratorConfig) TableName() string {
	return mc.tableName
}
func (mc MigratorConfig) MigrationNameColumn() string {
	return mc.migrationNameColumn
}
func (mc MigratorConfig) AppliedAtColumn() string {
	return mc.appliedAtColumn
}
func (mc MigratorConfig) MaxMigrationNameLength() uint {
	return mc.maxMigrationNameLength
}
