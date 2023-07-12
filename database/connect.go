package database

import (
	"app/config"
	"app/core"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/storage/postgres"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var Store *postgres.Storage

//--------------------------- TYPE STRUCT -------------------------------

type Header struct {
	CreatedBy  string     `gorm:"column:created_by;size:15"`
	UpdatedBy  string     `gorm:"column:updated_by;size:15"`	
	DeletedBy  string     `gorm:"column:deleted_by;size:15"`
}

// -----------------------------------------------------------------------


func Connect() bool {
	var err error
	status := true
	dbHost := config.Config("DB_HOST")
	dbPort := config.Config("DB_PORT")
	dbUser := config.Config("DB_USER")
	dbPassword := config.Config("DB_PASSWORD")
	dbName := config.Config("DB_NAME")
	dbSsh := config.Config("DB_SSH")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, dbSsh)

	DB, err = gorm.Open(postgresDriver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		status = false
		core.WriteLog("ERROR | DATABASE CONNECT")
	}

	if !status {
		return false
	}

	ConfigSession()

	return status
}

func ConfigSession() {

	host := config.Config("DB_HOST")
	port := config.Config("DB_PORT")
	user := config.Config("DB_USER")
	password := config.Config("DB_PASSWORD")
	name := config.Config("DB_NAME")
	sshmode := config.Config("SSH")
	post, _ := strconv.Atoi(port)

	Store = postgres.New(postgres.Config{

		Host:       host,
		Port:       post,
		Username:   user,
		Password:   password,
		Database:   name,
		Table:      "session",
		Reset:      false,
		GCInterval: 10 * time.Second,
		SslMode:    sshmode,
	})

}
