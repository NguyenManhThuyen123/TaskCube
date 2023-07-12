package middleware

import (
	"app/config"
	"app/database"
	"app/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Check Key App
func AppInfo(c *fiber.Ctx) error {
	appKey := c.Get("x-csv-key")
	response := new(config.DataResponse)

	if len(appKey) == 0 || appKey != config.Config("APP_KEY") {
		response.Status = false
		response.Message = config.GetMessageCode("KEY_NOT_FOUND")
		return c.JSON(response)
	}

	return c.Next()
}

// Authen
func AppAuthen(c *fiber.Ctx) error {
	store := database.Store
	isError := 0
	response := new(config.DataResponse)

	// Check token valid
	tokenData, err := utils.ExtractTokenData(c)
	if err != nil {
		response.Status = false
		response.Message = config.GetMessageCode("TOKEN_INCORRECT")
		return c.JSON(response)
	}

	// Check session Exist and comparse token
	sess, err := store.Get(tokenData.Username)
	authen := strings.Split(c.Get("x-csv-token"), " ")
	if err != nil || len(sess) == 0 || string(sess) != authen[1] {
		isError = 1
	}

	if isError == 1 {
		response.Status = false
		response.Message = config.GetMessageCode("TOKEN_INCORRECT")
		return c.JSON(response)
	}
	return c.Next()
}
