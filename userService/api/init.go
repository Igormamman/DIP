package api

import (
	"github.com/go-macaron/cors"
	"gopkg.in/macaron.v1"
)

func InitApiRoutes(m *macaron.Macaron) {
	m.Options("*", cors.CORS())
	m.Get("/logout", cors.CORS(), LogoutRouter)
	m.Get("/races", cors.CORS(), getRacesRouter)
	m.Get("/getUID", cors.CORS(), getUserIDRouter)
	m.Get("/getUserInfo", cors.CORS(), getUserInfoRouter)
	m.Get("/getRaceInfo",cors.CORS(),getRaceInfoRouter)
	m.Get("/getRaceAccess",cors.CORS(),getUserAccessRouter)
}
