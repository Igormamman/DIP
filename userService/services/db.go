package services

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	dbModel "photoservice/userService/models"
)

type DatabaseService struct {
	db *gorm.DB
}

var DB DatabaseService

func initDatabaseService() {
	db, err := gorm.Open(postgres.Open(Cfg.DatabaseConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("LOG_FATAL: Error: database connection failed - " + err.Error())
	}else{
		fmt.Println("DB SUCCESFULL SETUP")
	}
	DB.db = db
	db.AutoMigrate(
		&dbModel.UserCookie{},
		&dbModel.UserAccess{},
		&dbModel.Race{},
	)
}
