package api

import (
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cors"
	"gopkg.in/macaron.v1"
	apiModel "minio/api/apiModel"
	"sync"
)

var lock = sync.RWMutex{}

func InitApiRoutes(m *macaron.Macaron) {
	m.Options("*", cors.CORS())
	m.Post("/photo/load", cors.CORS(), LoadFile)
	m.Get("/photo/:year/:month/:day/:name", cors.CORS(), GetFile)
	m.Get("/delete/:year/:month/:day/:name", cors.CORS(), DeleteFile)
	m.Post("/photo/new",binding.Json([]apiModel.NewFile{}),cors.CORS(),NewFile)
}
