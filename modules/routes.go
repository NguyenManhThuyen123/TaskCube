package modules

import (
	authenRoute "app/modules/authen/routes"
	departmentRoute "app/modules/department/routes"
	groupRoute "app/modules/group/routes"
	teamRoute "app/modules/team/routes"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	authenRoute.InitAuthenRoutes(app)
	departmentRoute.InitDepartmentRoutes(app)
	groupRoute.InitGroupRoutes(app)
	teamRoute.InitTeamRoutes(app)
}
