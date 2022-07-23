package services

import (
	models "admin/models"
	"log"

	"gopkg.in/ini.v1"
)

var Cfg = models.AppConfig{}

func NewContext(path string, initDB bool) {
	dat, err := ini.Load(path)
	if err != nil {
		log.Panic("Error while reading confiuration file ", path)
	}
	err = dat.MapTo(&Cfg)
	if err != nil {
		log.Panic("Error while parsing confiuration file ", path)
	}
	if initDB {
		initDatabaseService()
	}
}

