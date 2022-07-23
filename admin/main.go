package main

import (
	"admin/api"
	service "admin/services"
	"fmt"
	"github.com/go-macaron/auth"
	"gopkg.in/macaron.v1"
	"net/http"
)


func runApi() error {
	cfgPath := "/etc/AppConfig.ini"
	service.NewContext(cfgPath, true)
	macaron.SetConfig(cfgPath)
	m := macaron.New()
	m.Use(macaron.Renderer())
	m.Use(macaron.Logger())
	m.Use(macaron.Renderer())
	m.Use(macaron.Recovery())
	m.Use(auth.BasicFunc(func(username, password string) bool {
		return auth.SecureCompare(username, "hjkliu") && auth.SecureCompare(password, "brb40min")
	}))
	m.SetDefaultCookieSecret(service.Cfg.CookieSecret)
	api.InitApiRoutes(m)
	return 	http.ListenAndServe(service.Cfg.Host+":"+service.Cfg.AdminPort, m)
}


func main() {
	err := runApi()
	if err != nil {
		fmt.Println("failed to start api")
	}
}
