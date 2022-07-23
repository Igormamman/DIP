package main

import (
	"gopkg.in/macaron.v1"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"minio/api"
	service "minio/services"
	"net/http"
)

func runApi() error {
	cfgPath := "/etc/AppConfig.ini"
	service.NewContext(cfgPath, true)
	macaron.SetConfig(cfgPath)
	log.SetOutput(L)
	log.SetFlags(log.Lshortfile)
	log.SetFlags(log.Ldate)
	log.SetFlags(log.Ltime)
	log.Println("LOG_INFO: Starting listening:",service.Cfg.MinioPort)
	m := macaron.New()
	m.Use(macaron.Renderer())
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	api.InitApiRoutes(m)
	service.InitMinio()
	return 	http.ListenAndServe(service.Cfg.Host+":"+service.Cfg.MinioPort, m)
}

var L = &lumberjack.Logger{Filename: "./logs/minioService.log",
	MaxSize:    50,
	MaxBackups: 2,
	MaxAge:     28,
	Compress:   true}

func main() {
}
