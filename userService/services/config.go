package services

import (
	"log"
	models "photoservice/userService/models"

	"gopkg.in/ini.v1"
)

var Cfg = models.AppConfig{}

func NewContext(path string, initDB bool) {
	dat, err := ini.Load(path)
	if err != nil {
		log.Panic("LOG_PANIC: Error while reading configuration file ", path)
	}
	err = dat.MapTo(&Cfg)
	if err != nil {
		log.Panic("LOG_PANIC: Error while parsing configuration file ", path)
	}
	if initDB {
		initDatabaseService()
	}
}
