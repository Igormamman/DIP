package services

import (
	dbModel "admin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type DatabaseService struct {
	db *gorm.DB
}

var DB DatabaseService

func initDatabaseService() {
	db, err := gorm.Open(postgres.Open(Cfg.DatabaseConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Error: database connection failed - " + err.Error())
	}
	DB.db = db
	db.AutoMigrate(
		&dbModel.Photo{},
		&dbModel.UploadInfo{},
		&dbModel.PhotoTag{},
	)
}
