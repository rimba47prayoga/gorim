package migrations

import (
	"github.com/rimba47prayoga/gorim.git/migrations"
	"github.com/rimba47prayoga/gorim.git/models"
)

var MigrationInstance *migrations.Migrations

func init() {
	MigrationInstance = &migrations.Migrations{}
	MigrationInstance.Models = []interface{}{
		&models.User{},
	}
	MigrationInstance.AddOperation(
		migrations.Operation{
			Name: "migrate_models",
			Func: MigrationInstance.RunMigrationModels(),
		},
	)
	MigrationInstance.AddOperation(
		migrations.Operation{
			Name: "fill_slug",
			Func: MigrationInstance.RunGo(FillSlug),
		},
	)
}
