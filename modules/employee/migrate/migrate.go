package teamMigrate

import (
	"app/database"
	model "app/modules/team/model"
)

func MigrateTbl() bool {
	db := database.DB

	db.AutoMigrate(&model.Team{})

	return true
}
