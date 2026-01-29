package migrations_test

import (
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/migrations"
)

func TestNewMigratorConfigDefaults(t *testing.T) {
	config := migrations.NewMigratorConfig()

	testutils.AssertEqualWithLabel(t, "__ormshift_migrations", config.TableName(), "MigratorConfig.TableName")
	testutils.AssertEqualWithLabel(t, "name", config.MigrationNameColumn(), "MigratorConfig.MigrationNameColumn")
	testutils.AssertEqualWithLabel(t, "applied_at", config.AppliedAtColumn(), "MigratorConfig.AppliedAtColumn")
	testutils.AssertEqualWithLabel(t, uint(250), config.MigrationNameMaxLength(), "MigratorConfig.MigrationNameMaxLength")
}

func TestNewMigratorConfigCustom(t *testing.T) {
	config := migrations.NewMigratorConfig().
		WithTableName("custom_migrations").
		WithColumnNames("migration_name", "applied_on").
		WithMigrationNameMaxLength(500)

	testutils.AssertEqualWithLabel(t, "custom_migrations", config.TableName(), "MigratorConfig.TableName")
	testutils.AssertEqualWithLabel(t, "migration_name", config.MigrationNameColumn(), "MigratorConfig.MigrationNameColumn")
	testutils.AssertEqualWithLabel(t, "applied_on", config.AppliedAtColumn(), "MigratorConfig.AppliedAtColumn")
	testutils.AssertEqualWithLabel(t, uint(500), config.MigrationNameMaxLength(), "MigratorConfig.MigrationNameMaxLength")
}
