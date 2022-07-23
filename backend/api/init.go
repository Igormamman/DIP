package api

import (
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cors"
	"gopkg.in/macaron.v1"
	json "photoservice/backend/api/apiModel"
)

func InitApiRoutes(m *macaron.Macaron) {
	m.Options("*", cors.CORS())
	m.Post("/ml-webhook", binding.Json(json.ClientResultResponse{}), cors.CORS(), webhookRouter)
	m.Post("/photo/task/new", binding.Json(json.TaskData{}), cors.CORS(), newTaskRouter)
	m.Post("/photo/task/meta", binding.Json(json.UploadMeta{}), cors.CORS(), newUploadMetaRouter)
	m.Get("/file/get", cors.CORS(), GetFile)
	m.Post("/ml/task/create", binding.Json(json.NewFile{}), cors.CORS(), createTask)
	m.Post("/file/upload", cors.CORS(), newPhotoRouter)
	m.Get("/meta", cors.CORS(), getPhotoMeta)
	m.Get("/getPreviews", cors.CORS(), getPhotoPreviewUrls)
	m.Get("/deletePhoto", cors.CORS(), deletePhotoRouter)
}
