package groupMigrate

import (
	"app/database"
	model "app/modules/group/model"
)

func MigrateTbl() bool {
	db := database.DB

	db.AutoMigrate(&model.Group{})

	return true
}
