package migrations

import (
	"gorim.org/gorim/migrations"
	"gorim.org/gorim/models"
)

var MigrationInstance *migrations.Migrations

func init() {
	MigrationInstance = &migrations.Migrations{}
	MigrationInstance.Models = []interface{}{
		&models.User{},
	}
	MigrationInstance.AddOperation(
		migrations.Operation{
			Name: "add_uuid_extensions",
			Func: MigrationInstance.RunGo(AddUuidExtensions),
		},
		migrations.Operation{
			Name: migrations.MIGRATE_MODELS,
			Func: MigrationInstance.RunMigrationModels(),
		},
		migrations.Operation{
			Name: "fill_slug",
			Func: MigrationInstance.RunGo(FillSlug),
		},
		migrations.Operation{
			Name: "test_migration_success1",
			Func: MigrationInstance.RunGo(TestMigrationSuccess),
		},
		migrations.Operation{
			Name: "test_migration_errors1",
			Func: MigrationInstance.RunGo(TestMigrationError),
		},
	)
}
