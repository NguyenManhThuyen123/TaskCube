package routes

import (
	"app/middleware"

	"app/modules/group/controller"

	"github.com/gofiber/fiber/v2"
)

func InitGroupRoutes(app *fiber.App) {
	group := app.Group("/group", middleware.AppInfo, middleware.AppAuthen)

	getList := group.Group("")
	getList.Get("/", controller.GetGroup)
	getList.Get("/all",controller.GetAllGroup) //error nếu swap get :id trước
	getList.Get("/:id", controller.GetGroupByID)


	
	group.Post("/", controller.CreateGroup)
	group.Put("/", controller.UpdateGroup)
	group.Delete("/:id", controller.DeleteGroup)
	group.Put("/restore/:id",controller.RestoreGroup)
}
