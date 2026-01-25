package migrations

type MigratorConfig struct {
	tableName              string
	migrationNameColumn    string
	migrationNameMaxLength uint
	appliedAtColumn        string
}

func NewMigratorConfig(pOptions ...func(*MigratorConfig)) MigratorConfig {
	lConfig := MigratorConfig{
		tableName:              "__ormshift_migrations",
		migrationNameColumn:    "name",
		migrationNameMaxLength: 250,
		appliedAtColumn:        "applied_at",
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

func WithMigrationNameMaxLength(pMaxLength uint) func(*MigratorConfig) {
	return func(mc *MigratorConfig) {
		mc.migrationNameMaxLength = pMaxLength
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
func (mc MigratorConfig) MigrationNameMaxLength() uint {
	return mc.migrationNameMaxLength
}
