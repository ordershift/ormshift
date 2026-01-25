package migrations_test

import (
	"testing"

	"github.com/ordershift/ormshift/internal/testutils"
	"github.com/ordershift/ormshift/migrations"
)

func TestNewMigratorConfigDefaults(t *testing.T) {
	lConfig := migrations.NewMigratorConfig()

	testutils.AssertEqualWithLabel(t, "__ormshift_migrations", lConfig.TableName(), "MigratorConfig.TableName")
	testutils.AssertEqualWithLabel(t, "name", lConfig.MigrationNameColumn(), "MigratorConfig.MigrationNameColumn")
	testutils.AssertEqualWithLabel(t, "applied_at", lConfig.AppliedAtColumn(), "MigratorConfig.AppliedAtColumn")
	testutils.AssertEqualWithLabel(t, uint(250), lConfig.MigrationNameMaxLength(), "MigratorConfig.MigrationNameMaxLength")
}

func TestNewMigratorConfigCustom(t *testing.T) {
	lConfig := migrations.NewMigratorConfig().
		WithTableName("custom_migrations").
		WithColumnNames("migration_name", "applied_on").
		WithMigrationNameMaxLength(500)

	testutils.AssertEqualWithLabel(t, "custom_migrations", lConfig.TableName(), "MigratorConfig.TableName")
	testutils.AssertEqualWithLabel(t, "migration_name", lConfig.MigrationNameColumn(), "MigratorConfig.MigrationNameColumn")
	testutils.AssertEqualWithLabel(t, "applied_on", lConfig.AppliedAtColumn(), "MigratorConfig.AppliedAtColumn")
	testutils.AssertEqualWithLabel(t, uint(500), lConfig.MigrationNameMaxLength(), "MigratorConfig.MigrationNameMaxLength")
}
