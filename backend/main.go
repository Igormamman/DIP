package main

import (
	"fmt"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"photoservice/backend/api"
	service "photoservice/backend/services"
)

var L = &lumberjack.Logger{Filename: "./logs/photoService.log",
	MaxSize:    50,
	MaxBackups: 2,
	MaxAge:     28,
	Compress:   true}


func runApi() error {
	cfgPath := "/etc/AppConfig.ini"
	service.NewContext(cfgPath, true)
	macaron.SetConfig(cfgPath)
	log.SetOutput(L)
	log.SetFlags(log.Lshortfile)
	log.SetFlags(log.Ldate)
	log.SetFlags(log.Ltime)
	log.Println("LOG_INFO: Starting listening:",service.Cfg.ApiPort)
	m := macaron.New()
	m.Use(macaron.Renderer())
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.SetDefaultCookieSecret(service.Cfg.CookieSecret)
	api.InitApiRoutes(m)
	service.InitMinio()
	go service.SetupWebhook()
	return http.ListenAndServe(service.Cfg.Host+":"+service.Cfg.ApiPort, m)
}

func main() {
	macaron.MaxMemory = int64(30 * 1024 * 1024)
	binding.MaxMemory = int64(30 * 1024 * 1024)
	err := runApi()
	if err != nil {
		fmt.Println("failed to start api")
	}
}
