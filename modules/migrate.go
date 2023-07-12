package modules

import (
	"app/modules/authen/migrate"
	"app/modules/department/migrate"
	"app/modules/group/migrate"
	"app/modules/team/migrate"
)

func MigrateModule() bool {
	migrate.MigrateAuthen()
	departmentMigrate.MigrateTbl()
	groupMigrate.MigrateTbl()
	teamMigrate.MigrateTbl()
	return true
}
