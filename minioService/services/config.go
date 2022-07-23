package services

import (
	"log"
	models "minio/models"

	"gopkg.in/ini.v1"
)

var Cfg = models.AppConfig{}

func NewContext(path string, initDB bool) {
	dat, err := ini.Load(path)
	if err != nil {
		log.Panic("LOG_ERROR: Error while reading confiuration file ", path)
	}
	err = dat.MapTo(&Cfg)
	if err != nil {
		log.Panic("LOG_ERROR: Error while parsing confiuration file ", path)
	}
}

