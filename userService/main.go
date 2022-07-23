package main

import (
	"fmt"
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"photoservice/userService/api"
	service "photoservice/userService/services"
)

var L = &lumberjack.Logger{Filename: "./logs/userService.log",
	MaxSize:    50,
	MaxBackups: 2,
	MaxAge:     28,
	Compress:   true}


func runApi() error {
	cfgPath := "/etc/AppConfig.ini"
	service.NewContext(cfgPath,true)
	macaron.SetConfig(cfgPath)
	log.SetOutput(L)
	log.SetFlags(log.Lshortfile)
	log.SetFlags(log.Ldate)
	log.SetFlags(log.Ltime)
	log.Println("LOG_INFO: Starting listening:",service.Cfg.UserPort)
	m := macaron.New()
	m.Use(macaron.Renderer())
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.SetDefaultCookieSecret(service.Cfg.CookieSecret)
	api.InitApiRoutes(m)
	return http.ListenAndServe(service.Cfg.Host+":"+service.Cfg.UserPort, m)
}

func main() {
	macaron.MaxMemory = int64(30 * 1024 * 1024)
	binding.MaxMemory = int64(30 * 1024 * 1024)
	service.InitScheduler()
	err := runApi()
	if err != nil {
		fmt.Println("failed to start api")
	}
}
