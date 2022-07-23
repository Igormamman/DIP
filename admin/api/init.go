package api

import (
	"github.com/go-macaron/cors"
	"gopkg.in/macaron.v1"
)

func InitApiRoutes(m *macaron.Macaron) {
	m.Options("*", cors.CORS())
	m.Get("/", cors.CORS(), getAdmin)
	m.Get("/all/update",cors.CORS(),UpdateAll)
}
