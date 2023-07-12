package main

import (
	"app/config"
	"app/database"
	"app/modules"
	"app/routes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"flag"
)

// @title Demo CSV API
// @version 1.0
// @description Language: Golang. Core: Fiber

// @contact.name CubeSystem Viet Nam
// @contact.url https://vn-cubesystem.com/
// @contact.email info@vn-cubesystem.com
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-csv-key
// @securityDefinitions.apikey ApiTokenAuth
// @in header
// @name x-csv-token
func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PATCH,PUT,DELETE",
	}))

	// Static folder
	app.Static("/assets", "./assets")

	// Connect Database
	if !database.Connect() || !modules.MigrateModule() {
		return
	}

	// Init router
	routes.InitRoutes(app)
	modules.InitRoutes(app)

	// Handle Error
	dirPath := "./assets/log/system"
	fileName := fmt.Sprintf("%s/%s.txt", dirPath, time.Now().Format("2006-01-02"))
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// dir_path does not exist
		os.MkdirAll(dirPath, os.ModePerm)
	}
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${ip} | ${method} | ${status} - ${error} | ${path} \n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   config.Config("APP_TIME_ZONE"),
		Output:     file,
	}))

	// Run app
	// port := config.Config("ENV_PORT")

	// addr := flag.String("addr", ":"+port, "http service address")
	addr := flag.String("addr", "192.168.11.162:8080", "http service address")
	flag.Parse()
	log.Fatal(app.Listen(*addr))
}
