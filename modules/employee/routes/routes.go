package routes

import (
	"app/middleware"

	"app/modules/team/controller"

	"github.com/gofiber/fiber/v2"
)

func InitTeamRoutes(app *fiber.App) {
	team := app.Group("/team", middleware.AppInfo, middleware.AppAuthen)

	getList := team.Group("")
	getList.Get("/", controller.GetTeam)
	getList.Get("/all",controller.GetAllTeam) //error nếu swap get :id trước
	getList.Get("/:id", controller.GetTeamByID)


	
	team.Post("/", controller.CreateTeam)
	team.Put("/", controller.UpdateTeam)
	team.Delete("/:id", controller.DeleteTeam)
	team.Put("/restore/:id",controller.RestoreTeam)
}
