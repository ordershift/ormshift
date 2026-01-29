package migrations

type MigratorConfig struct {
	table        string
	nameCol      string
	nameMaxLen   uint
	appliedAtCol string
}

func NewMigratorConfig() *MigratorConfig {
	config := MigratorConfig{
		table:        "__ormshift_migrations",
		nameCol:      "name",
		nameMaxLen:   250,
		appliedAtCol: "applied_at",
	}
	return &config
}

func (mc *MigratorConfig) WithTableName(table string) *MigratorConfig {
	mc.table = table
	return mc
}
func (mc *MigratorConfig) WithColumnNames(nameCol, appliedAtCol string) *MigratorConfig {
	mc.nameCol = nameCol
	mc.appliedAtCol = appliedAtCol
	return mc
}
func (mc *MigratorConfig) WithMigrationNameMaxLength(maxLength uint) *MigratorConfig {
	mc.nameMaxLen = maxLength
	return mc
}

func (mc *MigratorConfig) TableName() string {
	return mc.table
}
func (mc *MigratorConfig) MigrationNameColumn() string {
	return mc.nameCol
}
func (mc *MigratorConfig) MigrationNameMaxLength() uint {
	return mc.nameMaxLen
}
func (mc *MigratorConfig) AppliedAtColumn() string {
	return mc.appliedAtCol
}
